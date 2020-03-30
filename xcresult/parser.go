package xcresult

import (
	"encoding/json"
	"github.com/eroshenkoam/xctester/allure"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Attachment struct {
	name string
	ref  string
}

func Export(input string, output string) {
	testRefs := make(chan string)
	go extractTestRefs(input, testRefs)

	summaryRefs := make(chan string)
	go extractSummaryRefs(input, testRefs, summaryRefs)

	attachments := make(chan Attachment)
	go exportAttachments(input, output, attachments)

	results := make(chan allure.TestResult)
	go exportResults(output, results)

	exportSummaryRefs(input, summaryRefs, results, attachments)
}

func extractTestRefs(path string, testRefIds chan string) {
	resultSummary := readSummary(path)
	refs := resultSummary.Get("actions._values.#.actionResult.testsRef.id._value").Array()
	for _, ref := range refs {
		testRefIds <- ref.Str
	}
	close(testRefIds)
}

func extractSummaryRefs(input string, testRefIds chan string, summaryRefIds chan string) {
	for ref := range testRefIds {
		testRef := readReference(input, ref)
		for _, summary := range testRef.Get("summaries._values").Array() {
			for _, testableSummary := range summary.Get("testableSummaries._values").Array() {
				for _, test := range testableSummary.Get("tests._values").Array() {
					collectTestSummaryRefs(test, summaryRefIds)
				}
			}
		}
	}
	close(summaryRefIds)
}

func collectTestSummaryRefs(test gjson.Result, summaryRefIds chan string) {
	if test.Get("summaryRef.id._value").Exists() {
		summaryRefIds <- test.Get("summaryRef.id._value").String()
	}
	for _, subtest := range test.Get("subtests._values").Array() {
		collectTestSummaryRefs(subtest, summaryRefIds)
	}
}

func extractAttachments(activity gjson.Result, attachments chan Attachment) {
	for _, child := range activity.Get("subactivities._values").Array() {
		extractAttachments(child, attachments)
	}
	for _, attachment := range activity.Get("attachments._values").Array() {
		if attachment.Get("payloadRef").Exists() {
			name := attachment.Get("filename._value").Str
			ref := attachment.Get("payloadRef.id._value").Str
			attachments <- Attachment{name: name, ref: ref}
		}
	}
}

func exportSummaryRefs(path string, refs chan string, results chan allure.TestResult, attachments chan Attachment) {
	for summaryRef := range refs {
		summary := readReference(path, summaryRef)
		results <- convertSummary(summary)
		for _, activitySummary := range summary.Get("activitySummaries._values").Array() {
			extractAttachments(activitySummary, attachments)
		}
	}
	close(attachments)
	close(results)
}

func exportAttachments(path string, output string, attachments chan Attachment) {
	for attachment := range attachments {
		exportReference(path, attachment.ref, filepath.Join(output, attachment.name))
	}
}

func exportResults(output string, results chan allure.TestResult) {
	for result := range results {
		resultJson, _ := json.Marshal(result)
		resultFile := filepath.Join(output, uuid.New().String()+"-result.json")
		if err := ioutil.WriteFile(resultFile, resultJson, 0744); err != nil {
			log.Fatalln("Can not create result file [", resultFile, "] because ", err)
		}
	}
}
