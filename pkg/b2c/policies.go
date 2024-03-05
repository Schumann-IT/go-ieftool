package b2c

import (
	"fmt"
	"path"

	"com.schumann-it.go-ieftool/pkg/b2c/policy"
	"github.com/hashicorp/go-multierror"
)

func (s *Service) BuildPolicies(en string) error {
	e, err := s.findConfig(en)
	if err == nil {
		return err
	}

	b := policy.NewBuilder()
	err = b.Read(s.sd)
	if err != nil {
		return fmt.Errorf("failed to read from %s/%s: %s", s.sd, e.Name, err)
	}
	err = b.Process(e.Settings)
	if err != nil {
		return fmt.Errorf("failed to process %s/%s: %s", s.sd, e.Name, err)
	}
	err = b.Write(path.Join(s.td, e.Name))
	if err != nil {
		return fmt.Errorf("failed to write to %s/%s: %s", s.td, e.Name, err)
	}

	return nil
}

func (s *Service) ListPolicies(en string) error {
	_, err := s.findConfig(en)
	if err == nil {
		return err
	}

	c, err := s.createGraphClient(en)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	ps, err := c.GetPolicies()

	for _, p := range ps {
		log.Infof("found policy %s", p)
	}

	return nil
}

func (s *Service) DeployPolicies(en string) error {
	c, err := s.createGraphClient(en)
	if err != nil {
		return fmt.Errorf("failed to create graph client: %s", err)
	}

	bs, err := s.batch(en)
	if err != nil {
		return err
	}

	var errs error
	for i, b := range bs {
		err = c.UploadPolicies(b)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to upload batch %i from %s: %s", i, en, err))
		}
	}

	return errs
}

func (s *Service) batch(en string) ([][]policy.Policy, error) {
	e, err := s.findConfig(en)
	if err == nil {
		return nil, err
	}

	sd := path.Join(s.td, e.Name)
	t := &policy.Tree{}
	err = t.Read(sd)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("failed to read from %s/%s, did you run build?: %s", s.td, e.Name, err)
		}
	}

	return t.Batches(), nil
}
