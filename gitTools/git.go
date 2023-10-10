package gitTools

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	sshgit "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

type GitTools struct {
	Authentication transport.AuthMethod
	Repo           *git.Repository
	Worktree       *git.Worktree
	RepoPath       string
	Branch         string
	Branches       []string
	Tags           []string
}

func (gt *GitTools) Authenticate(repo string, username string, password string) error {
	if strings.HasPrefix(repo, "http") {
		gt.Authentication = &http.BasicAuth{
			Username: username,
			Password: password,
		}
		return nil
	}

	if strings.HasPrefix(repo, "git@") {
		sshPath := os.Getenv("HOME") + "/.ssh/id_rsa"
		sshKey, _ := ioutil.ReadFile(sshPath)

		var signer ssh.Signer
		var err error
		if len(password) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(password))
			if err != nil {
				return err
			}
		} else {
			signer, err = ssh.ParsePrivateKey(sshKey)
			if err != nil {
				return err
			}
		}

		gt.Authentication = &sshgit.PublicKeys{
			User:   "git",
			Signer: signer,
			HostKeyCallbackHelper: sshgit.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			},
		}

		return nil
	}

	gt.Authentication = &sshgit.Password{
		User:     username,
		Password: password,
		HostKeyCallbackHelper: sshgit.HostKeyCallbackHelper{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
	return nil
}

func (gt *GitTools) OpenRepo(dir string) (err error) {
	gt.Repo, err = git.PlainOpen(dir)
	if err != nil {
		return
	}

	h, err := gt.Repo.Head()
	if err != nil {
		return
	}
	gt.Branch = strings.TrimPrefix(string(h.Name()), "refs/heads/")

	gt.RepoPath = dir

	gt.Worktree, err = gt.Repo.Worktree()
	if err != nil {
		log.Printf("error: can't get worktree for repo %s - %s", dir, err)
		return
	}

	return
}

func (gt *GitTools) OpenRepoCwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	return gt.OpenRepo(dir)
}

func (gt *GitTools) Pull() (err error) {
	remotes, err := gt.Repo.Remotes()
	if err != nil {
		log.Println(err)
		return
	}
	err = gt.Authenticate(remotes[0].Config().URLs[0], "", "")
	if err != nil {
		log.Println(err)
		return
	}

	gitPullOptions := git.PullOptions{
		Auth:              gt.Authentication,
		RemoteName:        remotes[0].Config().Name,
		ReferenceName:     plumbing.NewBranchReferenceName(gt.Branch),
		SingleBranch:      true,
		Depth:             0,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          nil,
	}

	err = gt.Worktree.Pull(&gitPullOptions)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Printf("error: pulling repo %s - %s", gt.RepoPath, err)
		return
	}

	return
}

func (gt *GitTools) Clone(repository string, reference string, destination string) (err error) {
	err = gt.Authenticate(repository, "", "")
	if err != nil {
		log.Println(err)
		return
	}

	gitCloneOptions := git.CloneOptions{
		URL:               repository,
		Auth:              gt.Authentication,
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

	return
}

func (gt *GitTools) Status() (err error) {
	gitStatus, err := gt.Worktree.Status()
	if err != nil {
		log.Printf("error: status repo %s - %s", gt.RepoPath, err)
		return
	}

	if gitStatus.IsClean() {
		log.Printf("git: status repo %s - %s", gt.RepoPath, "clean")
	} else {
		log.Printf("git: status repo %s - %s", gt.RepoPath, "dirty")
	}

	return
}
