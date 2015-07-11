package simplegit

import (
	"fmt"
	"gopkg.in/libgit2/git2go.v22"
	"os"
	"strings"
)

type SimpleGit struct {
	*git.Repository

	User          string
	plainTextPass string
}

func NewSimpleGit(repo *git.Repository, creds ...string) *SimpleGit {
	var (
		user = ""
		pass = ""
	)

	switch len(creds) {
	case 1:
		user = creds[0]
		break
	case 2:
		user = creds[0]
		pass = creds[1]
		break
	default:
		panic(fmt.Errorf("Invalid cred args!"))
	}

	return &SimpleGit{repo, user, pass}
}

func (s *SimpleGit) detectGitUser() string {
	if s.User != "" {
		return s.User
	}

	sgitUser := os.Getenv("SIMPLEGIT_USER")
	if sgitUser == "" {
		sgitUser = os.Getenv("USER")
	}
	return sgitUser
}

func (s *SimpleGit) AddFiles(files ...string) error {
	idx, err := s.Index()
	if err != nil {
		return err
	}

	for _, fl := range files {
		if err = idx.AddByPath(fl); err != nil {
			return err
		}
	}
	return nil
}

func (s *SimpleGit) RemoveFiles(files ...string) error {
	idx, err := s.Index()
	if err != nil {
		return err
	}
	for _, fl := range files {
		if err = idx.RemoveByPath(fl); err != nil {
			return err
		}
	}
	return nil
}

/* Write tree and index */
func (s *SimpleGit) WriteTreeIndex() (treeId *git.Oid, err error) {
	var idx *git.Index

	if idx, err = s.Index(); err != nil {
		return
	}

	if treeId, err = idx.WriteTree(); err != nil {
		return
	}

	err = idx.Write()
	return
}

func (s *SimpleGit) plainCredentialsCallback(url string, url_username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredUserpassPlaintext(s.detectGitUser(), s.plainTextPass)
	return git.ErrorCode(ret), &cred
}

func (s *SimpleGit) sshCredentialsCallback(url string, url_username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	ret, cred := git.NewCredSshKeyFromAgent(s.detectGitUser())
	return git.ErrorCode(ret), &cred
}

func (b *SimpleGit) PushHead(remote *git.Remote, projName string) error {
	cbs := &git.RemoteCallbacks{}
	if strings.HasPrefix(remote.Url(), "http") {
		cbs.CredentialsCallback = b.plainCredentialsCallback
	} else {
		cbs.CredentialsCallback = b.sshCredentialsCallback
	}

	if err := remote.SetCallbacks(cbs); err != nil {
		return err
	}

	sign, err := b.DefaultSignature()
	if err != nil {
		return err
	}
	return remote.Push([]string{"refs/heads/master"}, nil, sign, "")
}

func (b *SimpleGit) CommitHead(treeId *git.Oid, message string) error {
	var (
		tree *git.Tree
		sign *git.Signature
		err  error
	)

	tree, err = b.LookupTree(treeId)
	if err != nil {
		return err
	}

	sign, err = b.DefaultSignature()
	if err != nil {
		return err
	}

	head, err := b.Head()
	if err != nil {
		return err
	}

	headCommit, err := b.LookupCommit(head.Target())
	if err != nil {
		return err
	}

	_, err = b.CreateCommit("HEAD", sign, sign, message, tree, headCommit)
	return err
}
