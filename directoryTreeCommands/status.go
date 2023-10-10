package directoryTreeCommands

// Importing packages
import (
	"log"

	"aland-mariwan.de/githelper/gitTools"
	"aland-mariwan.de/githelper/helper"
)

type Status struct {
	gitTools gitTools.GitTools
	//wg      sync.WaitGroup
}

func (status *Status) StatusCmd() {
	status.gitTools = gitTools.GitTools{}

	startPath := helper.FindUpwards(".", ".git")
	if startPath == "" {
		return
	}

	status.status(startPath)

	// Wait for the goroutines to finish.
	//status.wg.Wait()
}

func (status *Status) status(startPath string) {
	err := status.gitTools.OpenRepo(startPath)
	if err != nil {
		log.Fatal(err)
	}
	status.gitTools.Status()
	helper.FindDownwards(startPath, ".git", status.status)
}
