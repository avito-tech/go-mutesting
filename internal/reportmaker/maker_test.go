package reportmaker

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/avito-tech/go-mutesting/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeHTMLReport(t *testing.T) {
	tests := []struct {
		name    string
		report  models.Report
		wantErr bool
	}{
		{
			name: "successful report creation",
			report: models.Report{
				Stats: models.Stats{
					Msi: 85.5,
				},
				Escaped: []models.Mutant{
					{
						Mutator: models.Mutator{
							OriginalFilePath: "test.go",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty report",
			report: models.Report{
				Stats:   models.Stats{},
				Escaped: []models.Mutant{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MakeHTMLReport(tt.report)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := os.Stat(models.ReportHTMLFileName)
				assert.NoError(t, err)
				defer func(name string) {
					err := os.Remove(name)
					if err != nil {

					}
				}(models.ReportHTMLFileName)
			}
		})
	}
}

func TestMakeJSONReport(t *testing.T) {
	tests := []struct {
		name    string
		report  models.Report
		wantErr bool
	}{
		{
			name: "successful json report creation",
			report: models.Report{
				Stats: models.Stats{
					Msi: 90.0,
				},
				Escaped: []models.Mutant{
					{
						Mutator: models.Mutator{
							OriginalFilePath: "main.go",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty json report",
			report: models.Report{
				Stats:   models.Stats{},
				Escaped: []models.Mutant{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MakeJSONReport(tt.report)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := os.Stat(models.ReportFileName)
				assert.NoError(t, err)

				content, err := os.ReadFile(models.ReportFileName)
				require.NoError(t, err)

				var parsedReport models.Report
				err = json.Unmarshal(content, &parsedReport)
				require.NoError(t, err)

				assert.Equal(t, tt.report.Stats.Msi, parsedReport.Stats.Msi)

				defer func(name string) {
					err := os.Remove(name)
					if err != nil {

					}
				}(models.ReportFileName)
			}
		})
	}
}

func TestCreateOrTruncateReportFile(t *testing.T) {
	filename := "test_report.tmp"

	file, err := createOrTruncateReportFile(filename)
	require.NoError(t, err)
	require.NotNil(t, file)

	_, err = file.WriteString("test content")
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	file, err = createOrTruncateReportFile(filename)
	require.NoError(t, err)
	require.NotNil(t, file)

	_, err = os.ReadFile(filename)
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)
	err = os.Remove(filename)
	require.NoError(t, err)
}

func TestCloseReportFile(t *testing.T) {
	file, err := os.CreateTemp("", "test_close")
	require.NoError(t, err)

	filename := file.Name()

	closeReportFile(file, filename)

	_, err = os.Stat(filename)
	assert.NoError(t, err)

	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestFuncMap(t *testing.T) {
	diff := "line1\nline2\nline3"
	lines := funcMap["splitDiff"].(func(string) []string)(diff)
	expected := []string{"line1", "line2", "line3"}
	assert.Equal(t, expected, lines)

	result := funcMap["hasPrefix"].(func(string, string) bool)("test_string", "test")
	assert.True(t, result)

	result = funcMap["hasPrefix"].(func(string, string) bool)("test_string", "string")
	assert.False(t, result)
}
