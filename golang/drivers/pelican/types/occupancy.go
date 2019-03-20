package types

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	OCCUPANCY_UNKNOWN = iota
	OCCUPANCY_OCCUPIED
	OCCUPANCY_UNOCCUPIED
)

type occupancyResponse struct {
	Success    int        `xml:"success"`
	Thermostat thermostat `xml:"Thermostat"`
	Message    string     `xml:"message"`
}

type thermostat struct {
	Sensors []childSensor `xml:"slaves"`
}

type childSensor struct {
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Value string `xml:"value"`
}

func (pel *Pelican) GetOccupancy() (int, error) {
	resp, _, errs := pel.occupancyReq.Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "get").
		Param("object", "thermostat").
		Param("selection", fmt.Sprintf("name:%s;", pel.Name)).
		Param("value", "slaves;").
		End()

	if errs != nil {
		return 0, fmt.Errorf("Error retrieving thermostat occupancy data: %s", errs)
	}
	defer resp.Body.Close()

	var occResp occupancyResponse
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&occResp); err != nil {
		return 0, fmt.Errorf("Failed to decode occupancy API result: %s", err)
	}
	if occResp.Success == 0 {
		return 0, fmt.Errorf("Error retrieving thermostat occupancy data: %s", occResp.Message)
	}

	status := OCCUPANCY_UNKNOWN
	for _, sensor := range occResp.Thermostat.Sensors {
		if strings.ToLower(sensor.Type) == "occupancy sensor" {
			if strings.ToLower(sensor.Value) == "occupied" {
				status = OCCUPANCY_OCCUPIED
			} else {
				status = OCCUPANCY_UNOCCUPIED
			}
			break
		}
	}
	return status, nil
}
