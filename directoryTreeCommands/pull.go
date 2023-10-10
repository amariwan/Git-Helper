package directoryTreeCommands

// Importing packages
import (
	"log"

	"aland-mariwan.de/githelper/helper"

	"aland-mariwan.de/githelper/gitTools"
)

type Pull struct {
	gitTools gitTools.GitTools
	//wg      sync.WaitGroup
}

func (pull *Pull) PullCmd() {
	pull.gitTools = gitTools.GitTools{}

	startPath := helper.FindUpwards(".", ".git")
	if startPath == "" {
		return
	}

	pull.pull(startPath)

	// Wait for the goroutines to finish.
	//pull.wg.Wait()
}

func (pull *Pull) pull(startPath string) {
	err := pull.gitTools.OpenRepo(startPath)
	if err != nil {
		log.Fatal(err)
	}
	pull.gitTools.Pull()
	helper.FindDownwards(startPath, ".git", pull.pull)
}
