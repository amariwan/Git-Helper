package githelperJsonCommands

// Importing packages
import (
	"log"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"aland-mariwan.de/githelper/gitTools"
	"aland-mariwan.de/githelper/helper"

	"aland-mariwan.de/githelper/githelperConfig"
)

type CloneAndPull struct {
	DoClone  bool
	DoPull   bool
	//wg      sync.WaitGroup
}

func (cloneAndPull *CloneAndPull) CloneAndPullCmd() {
	cloneAndPull.CloneAndPull(".")

	// Wait for the goroutines to finish.
	//cloneAndPull.wg.Wait()
}

func (cloneAndPull *CloneAndPull) CloneAndPull(startPath string) {
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
			if cloneAndPull.DoPull {

				referenceToClone := githelperConfig.VersionReferenced(v.Versions, gitTools.Branch)
				log.Printf("git \n"+
					"\t          pull: %s \n"+
					"\tusing reference: %s \n"+
					"\t             to: %s",
					v.Repository, referenceToClone, repositoryDir)
				//cloneAndPull.wg.Add(1)
				cloneAndPull.pullSingle(v.Repository, referenceToClone, repositoryDir)
			}
		} else {
			if cloneAndPull.DoClone {
				referenceToClone := githelperConfig.VersionReferenced(v.Versions, gitTools.Branch)
				log.Printf("git \n"+
					"\t          clone: %s \n"+
					"\tusing reference: %s \n"+
					"\t             to: %s",
					v.Repository, referenceToClone, repositoryDir)
				//cloneAndPull.wg.Add(1)
				cloneAndPull.cloneSingle(v.Repository, referenceToClone, repositoryDir)
			} else {
				if cloneAndPull.DoPull {
					log.Fatalf(ErrGitRepoNotExists.Error(), repositoryDir)
				}
			}
		}
	}
}

func (cloneAndPull *CloneAndPull) cloneSingle(repository string, reference string, destination string) (err error) {

	gitTools := gitTools.GitTools{}
	err = gitTools.Authenticate(repository, "", "")
	if err != nil {
		log.Println(err)
		return
	}

	gitCloneOptions := git.CloneOptions{
		URL:               repository,
		Auth:              gitTools.Authentication,
		RemoteName:        "origin",
		ReferenceName:     plumbing.NewBranchReferenceName(reference),
		SingleBranch:      true,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          nil,
		Tags:              git.AllTags,
	}

	_, err = git.PlainClone(destination, false, &gitCloneOptions)
	if err != nil {
		log.Printf("error: cloning repo %s - %s", repository, err)
		return
	}

	helper.FindDownwards(destination, "githelper.json", cloneAndPull.CloneAndPull)

	return
}

func (cloneAndPull *CloneAndPull) pullSingle(repository string, reference string, destination string) (err error) {
	gitTools := gitTools.GitTools{}
	err = gitTools.Authenticate(repository, "", "")
	if err != nil {
		log.Println(err)
		return
	}

	repo, err := git.PlainOpen(destination)
	if err != nil {
		return
	}

	gitPullOptions := git.PullOptions{
		Auth:              gitTools.Authentication,
		RemoteName:        "origin",
		ReferenceName:     plumbing.NewBranchReferenceName(reference),
		SingleBranch:      true,
		Depth:             0,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          nil,
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Printf("error: pulling repo %s - %s", repository, err)
		return
	}

	status, err := worktree.Status()
	if err != nil {
		log.Printf("error: pulling repo %s - %s", repository, err)
		return
	}

	if status.IsClean() {
		err = worktree.Pull(&gitPullOptions)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Printf("error: pulling repo %s - %s", repository, err)
			return
		}
	}

	helper.FindDownwards(destination, "githelper.json", cloneAndPull.CloneAndPull)

	return
}
