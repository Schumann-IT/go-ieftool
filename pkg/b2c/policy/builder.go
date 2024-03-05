package policy

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"com.schumann-it.go-ieftool/pkg/b2c/policy/content"
	"github.com/hashicorp/go-multierror"
)

// Builder is a type that represents a content builder. It is used to read, process, and write content.
type Builder struct {
	s content.Source
	p content.Processed
}

// NewBuilder creates a new instance of Builder and initializes it with empty File and Processed contents
func NewBuilder() *Builder {
	b := &Builder{}
	b.Reset()

	return b
}

// Reset sets the Builder instance back to an initial state, with empty source root and empty content fields
func (b *Builder) Reset() {
	b.s = content.Source{}
	b.p = content.Processed{}
}

// Len returns the length of the content.Source map in the Builder instance.
func (b *Builder) Len() int {
	return b.s.Len()
}

// Read reads files from a specified directory and populates the Builder's source field.
// The absolute path to the directory is passed as the 'from' parameter.
// If the provided path is not absolute, an error is returned.
// The method first calls the Reset() method to set the Builder instance to its initial state.
// Then, it sets the source root field to the provided path.
// It uses filepath.WalkDir to walk through the files and directories in the source root directory.
// For each file with the extension ".xml", it reads the file contents using os.ReadFile.
// The file path is modified by removing the root directory path.
// The modified path and file contents are stored in the Builder's source map.
// If any errors occur during the file reading or path modification, the method returns the error.
// Finally, the method returns the error encountered during the filepath.WalkDir operation, if any.
func (b *Builder) Read(from string) error {
	if !filepath.IsAbs(from) {
		return fmt.Errorf("path must be absolute: %s", from)
	}

	b.Reset()
	err := filepath.WalkDir(from, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			log.Debugf("reading dir: %s", p)
			return nil
		}
		if filepath.Ext(e.Name()) == ".xml" {
			log.Debugf("found: %s", p)
			cb, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			// remove root dir path from file path
			pp := path.Join(strings.ReplaceAll(filepath.Dir(p), from, ""), e.Name())
			b.s[pp] = cb
		}
		return nil
	})

	return err
}

// Process processes the content of the Builder instance by replacing variables with their corresponding values.
// It takes a map of variables and their values as input.
// The method returns an error if there is nothing to process or if there are any errors encountered during the variable
// replacement.
func (b *Builder) Process(c map[string]string) error {
	if b.Len() == 0 {
		return fmt.Errorf("nothing to process")
	}
	if len(b.p) > 0 {
		// already processed
		return nil
	}

	var errs error
	for p, r := range b.s {
		err := r.ReplaceVariables(c, p)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		b.p[p] = r
	}

	return errs
}

// Result returns the processed content of the Builder instance.
// It returns a map[string][]byte where the key is the file path and the value is the processed content.
// If the Builder instance has not been processed yet, an empty map is returned.
func (b *Builder) Result() content.Processed {
	return b.p
}

// Write writes the processed content to the specified directory.
// It creates directories as needed and writes the content files.
// The "to" parameter must be an absolute path.
// If "to" is not absolute, it returns an error with a message indicating that the path must be absolute.
// If there is no content to write, it returns an error with a message indicating that there is nothing to write.
// It returns an error if any error occurs during directory creation or file writing.
func (b *Builder) Write(to string) error {
	if !filepath.IsAbs(to) {
		return fmt.Errorf("path must be absolute: %s", to)
	}

	if b.p.Len() == 0 {
		return fmt.Errorf("nothing to write")
	}

	var errs error
	for f, r := range b.Result() {
		p := path.Join(to, f)
		err := os.MkdirAll(filepath.Dir(p), os.ModePerm)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		err = os.WriteFile(p, r, os.ModePerm)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}
