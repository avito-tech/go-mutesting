package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/avito-tech/go-mutesting/internal/models"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainSimple(t *testing.T) {
	testMain(
		t,
		"../../example",
		[]string{"--debug", "--exec-timeout", "1"},
		returnOk,
		"The mutation score is 0.551724 (16 passed, 13 failed, 8 duplicated, 0 skipped, total is 29)",
	)
}

func TestMainRecursive(t *testing.T) {
	testMain(
		t,
		"../../example",
		[]string{"--debug", "--exec-timeout", "1", "./..."},
		returnOk,
		"The mutation score is 0.580645 (18 passed, 13 failed, 8 duplicated, 0 skipped, total is 31)",
	)
}

func TestMainFromOtherDirectory(t *testing.T) {
	testMain(
		t,
		"../..",
		[]string{"--debug", "--exec-timeout", "1", "github.com/avito-tech/go-mutesting/example"},
		returnOk,
		"The mutation score is 0.551724 (16 passed, 13 failed, 8 duplicated, 0 skipped, total is 29)",
	)
}

func TestMainMatch(t *testing.T) {
	testMain(
		t,
		"../../example",
		[]string{"--debug", "--exec", "../scripts/exec/test-mutated-package.sh", "--exec-timeout", "1", "--match", "baz", "./..."},
		returnOk,
		"The mutation score is 0.500000 (2 passed, 2 failed, 0 duplicated, 0 skipped, total is 4)",
	)
}

func TestMainSkipWithoutTest(t *testing.T) {
	testMain(
		t,
		"../../example",
		[]string{"--debug", "--exec-timeout", "1", "--config", "../testdata/configs/configSkipWithoutTest.yml.test"},
		returnOk,
		"The mutation score is 0.592593 (16 passed, 11 failed, 8 duplicated, 0 skipped, total is 27)",
	)
}

func TestMainJSONReport(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "go-mutesting-main-test-")
	assert.NoError(t, err)

	reportFileName := "reportTestMainJSONReport.json"
	jsonFile := tmpDir + "/" + reportFileName
	if _, err := os.Stat(jsonFile); err == nil {
		err = os.Remove(jsonFile)
		assert.NoError(t, err)
	}

	models.ReportFileName = jsonFile

	testMain(
		t,
		"../../example",
		[]string{"--debug", "--exec-timeout", "1", "--config", "../testdata/configs/configForJson.yml.test"},
		returnOk,
		"The mutation score is 0.592593 (16 passed, 11 failed, 8 duplicated, 0 skipped, total is 27)",
	)

	info, err := os.Stat(jsonFile)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	defer func() {
		err = os.Remove(jsonFile)
		if err != nil {
			fmt.Println("Error while deleting temp file")
		}
	}()

	jsonData, err := ioutil.ReadFile(jsonFile)
	assert.NoError(t, err)

	var mutationReport models.Report
	err = json.Unmarshal(jsonData, &mutationReport)
	assert.NoError(t, err)

	expectedStats := models.Stats{
		TotalMutantsCount:    27,
		KilledCount:          16,
		NotCoveredCount:      0,
		EscapedCount:         11,
		ErrorCount:           0,
		SkippedCount:         0,
		TimeOutCount:         0,
		Msi:                  0.5925925925925926,
		MutationCodeCoverage: 0,
		CoveredCodeMsi:       0,
		DuplicatedCount:      0,
	}

	assert.Equal(t, expectedStats, mutationReport.Stats)
	assert.Equal(t, 11, len(mutationReport.Escaped))
	assert.Nil(t, mutationReport.Timeouted)
	assert.Equal(t, 16, len(mutationReport.Killed))
	assert.Nil(t, mutationReport.Errored)

	for i := 0; i < len(mutationReport.Escaped); i++ {
		assert.Contains(t, mutationReport.Escaped[i].ProcessOutput, "FAIL")
	}
	for i := 0; i < len(mutationReport.Killed); i++ {
		assert.Contains(t, mutationReport.Killed[i].ProcessOutput, "PASS")
	}
}

func testMain(t *testing.T, root string, exec []string, expectedExitCode int, contains string) {
	saveStderr := os.Stderr
	saveStdout := os.Stdout
	saveCwd, err := os.Getwd()
	assert.Nil(t, err)

	r, w, err := os.Pipe()
	assert.Nil(t, err)

	os.Stderr = w
	os.Stdout = w
	assert.Nil(t, os.Chdir(root))

	bufChannel := make(chan string)

	go func() {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, r)
		assert.Nil(t, err)
		assert.Nil(t, r.Close())

		bufChannel <- buf.String()
	}()

	exitCode := mainCmd(exec)

	assert.Nil(t, w.Close())

	os.Stderr = saveStderr
	os.Stdout = saveStdout
	assert.Nil(t, os.Chdir(saveCwd))

	out := <-bufChannel

	assert.Equal(t, expectedExitCode, exitCode)
	assert.Contains(t, out, contains)
}
