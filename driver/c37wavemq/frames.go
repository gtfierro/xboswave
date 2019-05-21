package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/howeyc/crc16"
)

type SYNC_TYPE int

const (
	SYNC_TYPE_DATA   SYNC_TYPE = 0
	SYNC_TYPE_HEADER SYNC_TYPE = 1
	SYNC_TYPE_CFG1   SYNC_TYPE = 2
	SYNC_TYPE_CFG2   SYNC_TYPE = 3
	SYNC_TYPE_CFG3   SYNC_TYPE = 5
	SYNC_TYPE_CMD    SYNC_TYPE = 4
)

type CommonHeader struct {
	SYNC      uint16
	FRAMESIZE uint16
	IDCODE    uint16
	SOC       uint32
	FRACSEC   uint32
}
type Config12Frame struct {
	TIME_BASE uint32
	NUM_PMU   uint16
	Entries   []*Config12Entry
	DATA_RATE uint16
}
type Config12Entry struct {
	STN     string
	IDCODE  uint16
	FORMAT  uint16
	PHNMR   uint16
	ANNMR   uint16
	DGNMR   uint16
	PHCHNAM []string
	ANCHNAM []string
	DGCHNAM []string
	PHUNIT  []uint32
	ANUNIT  []uint32
	DGUNIT  []uint32
	FNOM    uint16
	CFGCNT  uint16
}

func ReadConfig12Frame(ch *CommonHeader, r io.Reader) (*Config12Frame, error) {
	rv := &Config12Frame{}
	subr := io.LimitReader(r, int64(ch.FRAMESIZE)-CommonHeaderLength-2)
	err := binary.Read(subr, binary.BigEndian, &rv.TIME_BASE)
	if err != nil {
		return nil, err
	}
	err = binary.Read(subr, binary.BigEndian, &rv.NUM_PMU)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(rv.NUM_PMU); i++ {
		e := &Config12Entry{}
		stn_arr := make([]byte, 16)
		_, err := subr.Read(stn_arr)
		if err != nil {
			return nil, err
		}
		zeroidx := bytes.IndexByte(stn_arr, 0)
		if zeroidx >= 0 {
			e.STN = string(bytes.TrimRight(stn_arr[:zeroidx], "\x00 "))
		} else {
			e.STN = string(bytes.TrimRight(stn_arr, "\x00 "))
		}
		err = binary.Read(subr, binary.BigEndian, &e.IDCODE)
		if err != nil {
			return nil, err
		}
		err = binary.Read(subr, binary.BigEndian, &e.FORMAT)
		if err != nil {
			return nil, err
		}
		err = binary.Read(subr, binary.BigEndian, &e.PHNMR)
		if err != nil {
			return nil, err
		}
		err = binary.Read(subr, binary.BigEndian, &e.ANNMR)
		if err != nil {
			return nil, err
		}
		err = binary.Read(subr, binary.BigEndian, &e.DGNMR)
		if err != nil {
			return nil, err
		}
		for phni := 0; phni < int(e.PHNMR); phni++ {
			barr := make([]byte, 16)
			_, err := subr.Read(barr)
			if err != nil {
				return nil, err
			}
			e.PHCHNAM = append(e.PHCHNAM, string(bytes.TrimRight(barr, "\x00 ")))
		}
		for annmri := 0; annmri < int(e.ANNMR); annmri++ {
			barr := make([]byte, 16)
			_, err := subr.Read(barr)
			if err != nil {
				return nil, err
			}
			e.ANCHNAM = append(e.ANCHNAM, string(bytes.TrimRight(barr, "\x00 ")))
		}
		for dgnmri := 0; dgnmri < int(e.DGNMR); dgnmri++ {
			for bit := 0; bit < 16; bit++ {
				barr := make([]byte, 16)
				_, err := subr.Read(barr)
				if err != nil {
					return nil, err
				}
				e.DGCHNAM = append(e.DGCHNAM, string(bytes.TrimRight(barr, "\x00 ")))
			}
		}

		for phni := 0; phni < int(e.PHNMR); phni++ {
			var unit uint32
			err = binary.Read(subr, binary.BigEndian, &unit)
			if err != nil {
				return nil, err
			}
			e.PHUNIT = append(e.PHUNIT, unit)
		}
		for annmri := 0; annmri < int(e.ANNMR); annmri++ {
			var unit uint32
			err = binary.Read(subr, binary.BigEndian, &unit)
			if err != nil {
				return nil, err
			}
			e.ANUNIT = append(e.ANUNIT, unit)
		}
		for dgnmri := 0; dgnmri < int(e.DGNMR); dgnmri++ {
			var unit uint32
			err = binary.Read(subr, binary.BigEndian, &unit)
			if err != nil {
				return nil, err
			}
			e.DGUNIT = append(e.DGUNIT, unit)
		}
		err = binary.Read(subr, binary.BigEndian, &e.FNOM)
		if err != nil {
			return nil, err
		}
		err = binary.Read(subr, binary.BigEndian, &e.CFGCNT)
		if err != nil {
			return nil, err
		}
		rv.Entries = append(rv.Entries, e)
	}
	err = binary.Read(subr, binary.BigEndian, &rv.DATA_RATE)
	if err != nil {
		return nil, err
	}
	return rv, nil
}

const CommonHeaderLength = 14

func (ch *CommonHeader) SyncType() SYNC_TYPE {
	return SYNC_TYPE((ch.SYNC >> 4) & 7)
}

func (ch *CommonHeader) Version() int {
	return int(ch.SYNC) & 15
}

func (ch *CommonHeader) SetSOCToNow() {
	ch.SOC = uint32(time.Now().UTC().Unix())
}
func (ch *CommonHeader) SetSyncType(t SYNC_TYPE) {
	ch.SYNC = (uint16(t) << 4) | 0xAA01
}
func ReadCommonHeader(r io.Reader) (*CommonHeader, error) {
	rv := &CommonHeader{}
	err := binary.Read(r, binary.BigEndian, rv)
	return rv, err
}
func ReadChecksum(r io.Reader) (uint16, error) {
	var rv uint16
	err := binary.Read(r, binary.BigEndian, &rv)
	return rv, err
}

type CMD_WORD int

const (
	CMD_TURN_OFF_TX CMD_WORD = 1
	CMD_TURN_ON_TX  CMD_WORD = 2
	CMD_SEND_HDR    CMD_WORD = 3
	CMD_SEND_CFG1   CMD_WORD = 4
	CMD_SEND_CFG2   CMD_WORD = 5
	CMD_SEND_CFG3   CMD_WORD = 6
	CMD_EXFRAME     CMD_WORD = 8
)

type CommandFrame struct {
	CommonHeader
	CMD uint16
}

type Frame interface {
}

type DataFrame struct {
	CommonHeader
	TIMEBASE     int
	UTCUnixNanos int64
	TimeQual     uint8
	Data         []*PMUData
}
type PMUData struct {
	STN    string
	IDCODE uint16
	STAT   uint16
	//Converted
	PHASOR_NAMES  []string
	PHASOR_MAG    []float64
	PHASOR_ANG    []float64
	PHASOR_ISVOLT []bool
	//Converted to absolute freq
	FREQ float64
	//Converted to floating point but not altered
	DFREQ        float64
	ANALOG_NAMES []string
	ANALOG       []float64

	DIGITAL_NAMES []string
	DIGITAL       []int
}

func Checksum(arr []byte) uint16 {
	return uint16(crc16.ChecksumCCITTFalse(arr))
}

func WriteChecksummedFrame(frame Frame, w io.Writer) error {
	dat := &bytes.Buffer{}
	err := binary.Write(dat, binary.BigEndian, frame)
	if err != nil {
		return err
	}
	chk := Checksum(dat.Bytes())
	dat.WriteByte(byte(chk >> 8))
	dat.WriteByte(byte(chk & 0xFF))
	_, err = w.Write(dat.Bytes())
	return err
}

func ReadPhasor(format uint16, unit uint32, r io.Reader) (mag float64, ang float64, isvolt bool, err error) {
	isvolt = false
	if unit>>24 == 0 {
		isvolt = true
	}
	if format&1 == 0 {
		//Rectangular coordinates
		if format&2 == 0 {
			//16 bit integer
			var real uint16
			var imag uint16
			err := binary.Read(r, binary.BigEndian, &real)
			if err != nil {
				return 0, 0, false, err
			}
			err = binary.Read(r, binary.BigEndian, &imag)
			if err != nil {
				return 0, 0, false, err
			}
			freal := float64(real)
			fimag := float64(imag)
			freal = freal * float64(unit&0xFFFFFF) * 10e-5
			fimag = fimag * float64(unit&0xFFFFFF) * 10e-5

			phase := math.Atan2(fimag, freal)
			degrees := (phase / (2 * math.Pi)) * 360
			mag := math.Hypot(fimag, freal)
			return mag, degrees, isvolt, nil

		} else {
			//32 bit float
			var real float32
			var imag float32
			err := binary.Read(r, binary.BigEndian, &real)
			if err != nil {
				return 0, 0, false, err
			}
			err = binary.Read(r, binary.BigEndian, &imag)
			if err != nil {
				return 0, 0, false, err
			}
			phase := math.Atan2(float64(imag), float64(real))
			degrees := (phase / (2 * math.Pi)) * 360
			mag := math.Hypot(float64(imag), float64(real))
			return mag, degrees, isvolt, nil
		}
	} else {
		//Polar coordinates
		if format&2 == 0 {
			//16 bit integer
			var mag uint16
			var ang uint16
			err := binary.Read(r, binary.BigEndian, &mag)
			if err != nil {
				return 0, 0, false, err
			}
			err = binary.Read(r, binary.BigEndian, &ang)
			if err != nil {
				return 0, 0, false, err
			}

			fmag := float64(mag) * float64(unit&0xFFFFFF) * 10e-5
			fang := float64(ang) * 10e-4 //radians
			degrees := (fang / (2 * math.Pi)) * 360
			return fmag, degrees, isvolt, nil
		} else {
			//32 bit float
			var mag float32
			var ang float32
			err := binary.Read(r, binary.BigEndian, &mag)
			if err != nil {
				return 0, 0, false, err
			}
			err = binary.Read(r, binary.BigEndian, &ang)
			if err != nil {
				return 0, 0, false, err
			}
			//fmt.Printf("read raw float32 mag=%f ang=%f\n", mag, ang)
			fmag := float64(mag)
			fang := float64(ang)
			degrees := (fang / (2 * math.Pi)) * 360
			return fmag, degrees, isvolt, nil
		}
	}
}
func ReadAnalog(format uint16, unit uint32, r io.Reader) (float64, error) {
	//TODO analog scaling
	if format&4 == 0 {
		//16 bit integer
		var rv uint16
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return 0, err
		}
		return float64(rv), nil
	} else {
		//32 bit float
		var rv float32
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return 0, err
		}
		return float64(rv), nil
	}
}
func ReadDigital(unit uint32, r io.Reader) (int, error) {
	var rv uint16
	err := binary.Read(r, binary.BigEndian, &rv)
	if err != nil {
		return -1, err
	}
	return int(rv), nil
}
func ReadFreq(format uint16, fnom float64, r io.Reader) (float64, error) {
	if format&8 == 0 {
		//16 bit integer
		var rv uint16
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return -1, err
		}
		return (float64(rv) / 10e3) + fnom, nil
	} else {
		//32 bit float
		var rv float32
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return 0, err
		}
		return float64(rv), nil
	}
}
func ReadROCOF(format uint16, r io.Reader) (float64, error) {
	if format&8 == 0 {
		//16 bit integer
		var rv uint16
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return -1, err
		}
		return (float64(rv) / 10e3), nil
	} else {
		//32 bit float
		var rv float32
		err := binary.Read(r, binary.BigEndian, &rv)
		if err != nil {
			return 0, err
		}
		return float64(rv), nil
	}
}
func FreqFieldToHz(freq uint16) float64 {
	if freq == 0 {
		return 60
	}
	if freq == 1 {
		return 50
	}
	panic("invalid FNOM")
}
func ReadDataFrame(ch *CommonHeader, cfg *Config12Frame, r io.Reader) (*DataFrame, error) {
	rv := DataFrame{CommonHeader: *ch}
	nanos := (float64(ch.FRACSEC&0xFFFFFF) / float64(cfg.TIME_BASE)) * 1e9
	rv.UTCUnixNanos = int64(ch.SOC)*1e9 + int64(nanos)
	rv.TimeQual = uint8(ch.FRACSEC >> 24)
	rv.TIMEBASE = int(cfg.TIME_BASE)
	for pmu := 0; pmu < int(cfg.NUM_PMU); pmu++ {
		d := PMUData{}
		err := binary.Read(r, binary.BigEndian, &d.STAT)
		if err != nil {
			return nil, err
		}
		entry := cfg.Entries[pmu]
		d.STN = entry.STN
		d.IDCODE = entry.IDCODE
		for phi := 0; phi < int(entry.PHNMR); phi++ {
			mag, ang, isvolt, err := ReadPhasor(entry.FORMAT, entry.PHUNIT[phi], r)
			if err != nil {
				return nil, err
			}
			d.PHASOR_ANG = append(d.PHASOR_ANG, ang)
			d.PHASOR_MAG = append(d.PHASOR_MAG, mag)
			d.PHASOR_ISVOLT = append(d.PHASOR_ISVOLT, isvolt)
			d.PHASOR_NAMES = append(d.PHASOR_NAMES, entry.PHCHNAM[phi])
		}
		fr, err := ReadFreq(entry.FORMAT, FreqFieldToHz(entry.FNOM), r)
		if err != nil {
			return nil, err
		}
		d.FREQ = fr
		rocof, err := ReadROCOF(entry.FORMAT, r)
		if err != nil {
			return nil, err
		}
		d.DFREQ = rocof
		for anni := 0; anni < int(entry.ANNMR); anni++ {
			val, err := ReadAnalog(entry.FORMAT, entry.ANUNIT[anni], r)
			if err != nil {
				return nil, err
			}
			d.ANALOG = append(d.ANALOG, val)
			d.ANALOG_NAMES = append(d.ANALOG_NAMES, entry.ANCHNAM[anni])
		}
		for digi := 0; digi < int(entry.DGNMR); digi++ {
			val, err := ReadDigital(entry.DGUNIT[digi], r)
			if err != nil {
				return nil, err
			}
			d.DIGITAL = append(d.DIGITAL, val)
			d.DIGITAL_NAMES = append(d.DIGITAL_NAMES, fmt.Sprintf("DIGITAL%02d", digi))
		}
		rv.Data = append(rv.Data, &d)
	}
	return &rv, nil
}
