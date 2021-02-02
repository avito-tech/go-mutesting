package models

var ReportFileName string = "report.json"

// Report Structure for mutation report
type Report struct {
	Stats     Stats    `json:"stats"`
	Escaped   []Mutant `json:"escaped"`
	Timeouted []Mutant `json:"timeouted"`
	Killed    []Mutant `json:"killed"`
	Errored   []Mutant `json:"errored"`
}

type Stats struct {
	TotalMutantsCount    int64   `json:"totalMutantsCount"`
	KilledCount          int64   `json:"killedCount"`
	NotCoveredCount      int64   `json:"notCoveredCount"`
	EscapedCount         int64   `json:"escapedCount"`
	ErrorCount           int64   `json:"errorCount"`
	SkippedCount         int64   `json:"skippedCount"`
	TimeOutCount         int64   `json:"timeOutCount"`
	Msi                  float64 `json:"msi"`
	MutationCodeCoverage int64   `json:"mutationCodeCoverage"`
	CoveredCodeMsi       float64 `json:"coveredCodeMsi"`
	DuplicatedCount      int64   `json:"-"`
}

type Mutant struct {
	Mutator       Mutator `json:"mutator"`
	Diff          string  `json:"diff"`
	ProcessOutput string  `json:"processOutput,omitempty"`
}

type Mutator struct {
	MutatorName        string `json:"mutatorName"`
	OriginalSourceCode string `json:"originalSourceCode"`
	MutatedSourceCode  string `json:"mutatedSourceCode"`
	OriginalFilePath   string `json:"originalFilePath"`
	OriginalStartLine  int64  `json:"originalStartLine"`
}

func (report *Report) Calculate() {
	report.Stats.Msi = report.MsiScore()
	report.Stats.TotalMutantsCount = report.TotalCount()
}

func (report *Report) MsiScore() float64 {
	total := report.TotalCount()

	if total == 0 {
		return 0.0
	}

	return float64(report.Stats.KilledCount+report.Stats.ErrorCount+report.Stats.SkippedCount) / float64(total)
}

func (report *Report) TotalCount() int64 {
	return report.Stats.KilledCount + report.Stats.EscapedCount + report.Stats.ErrorCount + report.Stats.SkippedCount
}
