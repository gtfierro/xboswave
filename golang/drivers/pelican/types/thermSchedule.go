package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	rrule "github.com/teambition/rrule-go"
)

// Login, Authentication, Thermostat ID Retrieval Structs
type thermIDRequest struct {
	Resources []thermIDResources `json:"resources"`
}

type thermIDResources struct {
	Children    []thermIDChild `json:"children"`
	GroupId     string         `json:"groupId"`
	Permissions string         `json:"permissions"`
}

type thermIDChild struct {
	Id          string `json:"id"`
	Permissions string `json:"permissions"`
}

// Thermostat Settings Structs
type settingsRequest struct {
	Epnum    float64         `json:"epnum"`
	Id       string          `json:"id"`
	Nodename string          `json:"nodename"`
	Userdata settingsWrapper `json:"userdata"`
}

type settingsWrapper struct {
	Epnum    float64 `json:"epnum"`
	Fan      string  `json:"fan"`
	Nodename string  `json:"nodename"`
	Repeat   string  `json:"repeat"`
}

// Thermostat Schedule By Day Decoding Structs
type scheduleRequest struct {
	ClientData scheduleSetTimes `json:"clientdata"`
}

type scheduleSetTimes struct {
	SetTimes []scheduleTimeBlock `json:"setTimes"`
}

type scheduleTimeBlock struct {
	HeatSetting float64 `json:"heatSetting"`
	CoolSetting float64 `json:"coolSetting"`
	StartValue  string  `json:"startValue"`
	System      string  `json:"systemDisplay"`
}

// Thermostat Schedule Structs

// Struct mapping each day of the week to its daily schedule
type ThermostatSchedule struct {
	DaySchedules map[string]([]ThermostatBlockSchedule) `msgpack:"day_schedules"`
}

// Struct containing data defining the settings of each schedule block
type ThermostatBlockSchedule struct {
	CoolSetting float64 `msgpack:"cool_setting"`
	HeatSetting float64 `msgpack:"heat_setting"`
	System      string  `msgpack:"system"`
	Time        string  `msgpack:"time"`
}

// Time Constant for Cookie Refresh. 720 Hours = 30 Days
const cookieDuration = 720 * time.Hour

var week = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var weekRRule = [...]rrule.Weekday{rrule.SU, rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR, rrule.SA}

func (pel *Pelican) GetSchedule() (*ThermostatSchedule, error) {
	// Check if cookie needs to be refreshed
	if time.Since(pel.cookieTime) > cookieDuration {
		if cookieAndIDError := pel.setCookieAndID(); cookieAndIDError != nil {
			return nil, cookieAndIDError
		}
	}

	thermSchedule := ThermostatSchedule{
		DaySchedules: make(map[string]([]ThermostatBlockSchedule), len(week)),
	}

	// Retrieve Repeat Type (Daily, Weekly, Weekend/Weekday) and Nodename from Thermostat's Settings
	settings, settingsErr := pel.getSettings()
	if settingsErr != nil {
		return nil, fmt.Errorf("Failed to determine repeat type for thermostat %v: %v", pel.Name, settingsErr)
	}
	repeatType := settings.Repeat
	nodename := settings.Nodename
	epnum := settings.Epnum

	// Build Schedule by Repeat Type
	if repeatType == "Daily" {
		schedule, scheduleError := pel.getScheduleByDay(0, epnum, nodename)
		if scheduleError != nil {
			return nil, fmt.Errorf("Error retrieving schedule for thermostat %v: %v", nodename, scheduleError)
		}
		for _, day := range week {
			thermSchedule.DaySchedules[day] = *schedule
		}
	} else if repeatType == "Weekly" {
		for index, day := range week {
			schedule, scheduleError := pel.getScheduleByDay(index, epnum, nodename)
			if scheduleError != nil {
				return nil, fmt.Errorf("Error retrieving schedule for thermostat %v on %v (day %v): %v", nodename, day, index, scheduleError)
			}
			thermSchedule.DaySchedules[day] = *schedule
		}
	} else if repeatType == "Weekday/Weekend" {
		weekend, weekendError := pel.getScheduleByDay(0, epnum, nodename)
		if weekendError != nil {
			return nil, fmt.Errorf("Error retrieving schedule for thermostat %v on weekend (day 0): %v", nodename, weekendError)
		}
		for _, day := range []string{"Sunday", "Saturday"} {
			thermSchedule.DaySchedules[day] = *weekend
		}
		weekday, weekdayError := pel.getScheduleByDay(1, epnum, nodename)
		if weekdayError != nil {
			return nil, fmt.Errorf("Error retrieving schedule for thermostat %v on weekday (day 1): %v", nodename, weekdayError)
		}
		for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"} {
			thermSchedule.DaySchedules[day] = *weekday
		}
	} else {
		return nil, fmt.Errorf("Failed to recognize repeat type of thermostat %v's schedule: %v", nodename, repeatType)
	}

	return &thermSchedule, nil
}

func (pel *Pelican) getSettings() (*settingsWrapper, error) {
	var requestURL bytes.Buffer
	requestURL.WriteString(fmt.Sprintf("https://%s.officeclimatecontrol.net/ajaxThermostat.cgi?id=", pel.sitename))
	requestURL.WriteString(pel.id)
	requestURL.WriteString(":Thermostat&request=GetSchedule")

	resp, _, errs := pel.scheduleReq.Get(requestURL.String()).Type("form").AddCookie(pel.cookie).End()
	if errs != nil {
		return nil, fmt.Errorf("Failed to retrieve schedule settings for thermostat %v: %v", pel.id, errs)
	}
	var result settingsRequest
	decoder := json.NewDecoder(resp.Body)
	if decodeError := decoder.Decode(&result); decodeError != nil {
		return nil, fmt.Errorf("Failed to decode schedule settings for thermostat %v: %v", pel.id, decodeError)
	}
	return &result.Userdata, nil
}

func (pel *Pelican) getScheduleByDay(dayOfWeek int, epnum float64, thermostatID string) (*[]ThermostatBlockSchedule, error) {
	// Construct Request URL for Thermostat Schedule by Day of Week
	var requestURL bytes.Buffer
	requestURL.WriteString(fmt.Sprintf("https://%s.officeclimatecontrol.net/thermDayEdit.cgi?section=json&nodename=", pel.sitename))
	requestURL.WriteString(thermostatID)
	requestURL.WriteString("&epnum=")
	requestURL.WriteString(fmt.Sprintf("%.0f", epnum))
	requestURL.WriteString("&dayofweek=")
	requestURL.WriteString(strconv.Itoa(dayOfWeek))

	// Make Request, Decode into Response Struct
	resp, _, errs := pel.scheduleReq.Get(requestURL.String()).Type("form").AddCookie(pel.cookie).End()
	if errs != nil {
		return nil, fmt.Errorf("Failed to retrieve schedule for thermostat %v on day of week %v: %v", thermostatID, dayOfWeek, errs)
	}
	var result scheduleRequest
	decoder := json.NewDecoder(resp.Body)
	if decodeError := decoder.Decode(&result); decodeError != nil {
		return nil, fmt.Errorf("Failed to decode schedule for thermostat %v on day of week %v: %v", thermostatID, dayOfWeek, decodeError)
	}

	// Transfer Response Struct Data into return struct
	var daySchedule []ThermostatBlockSchedule
	for _, block := range result.ClientData.SetTimes {
		returnBlock := ThermostatBlockSchedule{
			CoolSetting: block.CoolSetting,
			HeatSetting: block.HeatSetting,
			System:      block.System,
		}

		if rruleTime, rruleError := convertTimeToRRule(dayOfWeek, block.StartValue, pel.timezone); rruleError != nil {
			return nil, fmt.Errorf("Failed to convert time in string format %v to rrule format: %v", block.StartValue, rruleError)
		} else {
			returnBlock.Time = rruleTime
		}
		daySchedule = append(daySchedule, returnBlock)
	}
	return &daySchedule, nil
}

func convertTimeToRRule(dayOfWeek int, blockTime string, timezone *time.Location) (string, error) {
	timeParsed, timeParsedErr := time.Parse("15:04:PM", blockTime)
	if timeParsedErr != nil {
		return "", fmt.Errorf("Error parsing time %v with time.Parse function: %v\n: ", blockTime, timeParsedErr)
	}

	rruleSched, rruleSchedErr := rrule.NewRRule(rrule.ROption{
		Freq:    rrule.WEEKLY,
		Wkst:    weekRRule[dayOfWeek],
		Dtstart: time.Date(0, 0, 0, timeParsed.Hour(), timeParsed.Minute(), 0, 0, timezone),
	})
	if rruleSchedErr != nil {
		return "", fmt.Errorf("Error creating rruleSchedule object: %v\n", rruleSchedErr)
	}

	return rruleSched.String(), nil
}

// Handles setting the Pelican fields (cookie, id) that can only be retrieved by AJAX Requests
func (pel *Pelican) setCookieAndID() error {
	// Set the cookie field using the pelican's login information
	loginInfo := map[string]interface{}{
		"username": pel.username,
		"password": pel.password,
		"sitename": pel.sitename,
	}
	respLogin, _, errsLogin := pel.scheduleReq.Post(fmt.Sprintf("https://%s.officeclimatecontrol.net/#_loginPage", pel.sitename)).Type("form").Send(loginInfo).End()
	if (errsLogin != nil) || (respLogin.StatusCode != 200) {
		return fmt.Errorf("Error logging into climate control website to retrieve cookie: %v", errsLogin)
	}
	pel.cookie = (*http.Response)(respLogin).Cookies()[0]
	pel.cookieTime = time.Now()

	// Set the Thermostat ID using the thermostat resources AJAX Request
	respTherms, _, errsTherms := pel.scheduleReq.Get(fmt.Sprintf("https://%s.officeclimatecontrol.net/ajaxSchedule.cgi?request=getResourcesExtended&resourceType=Thermostats", pel.sitename)).Type("form").AddCookie(pel.cookie).End()
	if (errsTherms != nil) || (respTherms.StatusCode != 200) {
		return fmt.Errorf("Error retrieving Thermostat IDs: %v", errsTherms)
	}
	var IDRequest thermIDRequest
	decoder := json.NewDecoder(respTherms.Body)
	if decodeError := decoder.Decode(&IDRequest); decodeError != nil {
		return fmt.Errorf("Failed to decode Thermostat ID response JSON: %v\n", decodeError)
	}
	pel.id = IDRequest.Resources[0].Children[0].Id
	return nil
}
