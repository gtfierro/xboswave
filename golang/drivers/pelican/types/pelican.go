package types

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

var modeNameMappings = map[string]int32{
	"Off":  0,
	"Heat": 1,
	"Cool": 2,
	"Auto": 3,
}
var modeValMappings = []string{"Off", "Heat", "Cool", "Auto"}

var stateMappings = map[string]int32{
	"Off":         0,
	"Heat-Stage1": 1,
	"Heat-Stage2": 4,
	"Cool-Stage1": 2,
	"Cool-Stage2": 5,
}

// TODO Support case where the thermostat is configured to use Celsius

type Pelican struct {
	username      string
	password      string
	sitename      string
	id            string
	Name          string
	HeatingStages int32
	CoolingStages int32
	TimezoneName  string
	target        string
	cookieTime    time.Time
	timezone      *time.Location
	cookie        *http.Cookie
	req           *gorequest.SuperAgent
	drReq         *gorequest.SuperAgent
	occupancyReq  *gorequest.SuperAgent
	scheduleReq   *gorequest.SuperAgent
}

type PelicanStatus struct {
	Temperature       float64 `msgpack:"temperature"`
	RelHumidity       float64 `msgpack:"relative_humidity"`
	HeatingSetpoint   float64 `msgpack:"heating_setpoint"`
	CoolingSetpoint   float64 `msgpack:"cooling_setpoint"`
	Override          bool    `msgpack:"override"`
	Fan               bool    `msgpack:"fan"`
	Mode              int32   `msgpack:"mode"`
	State             int32   `msgpack:"state"`
	EnabledHeatStages int32   `msgpack:"enabled_heat_stages"`
	EnabledCoolStages int32   `msgpack:"enabled_cool_stages"`
	Time              int64   `msgpack:"time"`
}

type PelicanSetpointParams struct {
	HeatingSetpoint *float64
	CoolingSetpoint *float64
}

type PelicanStateParams struct {
	HeatingSetpoint *float64
	CoolingSetpoint *float64
	Override        *float64
	Mode            *float64
	Fan             *float64
}

type PelicanStageParams struct {
	HeatingStages *int32
	CoolingStages *int32
}

// Thermostat Object API Result Structs
type apiResult struct {
	Thermostat apiThermostat `xml:"Thermostat"`
	Success    int32         `xml:"success"`
	Message    string        `xml:"message"`
}

type apiThermostat struct {
	Temperature     float64 `xml:"temperature"`
	RelHumidity     int32   `xml:"humidity"`
	HeatingSetpoint int32   `xml:"heatSetting"`
	CoolingSetpoint int32   `xml:"coolSetting"`
	SetBy           string  `xml:"setBy"`
	Schedule        string  `xml:"schedule"`
	HeatNeedsFan    string  `xml:"HeatNeedsFan"`
	System          string  `xml:"system"`
	RunStatus       string  `xml:"runStatus"`
	HeatStages      int32   `xml:"heatStages"`
	CoolStages      int32   `xml:"coolStages"`
	StatusDisplay   string  `xml:"statusDisplay"`
}

type thermostatInfo struct {
	Name          string `xml:"name"`
	HeatingStages int32  `xml:"heatStages"`
	CoolingStages int32  `xml:"coolStages"`
}

type discoverAPIResult struct {
	Thermostats []thermostatInfo `xml:"Thermostat"`
	Success     int32            `xml:"success"`
	Message     string           `xml:"message"`
}

// Thermostat History Object API Result Structs
type apiResultHistory struct {
	XMLName xml.Name   `xml:"result"`
	Success int        `xml:"success"`
	Message string     `xml:"message"`
	Records apiRecords `xml:"ThermostatHistory"`
}

type apiRecords struct {
	Name    string       `xml:"name"`
	History []apiHistory `xml:"History"`
}

type apiHistory struct {
	TimeStamp string `xml:"timestamp"`
}

// Thermostat Site Object API Result Structs
type apiResultSite struct {
	XMLName   xml.Name    `xml:"result"`
	Success   int         `xml:"success"`
	Attribute apiTimezone `xml:"attribute"`
}

type apiTimezone struct {
	Timezone string `xml:"timeZone"`
}

type NewPelicanParams struct {
	Username      string
	Password      string
	Sitename      string
	Name          string
	HeatingStages int32
	CoolingStages int32
	Timezone      string
}

func NewPelican(params *NewPelicanParams) (*Pelican, error) {
	timezone, err := time.LoadLocation(params.Timezone)
	if err != nil {
		return nil, err
	}

	newPelican := &Pelican{
		username:      params.Username,
		password:      params.Password,
		sitename:      params.Sitename,
		target:        fmt.Sprintf("https://%s.officeclimatecontrol.net/api.cgi", params.Sitename),
		Name:          params.Name,
		HeatingStages: params.HeatingStages,
		CoolingStages: params.CoolingStages,
		TimezoneName:  params.Timezone,
		timezone:      timezone,
		req:           gorequest.New(),
		drReq:         gorequest.New(),
		occupancyReq:  gorequest.New(),
		scheduleReq:   gorequest.New(),
	}
	if error := newPelican.setCookieAndID(); error != nil {
		return nil, error
	}
	return newPelican, nil
}

func DiscoverPelicans(username, password, sitename string) ([]*Pelican, error) {
	// Time zone retrieval logic
	targetTimezone := fmt.Sprintf("https://%s.officeclimatecontrol.net/api.cgi", sitename)
	respTimezone, _, errsTimezone := gorequest.New().Get(targetTimezone).
		Param("username", username).
		Param("password", password).
		Param("request", "get").
		Param("object", "Site").
		Param("value", "timeZone;").
		End()
	if errsTimezone != nil {
		return nil, fmt.Errorf("Error retrieving object result from %s: %s", targetTimezone, errsTimezone)
	}
	defer respTimezone.Body.Close()
	var resultTimezone apiResultSite
	decTimezone := xml.NewDecoder(respTimezone.Body)
	if err := decTimezone.Decode(&resultTimezone); err != nil {
		return nil, fmt.Errorf("Failed to decode response XML: %v", err)
	}
	timezoneName := resultTimezone.Attribute.Timezone

	target := fmt.Sprintf("https://%s.officeclimatecontrol.net/api.cgi", sitename)
	resp, _, errs := gorequest.New().Get(target).
		Param("username", username).
		Param("password", password).
		Param("request", "get").
		Param("object", "Thermostat").
		Param("value", "name;heatStages;coolStages").
		End()
	if errs != nil {
		return nil, fmt.Errorf("Error retrieving thermostat name from %s: %s", target, errs)
	}

	defer resp.Body.Close()
	var result discoverAPIResult
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response XML: %v", err)
	}
	if result.Success == 0 {
		return nil, fmt.Errorf("Error retrieving thermostat info from %s: %s", resp.Request.URL, result.Message)
	}

	pelicans := make([]*Pelican, 0)
	for _, thermInfo := range result.Thermostats {
		// rewrite default heat/cool stage info
		if thermInfo.HeatingStages == 0 {
			thermInfo.HeatingStages = 1
		}
		if thermInfo.CoolingStages == 0 {
			thermInfo.CoolingStages = 1
		}
		if thermInfo.Name != "" {
			newPelican, err := NewPelican(&NewPelicanParams{
				Username:      username,
				Password:      password,
				Sitename:      sitename,
				Name:          thermInfo.Name,
				HeatingStages: thermInfo.HeatingStages,
				CoolingStages: thermInfo.CoolingStages,
				Timezone:      timezoneName,
			})
			if err != nil {
				return nil, fmt.Errorf("Error creating thermostat: %s", err)
			}
			pelicans = append(pelicans, newPelican)
		}
	}
	return pelicans, nil
}

func (pel *Pelican) GetStatus() (*PelicanStatus, error) {
	resp, _, errs := pel.req.Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "get").
		Param("object", "Thermostat").
		Param("selection", fmt.Sprintf("name:%s;", pel.Name)).
		Param("value", "temperature;humidity;heatSetting;coolSetting;setBy;HeatNeedsFan;system;runStatus;statusDisplay;schedule;heatStages;coolStages").
		End()
	if errs != nil {
		return nil, fmt.Errorf("Error retrieving thermostat status from %s: %v", pel.target, errs)
	}

	defer resp.Body.Close()
	var result apiResult
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response XML: %v", err)
	}
	if result.Success == 0 {
		return nil, fmt.Errorf("Error retrieving thermostat status from %s: %s", resp.Request.URL, result.Message)
	}

	thermostat := result.Thermostat

	if thermostat.StatusDisplay == "Unreachable" {
		fmt.Printf("Thermostat %s is unreachable\n", pel.Name)
		return nil, nil
	}

	var fanState bool
	if strings.HasPrefix(thermostat.RunStatus, "Heat") {
		fanState = thermostat.HeatNeedsFan == "Yes"
	} else if thermostat.RunStatus != "Off" {
		fanState = true
	} else {
		fanState = false
	}
	thermState, ok := stateMappings[thermostat.RunStatus]
	if !ok {
		// Thermostat is not calling for heating or cooling
		if thermostat.System == "Off" {
			thermState = 0 // Off
		} else {
			// Thermostat is not heating or cooling, but fan is still running
			// Report this as off
			thermState = 0 //Off
		}
	}

	// Thermostat History Object Request to retrieve time stamps from past hour
	endTime := time.Now().In(pel.timezone).Format(time.RFC3339)
	startTime := time.Now().Add(-1 * time.Hour).In(pel.timezone).Format(time.RFC3339)

	respHist, _, errsHist := pel.req.Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "get").
		Param("object", "ThermostatHistory").
		Param("selection", fmt.Sprintf("startDateTime:%s;endDateTime:%s;", startTime, endTime)).
		Param("value", "timestamp").
		End()
	defer respHist.Body.Close()

	if errsHist != nil {
		return nil, fmt.Errorf("Error retrieving thermostat status from %s: %v", pel.target, errsHist)
	}

	var histResult apiResultHistory
	histDec := xml.NewDecoder(respHist.Body)
	if histErr := histDec.Decode(&histResult); histErr != nil {
		return nil, fmt.Errorf("Failed to decode response XML: %v", histErr)
	}
	if histResult.Success == 0 {
		return nil, fmt.Errorf("Error retrieving thermostat status from %s: %s", respHist.Request.URL, histResult.Message)
	}

	if len(histResult.Records.History) > 0 {
		// Converting string timeStamp to int64 format
		match := histResult.Records.History[len(histResult.Records.History)-1]
		timestamp, timeErr := time.ParseInLocation("2006-01-02T15:04", match.TimeStamp, pel.timezone)
		if timeErr != nil {
			return nil, fmt.Errorf("Error parsing %v into Time struct: %v\n", match.TimeStamp, timeErr)
		}

		now := time.Now()
		if timestamp.Before(now.Add(-2 * time.Hour)) {
			fmt.Println("WARNING temperature data has not changed for 2 hours. This is not necessarily an error")
		}
	}

	return &PelicanStatus{
		Temperature:       thermostat.Temperature,
		RelHumidity:       float64(thermostat.RelHumidity),
		HeatingSetpoint:   float64(thermostat.HeatingSetpoint),
		CoolingSetpoint:   float64(thermostat.CoolingSetpoint),
		Override:          thermostat.Schedule != "On",
		Fan:               fanState,
		Mode:              modeNameMappings[thermostat.System],
		State:             thermState,
		EnabledHeatStages: thermostat.HeatStages,
		EnabledCoolStages: thermostat.CoolStages,
		Time:              time.Now().UnixNano(),
	}, nil
}

func (pel *Pelican) ModifySetpoints(params *PelicanSetpointParams) error {
	var value string
	// heating setpoint
	if params.HeatingSetpoint != nil {
		value += fmt.Sprintf("heatSetting:%d;", int(*params.HeatingSetpoint))
	}
	// cooling setpoint
	if params.CoolingSetpoint != nil {
		value += fmt.Sprintf("coolSetting:%d;", int(*params.CoolingSetpoint))
	}
	resp, _, errs := gorequest.New().Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "set").
		Param("object", "thermostat").
		Param("selection", fmt.Sprintf("name:%s;", pel.Name)).
		Param("value", value).
		End()
	if errs != nil {
		return fmt.Errorf("Error modifying thermostat temp settings: %v", errs)
	}

	defer resp.Body.Close()
	var result apiResult
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return fmt.Errorf("Failed to decode response XML: %v", err)
	}
	if result.Success == 0 {
		return fmt.Errorf("Error modifying thermostat temp settings: %v", result.Message)
	}

	return nil
}

func (pel *Pelican) ModifyState(params *PelicanStateParams) error {
	var value string

	// mode
	if params.Mode != nil {
		mode := int(*params.Mode)
		if mode < 0 || mode > 3 {
			return fmt.Errorf("Specified thermostat mode %d is invalid", mode)
		}
		systemVal := modeValMappings[mode]
		value += fmt.Sprintf("system:%s;", systemVal)
	}

	// override
	if params.Override != nil {
		var scheduleVal string
		if *params.Override == 1 {
			scheduleVal = "Off"
		} else {
			scheduleVal = "On"
		}
		value += fmt.Sprintf("schedule:%s;", scheduleVal)
	}

	// fan
	if params.Fan != nil {
		var fanVal string
		if *params.Fan == 1 {
			fanVal = "On"
		} else {
			fanVal = "Auto"
		}
		value += fmt.Sprintf("fan:%s;", fanVal)
	}

	// heating setpoint
	if params.HeatingSetpoint != nil {
		value += fmt.Sprintf("heatSetting:%d;", int(*params.HeatingSetpoint))
	}
	// cooling setpoint
	if params.CoolingSetpoint != nil {
		value += fmt.Sprintf("coolSetting:%d;", int(*params.CoolingSetpoint))
	}

	resp, _, errs := gorequest.New().Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "set").
		Param("object", "thermostat").
		Param("selection", fmt.Sprintf("name:%s;", pel.Name)).
		Param("value", value).
		End()
	if errs != nil {
		return fmt.Errorf("Error modifying thermostat state: %v (%s)", errs, resp.Request.URL)
	}

	defer resp.Body.Close()
	var result apiResult
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return fmt.Errorf("Failed to decode response XML: %v", err)
	}
	if result.Success == 0 {
		return fmt.Errorf("Error modifying thermostat state: %s", result.Message)
	}

	// rewrite default heat/cool stage info
	if result.Thermostat.HeatStages == 0 {
		result.Thermostat.HeatStages = 1
	}
	if result.Thermostat.CoolStages == 0 {
		result.Thermostat.CoolStages = 1
	}

	return nil
}

func (pel *Pelican) ModifyStages(params *PelicanStageParams) error {
	// Turn the thermostat off, saving its previously active mode
	status, err := pel.GetStatus()
	if err != nil {
		return fmt.Errorf("Error retrieving thermostat status: %s", err)
	}
	newMode := float64(0) // Off
	if err := pel.ModifyState(&PelicanStateParams{Mode: &newMode}); err != nil {
		return fmt.Errorf("Failed to turn thermostat off: %s", err)
	}

	// Restore the thermostat to its previous mode
	defer func() {
		if status != nil {
			oldMode := float64(status.Mode)
			if err := pel.ModifyState(&PelicanStateParams{Mode: &oldMode}); err != nil {
				fmt.Printf("Failed to restore thermostat to old mode: %s\n", err)
			}
		}
	}()

	// Change the thermostat's stage configuration
	var value string
	if params.HeatingStages != nil {
		value += fmt.Sprintf("heatStages:%d;", *params.HeatingStages)
	}
	if params.CoolingStages != nil {
		value += fmt.Sprintf("coolStages:%d;", *params.CoolingStages)
	}

	resp, _, errs := gorequest.New().Get(pel.target).
		Param("username", pel.username).
		Param("password", pel.password).
		Param("request", "set").
		Param("object", "thermostat").
		Param("selection", fmt.Sprintf("name:%s;", pel.Name)).
		Param("value", value).
		End()
	if errs != nil {
		return fmt.Errorf("Error modifying thermostat stages: %v (%s)", errs, resp.Request.URL)
	}

	defer resp.Body.Close()
	var result apiResult
	dec := xml.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return fmt.Errorf("Failed to decode response XML: %v", err)
	}
	if result.Success == 0 {
		return fmt.Errorf("Error modifying thermostat state: %s", result.Message)
	}

	// rewrite default heat/cool stage info
	if result.Thermostat.HeatStages == 0 {
		result.Thermostat.HeatStages = 1
	}
	if result.Thermostat.CoolStages == 0 {
		result.Thermostat.CoolStages = 1
	}
	return nil
}
