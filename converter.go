package main

import (
	"fmt"
	"github.com/eroshennkoam/xcresults/allure"
	"github.com/tidwall/gjson"
	"regexp"
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
	Attachments   string = "attachments._values"
	SubActivities string = "subactivities._values"
)

const (
	Filename string = "filename._value"
	Payload  string = "payloadRef.id._value"
)

const (
	Success string = "Success"
	Failure string = "Failure"
)

type Status struct {
	Name    allure.Status
	Details allure.StatusDetails
}

func convertSummary(summary gjson.Result) (result allure.TestResult) {
	result = allure.TestResult{}
	result.Name = value(summary, Name)
	result.FullName = value(summary, Identifier)
	result.Status = status(summary)

	status, labels, steps := parseSubActivities(summary.Get(ActivitySummaries).Array())

	result.Steps = steps
	result.Labels = labels
	result.Status = status.Name
	result.StatusDetails = status.Details

	result.Start = steps[0].Start
	result.Stop = steps[len(steps)-1].Stop

	return result
}

func parseSubActivities(activities []gjson.Result) (status Status, labels []allure.Label, steps []allure.StepResult) {
	steps = []allure.StepResult{}
	labels = []allure.Label{}
	status = Status{Name: allure.Passed}
	for _, activity := range activities {
		title := activity.Get(Title).Str

		labelMatcher := regexp.MustCompile(`allure\.label\.(?P<name>.*):(?P<value>.*)`)
		if labelMatcher.MatchString(title) {
			for _, group := range labelMatcher.FindAllStringSubmatch(title, -1) {
				labels = append(labels, allure.Label{Name: group[1], Value: group[2]})
			}
			continue
		}

		step := allure.StepResult{Name: title, Status: allure.Passed}
		if activity.Get(Start).Exists() && activity.Get(Finish).Exists() {
			step.Start = date(activity.Get(Start).Str)
			step.Stop = date(activity.Get(Finish).Str)
		}
		if strings.HasPrefix(title, "Assertion Failure:") {
			status = Status{Name: allure.Failed, Details: allure.StatusDetails{Message: title}}
			step.StatusDetails = status.Details
			step.Status = status.Name
		}
		substatus, sublabels, substeps := parseSubActivities(activity.Get(SubActivities).Array())
		if substatus.Name != allure.Passed {
			status = Status{Name: substatus.Name, Details: substatus.Details}
			step.Status = status.Name
		}
		step.Attachments = attachments(activity.Get(Attachments).Array())
		step.Steps = substeps

		steps = append(steps, step)
		labels = append(labels, sublabels...)
	}
	return status, labels, steps
}

func attachments(nodes []gjson.Result) (attachments []allure.Attachment) {
	attachments = []allure.Attachment{}
	for _, node := range nodes {
		attachment := allure.Attachment{Name: node.Get(Filename).Str, Source: node.Get(Filename).Str}
		attachments = append(attachments, attachment)
	}
	return attachments
}

func date(value string) int64 {
	layout := "2006-01-02T15:04:05.000+0300"
	t, err := time.Parse(layout, value)
	if err != nil {
		fmt.Println(err)
	}
	return t.UnixNano() / int64(time.Millisecond)
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
