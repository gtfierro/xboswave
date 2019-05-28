package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/driver"
	xbospb "github.com/gtfierro/xboswave/proto"
	"log"
	"math/rand"
	"time"
)

type VirtualThermostatDriver struct {
	*driver.Driver
	ntstats int

	temp  float64
	hsp   float64
	csp   float64
	state int
}

func newVirtualThermostatDriver(ntstats int) *VirtualThermostatDriver {
	return &VirtualThermostatDriver{
		ntstats: ntstats,
		temp:    73,
		hsp:     70,
		csp:     75,
		state:   0,
	}
}

func (drv *VirtualThermostatDriver) Init(cfg driver.Config) error {
	d, err := driver.NewDriver(cfg)
	drv.Driver = d
	return err
}

func (driver *VirtualThermostatDriver) Start() error {

	for tstatnum := 0; tstatnum < driver.ntstats; tstatnum++ {
		tstatnum := tstatnum

		instance := fmt.Sprintf("vtstat%d", tstatnum)
		err := driver.AddActuationCallback("virtual_thermostat", instance, func(msg *xbospb.XBOSIoTDeviceActuation, received time.Time) error {
			tstat := msg.Thermostat
			fmt.Println("actuation?", tstat, received)
			if tstat.HeatingSetpoint != nil {
				driver.hsp = tstat.HeatingSetpoint.Value
			}
			if tstat.CoolingSetpoint != nil {
				driver.csp = tstat.CoolingSetpoint.Value
			}
			reading := &xbospb.XBOSIoTDeviceState{
				Thermostat: &xbospb.Thermostat{
					Temperature:     &xbospb.Double{Value: driver.temp},
					HeatingSetpoint: &xbospb.Double{Value: driver.hsp},
					CoolingSetpoint: &xbospb.Double{Value: driver.csp},
					State:           &xbospb.HVACState{Value: xbospb.HVACStateValue(driver.state)},
				},
			}
			return driver.Respond("virtual_thermostat", instance, uint64(msg.Requestid), reading)
		})
		if err != nil {
			return err
		}

		err = driver.AddReport("virtual_thermostat", instance, func() (*xbospb.XBOSIoTDeviceState, error) {
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
			return reading, nil
		})
		if err != nil {
			return err
		}
	}

	return driver.BlockUntilError()
}

func main() {
	cfg, err := driver.ReadConfigFromFile("params.toml")
	if err != nil {
		log.Fatal(err)
	}
	tstats := newVirtualThermostatDriver(cfg.GetInt("num_tstats"))
	log.Fatal(driver.Manage(cfg, tstats))
}
