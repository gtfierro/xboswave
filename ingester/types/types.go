package types

import (
	"crypto/sha1"
	"encoding/base64"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Extract func(uri SubscriptionURI, msg xbospb.XBOS, add func(ExtractedTimeseries) error) error

type SubscriptionURI struct {
	// WAVE namespace (base64-encoded)
	Namespace string
	// '/'-delimited resource path to subscribe to
	Resource string
}

func (uri SubscriptionURI) Hash() string {
	h := sha1.New()
	h.Write([]byte(uri.Namespace))
	h.Write([]byte(uri.Resource))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

var NoMatch = errors.New("No Match")

type ExtractedTimeseries struct {
	// values extracted from the message
	Values []float64
	// corresponding times for each above value
	Times []int64
	// engineering units
	Unit string
	// BTRDB specific below this point
	// stream identifier
	UUID uuid.UUID
	// possibly temporary properties
	Annotations map[string]string
	// permanent properties
	Tags    map[string]string
	IntTags map[string]int64
	// collection name
	Collection string
}

func (ts ExtractedTimeseries) Empty() bool {
	return ts.UUID == nil
}

var _ns = uuid.Parse("d1c7c340-d0d4-11e8-a061-0cc47a0f7eea")

func GenerateUUID(uri SubscriptionURI, data []byte) uuid.UUID {
	data = append(data, []byte(uri.Namespace)...)
	data = append(data, []byte(uri.Resource)...)
	return uuid.NewSHA1(_ns, data)
}
