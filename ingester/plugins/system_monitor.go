package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

func has_system_status(msg xbospb.XBOS) bool {
	return msg.BasicServerStatus != nil
}

var units = map[string]string{
	"cpu_load": "percent",
	"phys_mem_available": "bytes",
	"disk_usage": "percent",
	"disk_available": "bytes",
}

var lookup = map[string]func(msg xbospb.XBOS) (float64, bool) {
	//"cpu_load": func(msg xbospb.XBOS) (float64, bool) {
	//	if has_system_status(msg) && msg.BasicServerStatus.CpuLoad != nil {
	//		return float64(msg.BasicServerStatus.CpuLoad), true
	//	}
	//	return 0, false
	//},
	"phys_mem_available": func(msg xbospb.XBOS) (float64, bool) {
		if has_system_status(msg) && msg.BasicServerStatus.PhysMemAvailable != nil {
			return float64(msg.BasicServerStatus.PhysMemAvailable.Value), true
		}
		return 0, false
	},
	"disk_usage": func(msg xbospb.XBOS) (float64, bool) {
		if has_system_status(msg) && msg.BasicServerStatus.DiskUsage != nil {
			return float64(msg.BasicServerStatus.DiskUsage.Value), true
		}
		return 0, false
	},
	"disk_available": func(msg xbospb.XBOS) (float64, bool) {
		if has_system_status(msg) && msg.BasicServerStatus.DiskAvailable != nil {
			return float64(msg.BasicServerStatus.DiskAvailable.Value), true
		}
		return 0, false
	},
}

func build(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {
	if extractfunc, found := lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.BasicServerStatus.Time)
			extracted.Values = append(extracted.Values, value)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(name))
			extracted.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
			extracted.Tags = map[string]string{
				"unit": units[name],
				"name": name,
				"hostname": msg.BasicServerStatus.Hostname,
			}
			return extracted
		}
	}

	return types.ExtractedTimeseries{}
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if has_system_status(msg) {
		for name := range lookup {
			extracted := build(uri, name, msg)
			if err := add(extracted); err != nil {
				return err
			}
		}

		for cpuidx, cpu := range msg.BasicServerStatus.CpuLoad {
			var extracted types.ExtractedTimeseries
			name := "cpu_load"
			value := cpu.Value
			time := int64(msg.BasicServerStatus.Time)
			extracted.Values = append(extracted.Values, value)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(name))
			extracted.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
			extracted.Tags = map[string]string{
				"unit": units[name],
				"name": name,
				"cpu_id": fmt.Sprintf("%d", cpuidx),
				"hostname": msg.BasicServerStatus.Hostname,
			}
			if err := add(extracted); err != nil {
				return err
			}
		}
	}
	return nil
}
