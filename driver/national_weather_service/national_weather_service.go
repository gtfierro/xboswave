package main

import (
	"encoding/json"
	"fmt"
	"github.com/gtfierro/xboswave/driver"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"log"
)

type datum struct {
	Resp    response
	Station string
}

type response struct {
	Properties struct {
		Temperature struct {
			UnitCode string   `json:"unitCode"`
			Value    *float64 `json:"value"`
		} `json:"temperature"`
		RelativeHumidity struct {
			UnitCode string   `json:"unitCode"`
			Value    *float64 `json:"value"`
		} `json:"relativeHumidity"`
		WindSpeed struct {
			UnitCode string   `json:"unitCode"`
			Value    *float64 `json:"value"`
		} `json:"windSpeed"`
		WindDirection struct {
			UnitCode string   `json:"unitCode"`
			Value    *float64 `json:"value"`
		} `json:"windDirection"`
		CloudLayers []struct {
			Base struct {
				UnitCode string   `json:"unitCode"`
				Value    *float64 `json:"value"`
			} `json:"base"`
			Amount string `json:"amount"`
		} `json:"cloudLayers"`
	} `json:"properties"`
}

type NationalWeatherServiceDriver struct {
	*driver.Driver
	stations []string
	url      string
	contact  string
	req      *gorequest.SuperAgent
}

func newDriver(stations []string, contact string) *NationalWeatherServiceDriver {
	if len(stations) == 0 {
		log.Fatal("No weather stations added!")
	}

	return &NationalWeatherServiceDriver{
		stations: stations,
		contact:  contact,
		url:      "https://api.weather.gov/stations/%s/observations/current",
		req:      gorequest.New(),
	}
}

func (nws *NationalWeatherServiceDriver) Init(cfg driver.Config) error {
	d, err := driver.NewDriver(cfg)
	nws.Driver = d
	return err
}

func (nws *NationalWeatherServiceDriver) Start() error {
	for _, station := range nws.stations {
		err := nws.AddReport("national_weather_service", station, func() (*xbospb.XBOSIoTDeviceState, error) {
			fmt.Println(station)
			datum, err := nws.read(station)
			if err != nil {
				return nil, err
			}
			fmt.Printf("%+v\n", datum.Resp)

			ws := &xbospb.WeatherStation{}
			if datum.Resp.Properties.Temperature.Value != nil {
				ws.Temperature = &xbospb.Double{Value: *datum.Resp.Properties.Temperature.Value}
			}
			if datum.Resp.Properties.RelativeHumidity.Value != nil {
				ws.Humidity = &xbospb.Double{Value: *datum.Resp.Properties.Temperature.Value}
			}
			if datum.Resp.Properties.WindSpeed.Value != nil {
				ws.WindSpeed = &xbospb.Double{Value: *datum.Resp.Properties.WindSpeed.Value}
			}
			if datum.Resp.Properties.WindDirection.Value != nil {
				ws.WindBearing = &xbospb.Double{Value: *datum.Resp.Properties.WindDirection.Value}
			}

			msg := &xbospb.XBOSIoTDeviceState{
				WeatherStation: ws,
			}
			return msg, nil
		})
		if err != nil {
			return err
		}
	}
	return nws.BlockUntilError()
}

func (nws *NationalWeatherServiceDriver) read(station string) (datum, error) {
	var d = datum{
		Station: station,
	}
	log.Println(fmt.Sprintf(nws.url, station))
	resp, _, errs := nws.req.Get(fmt.Sprintf(nws.url, station)).
		Set("User-Agent", nws.contact).
		Set("Accept", "*/*").
		End()
	if errs != nil {
		return d, errors.Wrap(errs[0], "Could not fetch URL")
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&d.Resp); err != nil {
		return d, errors.Wrap(err, "Could not decode response")
	}
	return d, nil

}

func main() {
	cfg, err := driver.ReadConfigFromFile("params.toml")
	if err != nil {
		log.Fatal(err)
	}
	stations := cfg.GetStringSlice("stations")
	contact := cfg.GetString("contact")
	nws := newDriver(stations, contact)
	log.Fatal(driver.Manage(cfg, nws))
}
