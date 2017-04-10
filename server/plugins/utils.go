package plugins

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	git2go "github.com/libgit2/git2go"
)

func GitCommits(headHash string, project Project, git Git) ([]GitCommit, error) {
	var err error
	var commits []GitCommit
	var repo *git2go.Repository
	var walk *git2go.RevWalk

	repo, err = GitCheckout(project, git)
	if err != nil {
		return nil, err
	}

	walk, err = repo.Walk()
	if err != nil {
		return nil, err
	}
	defer walk.Free()

	err = walk.PushHead()
	if err != nil {
		return nil, err
	}

	i := 0
	callback := func(obj *git2go.Commit) bool {
		i += 1

		if headHash == "" && i > 10 {
			return false
		}

		if headHash == obj.Id().String() {
			return false
		}

		commit := new(GitCommit)
		author := obj.Author()
		committer := obj.Committer()
		commit.Repository = project.Repository
		commit.Hash = obj.Id().String()
		if obj.Parent(0) != nil {
			commit.ParentHash = obj.Parent(0).Id().String()
		}
		commit.Ref = fmt.Sprintf("refs/heads/%s", git.Branch)
		commit.User = author.Name
		commit.Created = committer.When
		commit.Message = obj.Message()
		commits = append(commits, *commit)

		for p := 0; p < int(obj.ParentCount()); p++ {
			if obj.Parent(uint(i)) != nil {
				if headHash == obj.Parent(uint(i)).Id().String() {
					return false
				}
			}
		}

		return true
	}

	err = walk.Iterate(callback)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func GitFetch(project Project, git Git) (*git2go.Repository, error) {
	var repo *git2go.Repository
	var err error

	repoPath := fmt.Sprintf("%s/%s", git.Workdir, project.Repository)

	if _, err = os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			repo, err = git2go.Clone(git.Url, repoPath, GitCloneOptions(git))
		} else {
			return &git2go.Repository{}, err
		}
	} else {
		repo, err = git2go.OpenRepository(repoPath)
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return nil, err
	}

	err = remote.Fetch([]string{}, GitFetchOptions(git), "")
	if err != nil {
		return nil, err
	}

	return repo, err
}

func GitCheckout(project Project, git Git) (*git2go.Repository, error) {
	var repo *git2go.Repository
	var err error

	repo, err = GitFetch(project, git)
	if err != nil {
		return nil, err
	}

	checkoutOpts := &git2go.CheckoutOpts{
		Strategy: git2go.CheckoutForce,
	}

	//Getting the reference for the remote branch
	// remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + git.Branch)
	remoteBranch, err := repo.LookupBranch("origin/"+git.Branch, git2go.BranchRemote)
	if err != nil {
		log.Print("Failed to find remote branch: " + git.Branch)
		return nil, err
	}
	defer remoteBranch.Free()

	// Lookup for commit from remote branch
	commit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		log.Print("Failed to find remote branch commit: " + git.Branch)
		return nil, err
	}
	defer commit.Free()

	localBranch, err := repo.LookupBranch(git.Branch, git2go.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {
		// Creating local branch
		localBranch, err = repo.CreateBranch(git.Branch, commit, false)
		if err != nil {
			log.Print("Failed to create local branch: " + git.Branch)
			return nil, err
		}

		// Setting upstream to origin branch
		err = localBranch.SetUpstream("origin/" + git.Branch)
		if err != nil {
			log.Print("Failed to create upstream to origin/" + git.Branch)
			return nil, err
		}
	}

	if localBranch == nil {
		return nil, errors.New("Error while locating/creating local branch")
	}
	defer localBranch.Free()

	// Getting the tree for the branch
	localCommit, err := repo.LookupCommit(localBranch.Target())
	if err != nil {
		log.Print("Failed to lookup for commit in local branch " + git.Branch)
		return nil, err
	}
	defer localCommit.Free()

	tree, err := repo.LookupTree(localCommit.TreeId())
	if err != nil {
		log.Print("Failed to lookup for tree " + git.Branch)
		return nil, err
	}
	defer tree.Free()

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		log.Print("Failed to checkout tree " + git.Branch)
		return nil, err
	}

	// Setting the Head to point to our branch
	repo.SetHead("refs/heads/" + git.Branch)

	return repo, err
}

func GitCheckoutCommit(hash string, project Project, git Git) (*git2go.Repository, error) {
	var repo *git2go.Repository
	var err error

	repo, err = GitFetch(project, git)
	if err != nil {
		return nil, err
	}

	checkoutOpts := &git2go.CheckoutOpts{
		Strategy: git2go.CheckoutForce,
	}

	oid, err := git2go.NewOid(hash)
	if err != nil {
		return repo, err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return repo, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return repo, err
	}

	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return repo, err
	}

	err = repo.SetHeadDetached(oid)
	if err != nil {
		return repo, err
	}

	return repo, nil
}

type Git2goCallback struct {
	rsaPrivateKey string
	rsaPublicKey  string
}

func (x *Git2goCallback) GitCredentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	ret, cred := git2go.NewCredSshKeyFromMemory("git", x.rsaPublicKey, x.rsaPrivateKey, "")
	return git2go.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func (x *Git2goCallback) GitCertificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	return git2go.ErrorCode(0)
}

func GitFetchOptions(git Git) *git2go.FetchOptions {
	var fetchOptions git2go.FetchOptions

	git2goCallback := Git2goCallback{
		rsaPrivateKey: git.RsaPrivateKey,
		rsaPublicKey:  git.RsaPublicKey,
	}

	if git.Protocol == "SSH" {
		fetchOptions = git2go.FetchOptions{
			RemoteCallbacks: git2go.RemoteCallbacks{
				CredentialsCallback:      git2goCallback.GitCredentialsCallback,
				CertificateCheckCallback: git2goCallback.GitCertificateCheckCallback,
			},
		}
	} else {
		fetchOptions = git2go.FetchOptions{}
	}

	return &fetchOptions
}

func GitCloneOptions(git Git) *git2go.CloneOptions {
	cloneOptions := git2go.CloneOptions{
		FetchOptions: GitFetchOptions(git),
		CheckoutOpts: &git2go.CheckoutOpts{
			Strategy: git2go.CheckoutForce,
		},
	}
	return &cloneOptions
}

func GetRegexParams(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}
