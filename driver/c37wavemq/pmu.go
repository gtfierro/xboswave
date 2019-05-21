package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const QueueSize = 16000
const MaxBatch = 1000

type PMU struct {
	ctx     context.Context
	address string
	id      uint16
	name    string
	conn    *net.TCPConn
	br      *bufio.Reader

	currentconfig *Config12Frame

	cfgmu sync.Mutex
	cfgs  map[uint16]*Config12Frame

	output chan *DataFrame
}

func HandleDevice(ctx context.Context, name string, target string, id uint16) chan *DataFrame {
	pmu := &PMU{
		ctx:     ctx,
		address: target,
		name:    name,
		id:      id,
		cfgs:    make(map[uint16]*Config12Frame),
		output:  make(chan *DataFrame, 1000),
	}
	go pmu.dialLoop()
	return pmu.output
}

func (p *PMU) dialLoop() {
	for {
		err := p.dial()
		lg.Infof("[%s] abort due to error: %s", p.name, err)
		time.Sleep(30 * time.Second)
	}
}
func (p *PMU) dial() error {
	var conn *net.TCPConn
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	addr, err := net.ResolveTCPAddr("tcp", p.address)
	if err != nil {
		return err
	}
	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	lg.Infof("[%s] dial succeeded\n", p.name)
	br := bufio.NewReader(conn)
	p.br = br
	p.conn = conn
	p.initialConfigure()
	return p.process()
}

func (p *PMU) process() error {
	for {
		if p.ctx.Err() != nil {
			return p.ctx.Err()
		}
		ch, framez, err := p.readFrame()
		if err != nil {
			lg.Infof("[%s] frame read error: %v\n", p.name, err)
			p.conn.Close()
			return err
		}
		_ = ch
		for _, frame := range framez {
			cfg, ok := frame.(*Config12Frame)
			if ok {
				p.cfgmu.Lock()
				p.cfgs[ch.IDCODE] = cfg
				p.cfgmu.Unlock()
				p.sendStartCommand()
			}
			dat, ok := frame.(*DataFrame)
			if ok {
				select {
				case p.output <- dat:
				default:
					lg.Warning("[%s] WARNING QUEUE OVERFLOW. DROPPING DATA FROM ID %d\n", p.name, dat.IDCODE)
				}
			}
		}
	}
}

func (p *PMU) configKeepAlive() error {
	{
		c := &CommandFrame{}
		c.IDCODE = p.id
		c.SetSOCToNow()
		c.FRACSEC = 0
		c.SetSyncType(SYNC_TYPE_CMD)
		c.FRAMESIZE = CommonHeaderLength + 4
		c.CMD = uint16(CMD_SEND_CFG2)
		return WriteChecksummedFrame(c, p.conn)
	}
}
func (p *PMU) initialConfigure() {
	{
		c := &CommandFrame{}
		c.IDCODE = p.id
		c.SetSOCToNow()
		c.FRACSEC = 0
		c.SetSyncType(SYNC_TYPE_CMD)
		c.FRAMESIZE = CommonHeaderLength + 4
		c.CMD = uint16(CMD_SEND_CFG2)
		err := WriteChecksummedFrame(c, p.conn)
		if err != nil {
			panic(err)
		}
	}
	{
		c := &CommandFrame{}
		c.IDCODE = p.id
		c.SetSOCToNow()
		c.FRACSEC = 0
		c.SetSyncType(SYNC_TYPE_CMD)
		c.FRAMESIZE = CommonHeaderLength + 4
		c.CMD = uint16(CMD_TURN_ON_TX)
		err := WriteChecksummedFrame(c, p.conn)
		if err != nil {
			panic(err)
		}
	}
	lg.Infof("[%s] completed initial configuration commands\n", p.name)
}

func (p *PMU) sendStartCommand() {
	c := &CommandFrame{}
	c.IDCODE = p.id
	c.SetSOCToNow()
	c.SetSyncType(SYNC_TYPE_CMD)
	c.FRAMESIZE = CommonHeaderLength + 4
	c.CMD = uint16(CMD_TURN_ON_TX)
	err := WriteChecksummedFrame(c, p.conn)
	if err != nil {
		panic(err)
	}
}

func (p *PMU) ConfigFor(idcode uint16) (*Config12Frame, error) {
	p.cfgmu.Lock()
	cfg, ok := p.cfgs[idcode]
	p.cfgmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("No config found")
	}
	return cfg, nil
}
func (p *PMU) readFrame() (*CommonHeader, []Frame, error) {
	err := p.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		lg.Panicf("[%s] could not set deadline", p.name)
	}
	r := p.br
	initialByte, err := r.ReadByte()
	if err != nil {
		return nil, nil, err
	}
	skipped := 0
	for initialByte != 0xAA {
		skipped++
		initialByte, err = r.ReadByte()
		if err != nil {
			return nil, nil, err
		}
	}
	if skipped != 0 {
		lg.Warningf("[%s] SYNC LOSS DETECTED, SKIPPED %d BYTES RESYNCING\n", p.name, skipped)
	}
	raw := make([]byte, CommonHeaderLength)
	nread, err := io.ReadFull(r, raw[1:])
	if err != nil {
		return nil, nil, err
	}
	if nread != len(raw[1:]) {
		lg.Warningf("READ LEN MISMATCH %d vs %d\n", nread, len(raw[1:]))
	}
	raw[0] = initialByte
	rcopy := make([]byte, CommonHeaderLength)
	copy(rcopy, raw)

	ch, err := ReadCommonHeader(bytes.NewBuffer(rcopy))
	if err != nil {
		return nil, nil, err
	}

	rest := make([]byte, int(ch.FRAMESIZE)-CommonHeaderLength)
	_, err = io.ReadFull(r, rest)
	if err != nil {
		return nil, nil, err
	}

	raw = append(raw, rest[:len(rest)-2]...)
	realchk := Checksum(raw)
	expectedchk := (int(rest[len(rest)-2]) << 8) + int(rest[len(rest)-1])
	if expectedchk != int(realchk) {
		lg.Infof("[%s] frame checksum failure type=%d, got=%x expected=%x\n", p.name, ch.SyncType(), realchk, expectedchk)
		//the spec says silently ignore frames with bad checksums
		return ch, nil, nil
	}
	if ch.SyncType() == SYNC_TYPE_CFG2 {
		cfg2, err := ReadConfig12Frame(ch, bytes.NewBuffer(rest))
		if err != nil {
			return nil, nil, err
		}
		return ch, []Frame{cfg2}, nil
	}
	if ch.SyncType() == SYNC_TYPE_DATA {
		cfg, _ := p.ConfigFor(ch.IDCODE)
		if cfg == nil {
			lg.Infof("[%s] dropping data frame: no config\n", p.name)
			return ch, nil, nil
		}
		dat, err := ReadDataFrame(ch, cfg, bytes.NewBuffer(rest))
		if err != nil {
			return nil, nil, err
		}
		return ch, []Frame{dat}, nil
	}
	if ch.SyncType() == SYNC_TYPE_CFG1 {
		cfg1, err := ReadConfig12Frame(ch, bytes.NewBuffer(rest))
		if err != nil {
			return nil, nil, err
		}
		return ch, []Frame{cfg1}, nil
	}
	if ch.SyncType() == SYNC_TYPE_CFG3 {
		lg.Infof("[%s] WARN got CFG3 which is not supported\n", p.name)
		return ch, nil, nil
	}
	return ch, nil, fmt.Errorf("Unknown frame type")
}
