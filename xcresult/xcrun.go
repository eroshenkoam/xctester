package xcresult

import (
	"bytes"
	"github.com/tidwall/gjson"
	"log"
	"os/exec"
)

func readSummary(path string) (data gjson.Result) {
	cmd := exec.Command("xcresult", "xcresulttool", "get",
		"--format", "json",
		"--path", path)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return gjson.Parse(out.String())
}

func readReference(path string, id string) (data gjson.Result) {
	cmd := exec.Command("xcresult", "xcresulttool", "get",
		"--format", "json",
		"--path", path,
		"--id", id)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return gjson.Parse(out.String())
}

func exportReference(path string, id string, output string) {
	cmd := exec.Command("xcresult", "xcresulttool", "export",
		"--type", "file",
		"--path", path,
		"--id", id,
		"--output-path", output)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
