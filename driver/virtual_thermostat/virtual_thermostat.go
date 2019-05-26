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

func newVirtualThermostatDriver(ntstats int, cfg driver.Config) *VirtualThermostatDriver {
	driver, err := driver.NewDriver(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &VirtualThermostatDriver{
		Driver:  driver,
		ntstats: ntstats,
		temp:    73,
		hsp:     70,
		csp:     75,
		state:   0,
	}
}

func (driver *VirtualThermostatDriver) start() {
	for tstatnum := 0; tstatnum < driver.ntstats; tstatnum++ {
		tstatnum := tstatnum
		go func() {
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
				fmt.Println(driver.Respond("virtual_thermostat", instance, uint64(msg.Requestid), reading))
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			for {
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
				log.Println("Report?", reading, driver.Report("virtual_thermostat", instance, reading))
				//driver.ReportContext(instance)
				time.Sleep(2 * time.Second)
			}
		}()
	}
}

func main() {
	cfg := driver.Config{
		Namespace:  "GyAqNmq2Clb8Z17CAIJK95EPecGHb8MSK0CpWtYCLTP1JA==",
		EntityFile: "driver.ent",
		SiteRouter: "localhost:4516",
	}
	driver := newVirtualThermostatDriver(1, cfg)

	driver.start()
	select {}
}
