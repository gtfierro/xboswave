package main

import (
	"github.com/gtfierro/xboswave/ingester/types"
)

type pluginlist struct {
	mapping map[string]types.Extract
}

func newPluginlist() pluginlist {
	return pluginlist{
		mapping: make(map[string]types.Extract),
	}
}

func (pl pluginlist) add(filename string, extractFunc types.Extract) {
	pl.mapping[filename] = extractFunc
}
