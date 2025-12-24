package report_maker

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"

	"github.com/avito-tech/go-mutesting/internal/models"
)

//go:embed templates/report.html.gotpl
var reportTmpl string

var funcMap = template.FuncMap{
	"splitDiff": func(diff string) []string {
		return strings.Split(diff, "\n")
	},
	"hasPrefix": strings.HasPrefix,
}

// MakeHtmlReport is a function for creating an HTML report based on a stripped-down version of the models.Report model (not all fields are used)
func MakeHtmlReport(report models.Report) error {
	report.Stats.Msi = math.Round(report.Stats.Msi*10_000) / 100

	groupedMutants := make(map[string][]models.Mutant)
	for _, mutant := range report.Escaped {
		filePath := mutant.Mutator.OriginalFilePath
		groupedMutants[filePath] = append(groupedMutants[filePath], mutant)
	}

	t, err := template.New(models.ReportHtmlFileName).Funcs(funcMap).Parse(reportTmpl)
	if err != nil {
		return fmt.Errorf("Error while parse template: %w ", err)
	}

	file, err := createOrTruncateReportFile(models.ReportHtmlFileName)
	if err != nil {
		return fmt.Errorf("Error while open/create .html report file from template: %w ", err)
	}
	defer closeReportFile(file, models.ReportHtmlFileName)

	data := struct {
		Stats          models.Stats
		GroupedMutants map[string][]models.Mutant
	}{
		Stats:          report.Stats,
		GroupedMutants: groupedMutants,
	}

	err = t.Execute(file, data)
	if err != nil {
		return fmt.Errorf("Error while execute template for .html report: %w ", err)
	}

	return nil
}

// MakeJsonReport is a function for creating json report, which is based on models.Report
func MakeJsonReport(report models.Report) error {
	jsonContent, err := json.Marshal(report)
	if err != nil {
		return err
	}

	file, err := createOrTruncateReportFile(models.ReportFileName)
	if err != nil {
		return fmt.Errorf("Error while open/create .json report file from template: %w ", err)
	}
	defer closeReportFile(file, models.ReportFileName)

	if file == nil {
		return errors.New("cannot create file for .json report")
	}

	_, err = file.WriteString(string(jsonContent))
	if err != nil {
		return err
	}

	return nil
}

func createOrTruncateReportFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func closeReportFile(file *os.File, filename string) {
	if err := file.Close(); err != nil {
		fmt.Printf("Error while closing %s: %v\n", filename, err)
	}
}
