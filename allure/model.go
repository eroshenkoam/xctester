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

	Status        Status        `json:"status"`
	StatusDetails StatusDetails `json:"statusDetails,omitempty"`

	Steps []StepResult `json:"steps"`

	Start int64 `json:"start"`
	Stop  int64 `json:"stop"`

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
	MimeType string `json:"type,omitempty"`
}

type StepResult struct {
	Name string `json:"name"`

	Status        Status        `json:"status"`
	StatusDetails StatusDetails `json:"statusDetails,omitempty"`

	Start int64 `json:"start"`
	Stop  int64 `json:"stop"`

	Steps       []StepResult `json:"steps,omitempty"`
	Parameters  []Parameter  `json:"parameters,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Status string

type StatusDetails struct {
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}
