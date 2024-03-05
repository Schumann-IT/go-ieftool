package policy

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"com.schumann-it.go-ieftool/pkg/b2c/policy/tree"
)

// Tree is a type that represents a collection of policies.
type Tree []Policy

// Read reads policies from the specified absolute path and adds them to the Tree.
// It recursively searches for XML files in the specified directory and its subdirectories.
// Each XML file found is parsed as a Policy using the New method.
// If a Policy with the same ID already exists in the Tree, an error is returned.
// The Tree is modified by appending the new Policies.
// The method returns an error if any file or parsing error occurs during the process.
func (t *Tree) Read(from string) error {
	if !filepath.IsAbs(from) {
		return fmt.Errorf("path must be absolute: %s", from)
	}

	err := filepath.WalkDir(from, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		if filepath.Ext(e.Name()) == ".xml" {
			log.Debugf("found: %s", p)
			p, err := New(p)
			if err != nil {
				return err
			}
			err = t.add(p)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// Batches returns the policies organized into batches. Each batch is a slice of policies.
// It determines the batches by traversing the tree and appending policies to each batch.
// The resulting batches are returned as a two-dimensional slice.
//
// Example usage:
//
//	tree := Tree{}
//	batches := tree.Batches()
//
// Result:
//
//	batches: [][]Policy
func (t *Tree) Batches() [][]Policy {
	var r [][]Policy

	log.Debug("Building Policy Tree")
	rp := t.findRoot()
	rb := tree.NewBranch(rp)

	t.recursiveAddBranch(&rb)

	log.Debug("determining batches...")
	t.batchFrom([]tree.Branch[Policy]{rb}, &r)
	log.Debugf("found %d batches", len(r))

	return r
}

func (t *Tree) add(p Policy) error {
	e := t.find(p)
	if e != nil {
		return fmt.Errorf("policy with id %s already exists in file %s. Tried to add %s", e.Id(), e.File(), p.File())
	}

	*t = append(*t, p)

	return nil
}

func (t *Tree) find(p Policy) Policy {
	for _, e := range *t {
		if e.Id() == p.Id() {
			return e
		}
	}

	return nil
}

func (t *Tree) findRoot() Policy {
	var r Policy
	for _, p := range *t {
		if !p.HasParent() {
			r = p
			t.remove(p)
		}
	}

	return r
}

func (t *Tree) remove(r Policy) {
	var n []Policy
	for _, p := range *t {
		if p.Id() != r.Id() {
			n = append(n, p)
		}
	}

	*t = n
}

func (t *Tree) recursiveAddBranch(parent *tree.Branch[Policy]) {
	childPolicies := t.findChildPolicies(parent.Data())
	if len(childPolicies) == 0 {
		return
	}
	for _, child := range childPolicies {
		branch := tree.NewBranch(child)
		t.recursiveAddBranch(&branch)
		parent.AddChild(branch)
	}
}

func (t *Tree) findChildPolicies(p Policy) []Policy {
	var r []Policy
	for _, e := range *t {
		if e.HasParent() && e.Parent().Id() == p.Id() {
			r = append(r, e)
		}
	}
	return r
}

func (t *Tree) batchFrom(tree []tree.Branch[Policy], policies *[][]Policy) {
	var batch []Policy
	for _, branch := range tree {
		batch = append(batch, branch.Data())
	}
	for _, branch := range tree {
		if len(branch.Children()) > 0 {
			t.batchFrom(branch.Children(), policies)
		}
	}
	*policies = append([][]Policy{batch}, *policies...)
}
