package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const rdtpFilePath = "~/.rdtp"

// Statefile contains a map of ports in use to
// the time they were opened at (epoch nanoseconds)
type Statefile struct {
	Ports map[uint16]int64 `yaml:"ports"`
}

// getState reads state from a statefile at the given path
func getState(path string) (*Statefile, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not read statefile")
	}
	var state *Statefile
	if err := yaml.Unmarshal(dat, &state); err != nil {
		return nil, fmt.Errorf("could not unmarshal statefile: %s", err)
	}
	return state, nil
}

// commit overwrites a statefile at a given path, with the current state of s
func (s *Statefile) commit(path string) error {
	stateByt, err := yaml.Marshal(&s)
	if err != nil {
		return fmt.Errorf("could not marshal statefile: %s", err)
	}
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create padlfile: %s", err)
	}
	if _, err = fd.Write(stateByt); err != nil {
		return fmt.Errorf("could not write padlfile: %s", err)
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
