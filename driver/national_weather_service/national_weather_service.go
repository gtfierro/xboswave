package main

import (
	"encoding/json"
	"fmt"
	"github.com/gtfierro/xboswave/driver"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"log"
	"time"
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

func newDriver(stations []string, contact string, cfg driver.Config) *NationalWeatherServiceDriver {
	driver, err := driver.NewDriver(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(stations) == 0 {
		log.Fatal("No weather stations added!")
	}

	return &NationalWeatherServiceDriver{
		Driver:   driver,
		stations: stations,
		contact:  contact,
		url:      "https://api.weather.gov/stations/%s/observations/current",
		req:      gorequest.New(),
	}
}

func (driver *NationalWeatherServiceDriver) start() {

	for _, station := range driver.stations {
		fmt.Println(station)
		err := driver.AddReport("national_weather_service", station, func() (*xbospb.XBOSIoTDeviceState, error) {
			fmt.Println(station)
			datum, err := driver.read(station)
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
			log.Fatal(err)
		}
	}
}

func (driver *NationalWeatherServiceDriver) read(station string) (datum, error) {
	var d = datum{
		Station: station,
	}
	log.Println(fmt.Sprintf(driver.url, station))
	resp, _, errs := driver.req.Get(fmt.Sprintf(driver.url, station)).
		Set("User-Agent", driver.contact).
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
	cfg := driver.Config{
		Namespace:  "GyCetklhSNcgsCKVKXxSuCUZP4M80z9NRxU1pwfb2XwGhg==",
		EntityFile: "driver.ent",
		SiteRouter: "localhost:4516",
		ReportRate: 15 * time.Minute,
	}
	driver := newDriver([]string{"KOAK", "KTOA"}, "github.com/gtfierro/xboswave", cfg)

	driver.start()
	select {}
}
