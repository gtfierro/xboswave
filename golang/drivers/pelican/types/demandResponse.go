package types

import (
	"encoding/xml"
	"fmt"
	"time"
)

type DR_EVENT_STATUS int

const (
	DR_EVENT_STATUS_NOT_CONFIGURED DR_EVENT_STATUS = iota
	DR_EVENT_STATUS_UNUSABLE
	DR_EVENT_STATUS_INACTIVE
	DR_EVENT_STATUS_ACTIVE
)

type DR_EVENT_TYPE int

const (
	DR_EVENT_TYPE_NO_EVENT DR_EVENT_TYPE = iota
	DR_EVENT_TYPE_NORMAL
	DR_EVENT_TYPE_MODERATE
	DR_EVENT_TYPE_HIGH
	DR_EVENT_TYPE_SPECIAL
)

type ADREventWrapperAPI struct {
	Success int         `xml:"success"`
	Message string      `xml:"message"`
	Event   ADREventAPI `xml:"attribute"`
}

type ADREventAPI struct {
	End    string `xml:"OpenADREventEnd"`
	Start  string `xml:"OpenADREventStart"`
	Status string `xml:"OpenADRStatus"`
	Type   string `xml:"OpenADREventType"`
}

type ADREvent struct {
	EventEnd   int64           `msgpack:"event_end"`
	EventStart int64           `msgpack:"event_start"`
	EventType  DR_EVENT_TYPE   `msgpack:"event_type"`
	DRStatus   DR_EVENT_STATUS `msgpack:"dr_status"`
	Time       int64           `msgpack:"time"`
}

func (pel *Pelican) TrackDREvent() (*ADREvent, error) {
	resp, _, errs := pel.drReq.Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "get").
		Param("object", "Site").
		Param("value", "OpenADREventEnd;OpenADREventStart;OpenADRStatus;OpenADREventType").
		End()

	if errs != nil {
		return nil, fmt.Errorf("Error retrieving thermostat demand-response status from %s: %v", pel.target, errs)
	}

	defer resp.Body.Close()
	var result ADREventWrapperAPI
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode demand-response XML: %v", err)
	}
	if result.Success == 0 {
		return nil, fmt.Errorf("Error retrieving thermostat demand-response status from %s: %s", resp.Request.URL, result.Message)
	}

	event := result.Event
	var output ADREvent

	// Convert ADR Start Time from String to Int
	startTime, err := drTimeToUnix(event.Start, pel.timezone)
	if err != nil {
		return nil, fmt.Errorf("String to Unix Time Conversion Error: %v", err)
	}
	output.EventStart = startTime

	// Convert ADR End Time from String to Int
	endTime, err := drTimeToUnix(event.End, pel.timezone)
	if err != nil {
		return nil, fmt.Errorf("String to Unix Time Conversion Error: %v", err)
	}
	output.EventEnd = endTime

	eventStatus, statusErr := getEventStatus(event.Status)
	if statusErr != nil {
		return nil, fmt.Errorf("Event Status Error: %v", statusErr)
	}
	output.DRStatus = eventStatus

	eventType, typeErr := getEventType(event.Type)
	if typeErr != nil {
		return nil, fmt.Errorf("Event Type Error: %v", typeErr)
	}
	output.EventType = eventType

	output.Time = time.Now().UnixNano()

	return &output, nil
}

func drTimeToUnix(DRTime string, timezone *time.Location) (int64, error) {
	// Time field is empty or nil
	if len(DRTime) == 0 {
		return 0, nil
	}

	// Using Parse in Location to convert time string into correct time.Time value
	outputTime, timeErr := time.ParseInLocation("2006-01-02T15:04", DRTime, timezone)
	if timeErr != nil {
		return 0, fmt.Errorf("Error parsing %v into Time struct: %v\n", DRTime, timeErr)
	}

	return outputTime.UnixNano(), nil
}

// Map Status to Corresponding Integer Value
func getEventStatus(eventStatus string) (DR_EVENT_STATUS, error) {
	switch eventStatus {
	case "Not Configured":
		return DR_EVENT_STATUS_NOT_CONFIGURED, nil
	case "Unusable":
		return DR_EVENT_STATUS_UNUSABLE, nil
	case "Inactive":
		return DR_EVENT_STATUS_INACTIVE, nil
	case "Active":
		return DR_EVENT_STATUS_ACTIVE, nil
	default:
		return 0, fmt.Errorf("Event status not recognized")
	}
}

// Map Event Type to Corresponding Integer Value
func getEventType(eventType string) (DR_EVENT_TYPE, error) {
	switch eventType {
	case "Normal":
		return DR_EVENT_TYPE_NORMAL, nil
	case "Moderate":
		return DR_EVENT_TYPE_MODERATE, nil
	case "High":
		return DR_EVENT_TYPE_HIGH, nil
	case "Special":
		return DR_EVENT_TYPE_SPECIAL, nil
	case "None":
		return DR_EVENT_TYPE_NO_EVENT, nil
	default:
		return 0, fmt.Errorf("Event type not recognized")
	}
}
