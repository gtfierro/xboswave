package main

import (
	"github.com/gtfierro/xboswave/driver"
	xbospb "github.com/gtfierro/xboswave/proto"
	"log"
	"time"
)

// a driver struct embeds the driver.Driver struct which provides the API
type virtual_light struct {
	*driver.Driver
	name       string
	state      bool
	brightness int
}

// The configuration contains common parameters for WAVE/WAVEMQ
func newVirtualLight(name string) *virtual_light {
	// perform any specific initialization here
	return &virtual_light{
		state:      true,
		name:       name,
		brightness: 100,
	}
}

// Implement this: Instantiate the driver.Driver using
// driver.NewDriver(cfg) and insert it into your driver struct
func (vl *virtual_light) Init(cfg driver.Config) error {
	d, err := driver.NewDriver(cfg)
	vl.Driver = d
	return err
}

// Implement this:
// add reporting and actuation functions. Block until an error occurs
func (vl *virtual_light) Start() error {
	// each instance of each service reports at the configured interval (driver.Config)
	// AddReport takes the service name, device instance identifier and callback as args
	err := vl.AddReport("virtual_light", vl.name, func() (*xbospb.XBOSIoTDeviceState, error) {
		// read from device, service, etc
		// form and return the device state as xbospb.XBOSIoTDeviceState
		reading := &xbospb.XBOSIoTDeviceState{
			// use this to communicate any device-level error that occurred during reading
			// Error: &xbospb.Error{
			// 	Msg: "Error message (if any)",
			// },
			Light: &xbospb.Light{
				State:      &xbospb.Bool{Value: vl.state},
				Brightness: &xbospb.Int64{Value: int64(vl.brightness)},
			},
		}
		return reading, nil
	})
	if err != nil {
		return err
	}

	err = vl.AddActuationCallback("virtual_light", vl.name, func(msg *xbospb.XBOSIoTDeviceActuation, received time.Time) error {
		// remember to use 'received' to filter out old actuation messages if you want
		if received.Before(time.Now().Add(-1 * time.Minute)) {
			return nil
		}

		// pull the device configuration out of the actuation message
		if msg == nil || msg.Light == nil {
			return nil
		}
		newconfiguration := msg.Light

		// perform the "backend" actuation you need
		if newconfiguration.State != nil {
			vl.state = newconfiguration.State.Value
		}
		if newconfiguration.Brightness != nil {
			value := newconfiguration.Brightness.Value
			if value > 100 {
				value = 100
			} else if value < 0 {
				value = 0
			}
			vl.brightness = int(value)
		}

		// then report the state of the device after the actuation
		reading := &xbospb.XBOSIoTDeviceState{
			Light: &xbospb.Light{
				State:      &xbospb.Bool{Value: vl.state},
				Brightness: &xbospb.Int64{Value: int64(vl.brightness)},
			},
		}
		return vl.Respond("virtual_light", vl.name, uint64(msg.Requestid), reading)
	})
	if err != nil {
		return err
	}

	// blocks until the driver implementation throws an error talking to WAVEMQ
	return vl.BlockUntilError()
}

func main() {
	cfg, err := driver.ReadConfigFromFile("params.toml")
	if err != nil {
		log.Fatal(err)
	}
	vl1 := newVirtualLight("light1")
	vl2 := newVirtualLight("light2")
	driver.Manage(cfg, vl1, vl2)
	select {}
}
