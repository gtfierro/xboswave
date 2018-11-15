package main

import (
	"github.com/gtfierro/xboswave/ingester/types"
	"github.com/immesys/wavemq/mqpb"
)

// TODO: these go into sqlite table
type ArchiveRequest struct {
	Schema string
	Plugin string
	URI    types.SubscriptionURI
}

type subscription struct {
	S    mqpb.WAVEMQ_SubscribeClient
	stop chan struct{}
	uri  types.SubscriptionURI
}
