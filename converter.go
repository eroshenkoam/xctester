package main

import (
	"github.com/eroshennkoam/xcresults/allure"
	"github.com/tidwall/gjson"
)

func convertSummary(summary gjson.Result) (_ allure.TestResult) {
	testResult := allure.TestResult{}
	testResult.Name = getValue(summary, "name._value")
	testResult.FullName = getValue(summary, "identifier._value")
	testResult.Status = getStatus(summary, "testStatus._value")
	return testResult
}

func getValue(summary gjson.Result, key string) string {
	value := summary.Get(key)
	if value.Exists() {
		return value.Str
	} else {
		return ""
	}
}

func getStatus(summary gjson.Result, key string) string {
	value := summary.Get(key)
	if value.Exists() {
		if "Success" == value.Str {
			return "passed"
		}
		if "Failure" == value.Str {
			return "failed"
		}
	}
	return "undefined"
}
