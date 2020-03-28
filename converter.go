package main

import (
	"fmt"
	"github.com/eroshennkoam/xcresults/allure"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

const (
	Name              string = "name._value"
	Identifier        string = "identifier._value"
	TestStatus        string = "testStatus._value"
	ActivitySummaries string = "activitySummaries._values"
)

const (
	Title         string = "title._value"
	Start         string = "start._value"
	Finish        string = "finish._value"
	SubActivities string = "subactivities._values"
)

const (
	Success string = "Success"
	Failure string = "Failure"
)

type Status struct {
	Name    allure.Status
	Details allure.StatusDetails
}

func convertSummary(summary gjson.Result) (_ allure.TestResult) {
	testResult := allure.TestResult{}
	testResult.Name = value(summary, Name)
	testResult.FullName = value(summary, Identifier)
	testResult.Status = status(summary)

	status, steps := parseActivities(summary)
	testResult.Steps = steps
	testResult.Status = status.Name
	testResult.StatusDetails = status.Details

	return testResult
}

func parseActivities(node gjson.Result) (status Status, steps []allure.StepResult) {
	steps = []allure.StepResult{}
	status = Status{Name: allure.Passed}
	for _, subactibity := range node.Get(ActivitySummaries).Array() {
		title := subactibity.Get(Title).Str
		step := allure.StepResult{Name: title}

		substatus, substeps := parseSubActivities(subactibity)
		if substatus.Name == allure.Failed {
			status.Name = substatus.Name
			status.Details = substatus.Details
		}
		step.StatusDetails = substatus.Details
		step.Status = substatus.Name
		step.Steps = substeps

		steps = append(steps, step)
	}
	return status, steps
}

func parseSubActivities(node gjson.Result) (status Status, steps []allure.StepResult) {
	steps = []allure.StepResult{}
	status = Status{Name: allure.Passed}
	for _, subActibity := range node.Get(SubActivities).Array() {
		title := subActibity.Get(Title).Str
		if strings.HasPrefix(title, "Assertion Failure:") {
			status.Details = allure.StatusDetails{Message: title}
			status.Name = allure.Failed
			break
		}
		step := allure.StepResult{Name: title}
		if subActibity.Get(Start).Exists() && subActibity.Get(Finish).Exists() {
			step.Start = date(subActibity.Get(Start).Str)
			step.Start = date(subActibity.Get(Finish).Str)
		}

		substatus, substeps := parseSubActivities(subActibity)
		if substatus.Name == allure.Failed {
			status.Name = substatus.Name
			status.Details = substatus.Details
		}
		step.StatusDetails = substatus.Details
		step.Status = substatus.Name
		step.Steps = substeps

		steps = append(steps, step)
	}
	return status, steps
}

func date(value string) int64 {
	layout := "2006-01-02T15:04:05.000+0300"
	t, err := time.Parse(layout, value)
	if err != nil {
		fmt.Println(err)
	}
	return t.Unix()
}

func value(summary gjson.Result, key string) string {
	value := summary.Get(key)
	if value.Exists() {
		return value.Str
	} else {
		return ""
	}
}

func status(summary gjson.Result) allure.Status {
	value := summary.Get(TestStatus)
	if value.Exists() {
		if Success == value.Str {
			return allure.Passed
		}
		if Failure == value.Str {
			return allure.Failed
		}
	}
	return allure.Unknown
}
