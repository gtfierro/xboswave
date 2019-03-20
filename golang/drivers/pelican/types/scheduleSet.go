package types

import (
	"encoding/xml"
	"fmt"

	"github.com/parnurzeal/gorequest"
	rrule "github.com/teambition/rrule-go"
)

func (pel *Pelican) SetSchedule(newSchedule *ThermostatSchedule) error {
	// Retrieve thermostat's latest schedule for extra block deletion purposes
	currentSchedule, currentErr := pel.GetSchedule()
	if currentErr != nil {
		return fmt.Errorf("Error retrieving thermostat %v's current schedule: %v", pel.Name, currentErr)
	}

	for _, day := range week {
		if blocks, ok := newSchedule.DaySchedules[day]; ok {
			// Delete Unnecessary Blocks + Error Checking
			currentBlockCount := len(currentSchedule.DaySchedules[day])
			requiredBlockCount := len(blocks)
			for requiredBlockCount < currentBlockCount {
				respDelete, _, errsDelete := gorequest.New().Get(pel.target).
					Param("username", pel.username).
					Param("password", pel.password).
					Param("request", "set").
					Param("object", "thermostatSchedule").
					Param("selection", fmt.Sprintf("name:%s;dayOfWeek:%s;setTime:%v;", pel.Name, day, currentBlockCount)).
					Param("value", "delete").
					End()

				if errsDelete != nil {
					return fmt.Errorf("Error deleting thermostat schedule settings on day %v: %v\n", day, errsDelete)
				}
				defer respDelete.Body.Close()
				var resultDelete apiResult
				decDelete := xml.NewDecoder(respDelete.Body)
				if errDecodeDelete := decDelete.Decode(&resultDelete); errDecodeDelete != nil {
					return fmt.Errorf("Failed to decode thermostat schedule delete response XML: %v", errDecodeDelete)
				}
				if resultDelete.Success == 0 {
					return fmt.Errorf("Error deleting thermostat schedule settings on day %v: %v\n", day, resultDelete.Message)
				}
				currentBlockCount -= 1
			}

			// Construct and Set New Schedule Settings by Block
			for index, block := range blocks {
				var value string = ""
				value += fmt.Sprintf("coolSetting:%.0f;", block.CoolSetting)
				value += fmt.Sprintf("heatSetting:%.0f", block.HeatSetting)
				value += fmt.Sprintf("system:%s", block.System)

				// Convert Time to Pelican's Timezone
				timeRRule, timeRRuleErr := rrule.StrToRRule(block.Time)
				if timeRRuleErr != nil {
					return fmt.Errorf("Error converting time string %v to RRule format: %v\n", block.Time, timeRRuleErr)
				}
				timeLocal := timeRRule.OrigOptions.Dtstart.In(pel.timezone)
				value += fmt.Sprintf("startTime:%s;", timeLocal.Format("03:04"))

				// Set Request + Error Checking
				respSet, _, errsSet := gorequest.New().Get(pel.target).
					Param("username", pel.username).
					Param("password", pel.password).
					Param("request", "set").
					Param("object", "thermostatSchedule").
					Param("selection", fmt.Sprintf("name:%s;dayOfWeek:%s;setTime:%v;", pel.Name, day, index+1)).
					Param("value", value).
					End()

				if errsSet != nil {
					return fmt.Errorf("Error setting thermostat schedule on day %v: %v\n", day, errsSet)
				}
				defer respSet.Body.Close()
				var resultSet apiResult
				decSet := xml.NewDecoder(respSet.Body)
				if errDecodeSet := decSet.Decode(&resultSet); errDecodeSet != nil {
					return fmt.Errorf("Failed to decode thermostat schedule set response XML: %v", errDecodeSet)
				}
				if resultSet.Success == 0 {
					return fmt.Errorf("Error setting thermostat schedule on day %v: %v\n", day, resultSet.Message)
				}
			}
		} else {
			// Delete Entire Day's Schedule
			blockCount := len(currentSchedule.DaySchedules[day])
			for i := blockCount; i >= 1; i-- {
				respDelete, _, errsDelete := gorequest.New().Get(pel.target).
					Param("username", pel.username).
					Param("password", pel.password).
					Param("request", "set").
					Param("object", "thermostatSchedule").
					Param("selection", fmt.Sprintf("name:%s;dayOfWeek:%s;setTime:%v;", pel.Name, day, i)).
					Param("value", "delete").
					End()

				if errsDelete != nil {
					return fmt.Errorf("Error deleting thermostat schedule settings on day %v: %v\n", day, errsDelete)
				}
				defer respDelete.Body.Close()
				var resultDelete apiResult
				decDelete := xml.NewDecoder(respDelete.Body)
				if errDecodeDelete := decDelete.Decode(&resultDelete); errDecodeDelete != nil {
					return fmt.Errorf("Failed to decode thermostat schedule delete response XML: %v", errDecodeDelete)
				}
				if resultDelete.Success == 0 {
					return fmt.Errorf("Error deleting thermostat schedule settings on day %v: %v\n", day, resultDelete.Message)
				}
			}
		}
	}

	return nil
}
