package main

import (
	"github.com/eroshennkoam/xcresults/allure"
	"github.com/tidwall/gjson"
)

const (
	Name              string = "name._value"
	Identifier        string = "identifier._value"
	TestStatus        string = "testStatus._value"
	ActivitySummaries string = "activitySummaries._values"
)

const (
	Title         string = "title._value"
	SubActivities string = "subactivities._values"
)

const (
	Success string = "Success"
	Failure string = "Failure"
)

func convertSummary(summary gjson.Result) (_ allure.TestResult) {
	testResult := allure.TestResult{}
	testResult.Name = value(summary, Name)
	testResult.FullName = value(summary, Identifier)
	testResult.Status = status(summary)
	testResult.Steps = []allure.StepResult{}
	if summary.Get(ActivitySummaries).Exists() {
		for _, activity := range summary.Get(ActivitySummaries).Array() {
			step := parseActivities(activity)
			testResult.Steps = append(testResult.Steps, step)
		}
	}
	return testResult
}

func parseActivities(activity gjson.Result) (step allure.StepResult) {
	title := activity.Get(Title).Str
	step = allure.StepResult{Name: title, Status: allure.Passed}

	if activity.Get(SubActivities).Exists() {
		step.Steps = []allure.StepResult{}
		for _, subactibity := range activity.Get(SubActivities).Array() {
			substep := parseActivities(subactibity)
			step.Steps = append(step.Steps, substep)
		}
	}
	return step
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
