package githelperJsonCommands

// Importing packages
import (
	"log"
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"aland-mariwan.de/githelper/gitTools"

	"aland-mariwan.de/githelper/githelperConfig"
	"aland-mariwan.de/githelper/helper"
)

type Status struct {
	//wg      sync.WaitGroup
}

func (status *Status) StatusCmd() {
	status.status(".")

	// Wait for the goroutines to finish.
	//status.wg.Wait()
}

func (status *Status) status(startPath string) {
	gitHelperConfig, err := githelperConfig.ReadGithelperJson(startPath)
	if err != nil {
		log.Fatal(err)
	}

	gitTools := gitTools.GitTools{}
	err = gitTools.OpenRepo(startPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range gitHelperConfig {
		repositoryDir := filepath.Join(startPath, v.Folder, v.Name)

		if helper.DirExists(repositoryDir) {
			referenceToClone := githelperConfig.VersionReferenced(v.Versions, gitTools.Branch)
			log.Printf("git \n"+
				"\t         status: %s \n"+
				"\tusing reference: %s \n"+
				"\t             to: %s",
				v.Repository, referenceToClone, repositoryDir)
			//status.wg.Add(1)
			status.statusSingle(v.Repository, referenceToClone, repositoryDir)

		} else {
			log.Fatalf(ErrGitRepoNotExists.Error(), repositoryDir)
		}
	}
}

func (status *Status) statusSingle(repository string, reference string, destination string) (err error) {
	repo, err := git.PlainOpen(destination)
	if err != nil {
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Printf("error: pulling repo %s - %s", repository, err)
		return
	}

	gitStatus, err := worktree.Status()
	if err != nil {
		log.Printf("error: status repo %s - %s", repository, err)
		return
	}

	if gitStatus.IsClean() {
		log.Printf("git: status repo %s - %s", repository, "clean")
	} else {
		log.Printf("git: status repo %s - %s", repository, "dirty")
	}

	helper.FindDownwards(destination, "githelper.json", status.status)

	return
}
