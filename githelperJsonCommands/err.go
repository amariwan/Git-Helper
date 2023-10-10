package githelperJsonCommands

// Importing packages
import (
	"errors"
)

var (
	ErrGitRepoNotExists = errors.New("repo '%s' does not exist; try 'githelper clone' or 'githelper sync'")
)
