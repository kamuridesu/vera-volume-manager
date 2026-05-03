package state

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type States struct {
	States   map[string]bool
	FilePath string
}

func getAbsolutePath(configFile string) (string, error) {
	abs, err := filepath.Abs(configFile)
	if err != nil {
		return configFile, fmt.Errorf("failed to get abs path: %w", err)
	}
	return abs, nil
}

func New() (*States, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("error fetching user config dir, error is '%w'", err)
	}

	statePath := path.Join(configDir, "vvm")
	err = os.MkdirAll(statePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating config folder: %w", err)
	}

	stateFilePath := path.Join(statePath, "state.yaml")

	if _, err := os.Stat(stateFilePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(stateFilePath, []byte(""), 0644); err != nil {
				return nil, fmt.Errorf("failed to create initial state file: %w", err)
			}
		} else {
			return nil, err
		}
	}

	content, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	states := make(map[string]bool)

	if len(content) > 0 {
		err = yaml.Unmarshal(content, &states)
		if err != nil {
			return nil, fmt.Errorf("failed to parse state file: %w", err)
		}
	}

	return &States{States: states, FilePath: stateFilePath}, nil
}

func (s *States) writeToFile() error {
	content, err := yaml.Marshal(s.States)
	if err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}

	err = os.WriteFile(s.FilePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func (s *States) SaveState(configFile string, isMounted bool) error {
	configFile, err := getAbsolutePath(configFile)
	if err != nil {
		return err
	}
	s.States[configFile] = isMounted
	return s.writeToFile()
}
