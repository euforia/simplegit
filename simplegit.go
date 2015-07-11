package simplegit

import (
	"github.com/libgit2/git2go"
)

type SimpleGit struct {
	*git.Repository
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
