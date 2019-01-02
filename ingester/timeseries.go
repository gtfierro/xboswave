package main

import (
	"context"
	"github.com/gtfierro/xboswave/ingester/types"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	logrus "github.com/sirupsen/logrus"
	"gopkg.in/btrdb.v4"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var errStreamNotExist = errors.New("Stream does not exist")

type timeseriesInterface interface {
	Flush() error
	write(extracted types.ExtractedTimeseries) error
}

type btrdbClient struct {
	addresses       []string
	conn            *btrdb.BTrDB
	streamCache     map[string]*btrdb.Stream
	streamCacheLock sync.RWMutex

	writebuffer         map[[16]byte]*types.ExtractedTimeseries
	lastwrite           map[[16]byte]time.Time
	outstandingReadings int64
	writebufferLock     sync.Mutex
}

func newBTrDBv4(c *BTrDBConfig) *btrdbClient {
	b := &btrdbClient{
		addresses:   []string{c.Address},
		streamCache: make(map[string]*btrdb.Stream),
		writebuffer: make(map[[16]byte]*types.ExtractedTimeseries),
		lastwrite:   make(map[[16]byte]time.Time),
	}
	logrus.Infof("Connecting to BtrDBv4 at addresses %v...", b.addresses)
	conn, err := btrdb.Connect(context.Background(), b.addresses...)
	if err != nil {
		logrus.Warningf("Could not connect to btrdbv4: %v", err)
		return nil
	}
	b.conn = conn
	logrus.Info("Connected to BtrDB!")

	go func() {
		for _ = range time.NewTicker(10 * time.Second).C {
			b.writebufferLock.Lock()

			var commitset []*types.ExtractedTimeseries
			for uuid, lastwritten := range b.lastwrite {
				randomDuration := 30*time.Second + time.Duration(rand.Intn(30))*time.Second
				if time.Now().Sub(lastwritten) > randomDuration {
					commitset = append(commitset, b.writebuffer[uuid])
				}
			}
			logrus.Info("BTRDB STATUS writebuffers=", len(b.writebuffer), " committingBuffers=", len(commitset), " outstanding=", atomic.LoadInt64(&b.outstandingReadings))
			b.writebufferLock.Unlock()
			for _, buf := range commitset {
				if err := b.commitBuffer(buf); err != nil {
					logrus.Error(errors.Wrapf(err, "Could not commit buffer for %s (%s)", buf.UUID.String(), buf.Collection))
				}
			}
		}
	}()

	return b
}

func (bdb *btrdbClient) Flush() error {
	var commitset []*types.ExtractedTimeseries
	bdb.writebufferLock.Lock()
	for _, buf := range bdb.writebuffer {
		commitset = append(commitset, buf)
	}
	bdb.writebufferLock.Unlock()
	var reterr error
	for _, buf := range commitset {
		if err := bdb.commitBuffer(buf); err != nil {
			if reterr != nil {
				reterr = err
			}
		}
	}
	bdb.conn.Disconnect()
	return reterr
}

// Fetch the stream object so we can read/write. This will first check the internal in-memory
// cache of stream objects, then it will check the BtrDB client cache. If the stream
// is not found there, then this method will return errStreamNotExist and a nil stream
func (bdb *btrdbClient) getStream(extracted types.ExtractedTimeseries) (stream *btrdb.Stream, err error) {
	// first check cache
	streamuuid := extracted.UUID
	bdb.streamCacheLock.RLock()
	stream, found := bdb.streamCache[streamuuid.String()]
	bdb.streamCacheLock.RUnlock()
	if found {
		return // from cache
	}
	// then check BtrDB for existing stream
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, TimeseriesOperationTimeout)
	defer cancel()
	stream = bdb.conn.StreamFromUUID(streamuuid)

	exists, existsErr := stream.Exists(ctx)
	if existsErr != nil {
		logrus.Error("problem checking stream exists", streamuuid.String(), stream.UUID().String(), extracted.Collection)
		e := btrdb.ToCodedError(existsErr)
		if e.Code != 404 && e.Code != 501 { // if error is NOT unfound or permission denied, then return it
			err = errors.Wrap(existsErr, "Could not fetch stream")
			return
		}
	} else if exists {
		//logrus.Info("exists already", streamuuid.String())
		// if stream exists, add to cache
		bdb.streamCacheLock.Lock()
		bdb.streamCache[streamuuid.String()] = stream
		bdb.streamCacheLock.Unlock()
		return
	}

	logrus.Info("Initializing stream with uuid: ", extracted.UUID.String(), extracted.Collection)
	stream, err = bdb.conn.Create(ctx, extracted.UUID, extracted.Collection, extracted.Tags, extracted.Annotations)
	if err == nil {
		//logrus.Info("Created stream", extracted.Collection)
		bdb.streamCacheLock.Lock()
		bdb.streamCache[streamuuid.String()] = stream
		bdb.streamCacheLock.Unlock()
	} else {
		e := btrdb.ToCodedError(err)
		if e.Code == 501 {
			err = nil
		} else {
			logrus.Info("Error", e.Code)
		}

	}

	return
}

func (bdb *btrdbClient) write(extracted types.ExtractedTimeseries) error {

	bdb.writebufferLock.Lock()
	buf, found := bdb.writebuffer[extracted.UUID.Array()]
	if !found {
		bdb.writebuffer[extracted.UUID.Array()] = &extracted
		bdb.lastwrite[extracted.UUID.Array()] = time.Now()
		buf = &extracted
	} else {
		buf.Values = append(buf.Values, extracted.Values...)
		buf.Times = append(buf.Times, extracted.Times...)
	}
	bdb.writebufferLock.Unlock()
	atomic.AddInt64(&bdb.outstandingReadings, int64(len(extracted.Values)))

	if len(buf.Values) >= MaxInMemoryTimeseriesBuffer {
		bdb.commitBuffer(buf)
	}

	return nil
}

func (bdb *btrdbClient) commitBuffer(buf *types.ExtractedTimeseries) error {
	start := time.Now()
	stream, err := bdb.getStream(*buf)
	if err != nil {
		return err
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, TimeseriesOperationTimeout)
	defer cancel()
	if ctx.Err() == context.Canceled {
		logrus.Warning("CANCELLED")
	}

	bdb.writebufferLock.Lock()
	bdb.lastwrite[buf.UUID.Array()] = time.Now()
	bdb.writebufferLock.Unlock()

	if len(buf.Values) == 0 {
		return nil
	}

	// delete the range before inserting
	var minTime int64 = 1 << 62
	var maxTime int64 = -1
	for _, ts := range buf.Times {
		if ts < minTime {
			minTime = ts
		}
		if ts > maxTime {
			maxTime = ts
		}
	}
	// use nearest on the start time of the range.
	if minTime != maxTime {
		rv, _, err := stream.Nearest(ctx, minTime, 0, true)
		if err != nil {
			logrus.Warning(errors.Wrapf(err, "Could not get nearest for (%s %s)", buf.UUID.String(), buf.Collection))
		}
		// if this is true, then there is data inside the range we are going to commit; in this case, delete the old data
		if rv.Time >= minTime {
			if _, err := stream.DeleteRange(ctx, minTime, maxTime); err != nil {
				logrus.Warning(errors.Wrapf(err, "Could not delete range for (%s %s)", buf.UUID.String(), buf.Collection))
			}
		}
	}
	commitDuration := time.Duration(maxTime-minTime) * time.Nanosecond
	logrus.Debug("Committing collection=", buf.Collection, " uuid=", buf.UUID.String(), "(", stream.UUID().String(), ")", " #values=", len(buf.Values), " start=", time.Unix(0, minTime), " dur=", commitDuration)
	err = stream.InsertTV(ctx, buf.Times, buf.Values)
	if err != nil {
		return errors.Wrapf(err, "Could not insert into %s", buf.Collection)
	}
	atomic.AddInt64(&bdb.outstandingReadings, -int64(len(buf.Values)))
	pointsCommitted.Add(float64(len(buf.Values)))
	commitSizes.Observe(float64(len(buf.Values)))
	commitTimes.Observe(float64(time.Since(start).Nanoseconds() / 1e6))
	buf.Values = buf.Values[:0]
	buf.Times = buf.Times[:0]
	return nil
}

type influxClient struct {
	conn influx.Client

	dbname              string
	writebuffer         map[[16]byte]*types.ExtractedTimeseries
	lastwrite           map[[16]byte]time.Time
	outstandingReadings int64
	writebufferLock     sync.Mutex
}

func newInfluxDB(c *InfluxDBConfig) *influxClient {
	logrus.Infof("Connecting to InfluxDB at address %v...", c.Address)
	conn, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     c.Address,
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "could not connect influx"))
	}
	logrus.Info("Connected to InfluxDB!")

	i := &influxClient{
		conn:        conn,
		dbname:      c.Database,
		writebuffer: make(map[[16]byte]*types.ExtractedTimeseries),
		lastwrite:   make(map[[16]byte]time.Time),
	}

	return i
}

func (inf *influxClient) write(extracted types.ExtractedTimeseries) error {

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  inf.dbname,
		Precision: "s",
	})
	if err != nil {
		return errors.Wrap(err, "could not create new batch")
	}

	start := time.Now()
	for idx, t := range extracted.Times {
		val := extracted.Values[idx]
		tags := map[string]string{
			"unit": extracted.Unit,
			"uuid": extracted.UUID.String(),
		}
		for k, v := range extracted.Annotations {
			tags[k] = v
		}
		for k, v := range extracted.Tags {
			tags[k] = v
		}
		fields := map[string]interface{}{
			"value": val,
		}

		pt, err := influx.NewPoint(extracted.Collection, tags, fields, time.Unix(0, t))
		if err != nil {
			return errors.Wrap(err, "could not create new point")
		}
		bp.AddPoint(pt)
		// increment # of points we have processed
		pointsCommitted.Inc()
	}
	commitSizes.Observe(float64(len(extracted.Values)))
	commitTimes.Observe(float64(time.Since(start).Nanoseconds() / 1e6))

	return inf.conn.Write(bp)
}

func (inf *influxClient) Flush() error {
	return nil
}
