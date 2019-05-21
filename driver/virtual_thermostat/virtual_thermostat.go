package main

import (
	"github.com/gtfierro/xboswave/driver"
	xbospb "github.com/gtfierro/xboswave/proto"
	"log"
	"math/rand"
	"time"
)

type VirtualThermostatDriver struct {
	*driver.Driver
	temp  float64
	hsp   float64
	csp   float64
	state int
}

func newVirtualThermostatDriver(cfg driver.Config) *VirtualThermostatDriver {
	driver, err := driver.NewDriver(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &VirtualThermostatDriver{
		Driver: driver,
		temp:   73,
		hsp:    70,
		csp:    75,
		state:  0,
	}
}

func main() {
	cfg := driver.Config{
		Namespace:  "GyAqNmq2Clb8Z17CAIJK95EPecGHb8MSK0CpWtYCLTP1JA==",
		EntityFile: "driver.ent",
		SiteRouter: "localhost:4516",
	}
	driver := newVirtualThermostatDriver(cfg)

	for {
		time.Sleep(2 * time.Second)

		// update thermostat state
		if driver.state == 2 {
			driver.temp -= rand.Float64()
		} else if driver.state == 1 {
			driver.temp += rand.Float64()
		} else {
			driver.temp += 0.5 - rand.Float64()
		}

		if driver.temp <= driver.hsp-1 {
			driver.state = 1
		} else if driver.temp >= driver.csp+1 {
			driver.state = 2
		} else {
			driver.state = 0
		}

		reading := &xbospb.XBOSIoTDeviceState{
			Thermostat: &xbospb.Thermostat{
				Temperature:     &xbospb.Double{Value: driver.temp},
				HeatingSetpoint: &xbospb.Double{Value: driver.hsp},
				CoolingSetpoint: &xbospb.Double{Value: driver.csp},
				State:           &xbospb.HVACState{Value: xbospb.HVACStateValue(driver.state)},
			},
		}
		log.Println("Report?", reading, driver.Report("test", reading))
		driver.ReportContext()
	}

	_ = driver
}
