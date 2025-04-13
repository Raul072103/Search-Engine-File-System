package usn

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrParsingJournalQuery = errors.New("couldn't find Next USN Field")
)

type executor struct {
	config ExecutorConfig
}

type ExecutorConfig struct {
	USNLogsPath string
	NextUSNPath string
	CurrentUSN  string
	NextUSN     string
}

type Query struct {
	NextUSN string `json:"next_usn"`
}

func newExecutor(config ExecutorConfig) *executor {
	if data, err := os.ReadFile(config.NextUSNPath); err == nil {
		var progress struct {
			NextUSN string `json:"next_usn"`
		}
		if err := json.Unmarshal(data, &progress); err == nil {
			config.NextUSN = progress.NextUSN
		}
	}
	return &executor{config: config}
}

// ExecuteReadUSNJournal executes the command to read USN Journal and saves the logs at the given path in configuration
func (e *executor) ExecuteReadUSNJournal() error {
	var cmdStr string
	if e.config.NextUSN != "" {
		cmdStr = fmt.Sprintf("fsutil usn readjournal C: startusn=%s", e.config.NextUSN)
	} else {
		cmdStr = "fsutil usn readjournal C:"
	}

	cmd := exec.Command("powershell", "-Command", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	err = os.WriteFile(e.config.USNLogsPath, output, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ExecuteQueryUSNJournal executes the command to query the USN Journal and saves next USN journal log
func (e *executor) ExecuteQueryUSNJournal() error {
	cmd := exec.Command("powershell", "-Command", "fsutil usn queryjournal C:")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	var nextUSN string
	for _, line := range lines {
		if strings.HasPrefix(line, "Next Usn") {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				return ErrParsingJournalQuery
			} else {
				nextUSN = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if nextUSN != "" {
		// check for USN Journal Overflow
		if e.config.NextUSN != "" {
			if strings.Compare(e.config.NextUSN, nextUSN) > 0 {
				// USN Journal Overflow at some point
				nextUSN = ""
			}
		}
		e.config.NextUSN = nextUSN
		err := e.saveNextUSNJournalNumber()
		if err != nil {
			return nil
		}
	} else {
		return ErrParsingJournalQuery
	}

	return nil
}

// saveNextUSNJournalNumber saves the next USN at the specified path in the configuration
func (e *executor) saveNextUSNJournalNumber() error {
	progress := Query{NextUSN: e.config.NextUSN}

	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(e.config.NextUSNPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
