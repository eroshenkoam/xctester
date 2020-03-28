package allure

const (
	Skipped Status = "skipped"
	Passed  Status = "passed"
	Failed  Status = "failed"
	Broken  Status = "broken"
	Unknown Status = "unknown"
)

type TestResult struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	HistoryId string `json:"historyId,omitempty"`
	Status    Status `json:"status"`

	Steps []StepResult `json:"steps"`

	Start int `json:"start"`
	Stop  int `json:"stop"`

	Labels      []Label      `json:"labels"`
	Parameters  []Parameter  `json:"parameters"`
	Attachments []Attachment `json:"attachments"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attachment struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	MimeType string `json:"type"`
}

type StepResult struct {
	Name   string `json:"name"`
	Status Status `json:"status"`

	Start int `json:"start"`
	Stop  int `json:"stop"`

	Steps       []StepResult `json:"steps"`
	Parameters  []Parameter  `json:"parameters"`
	Attachments []Attachment `json:"attachments"`
}

type Status string

type StatusDetails struct {
	Message string `json:"message"`
	Trace   string `json:"trace,omitempty"`
}
