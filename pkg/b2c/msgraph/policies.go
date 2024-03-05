package msgraph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"com.schumann-it.go-ieftool/pkg/b2c/policy"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/go-multierror"
)

func (s *ServiceClient) GetPolicies() ([]string, error) {
	d, err := s.gc.TrustFramework().Policies().Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %s", err)
	}

	var r []string
	for _, p := range d.GetValue() {
		r = append(r, to.String(p.GetId()))
	}

	return r, nil
}

func (s *ServiceClient) DeletePolicies() error {
	ps, err := s.GetPolicies()
	if err != nil {
		return err
	}

	var errs error
	for _, id := range ps {
		err = s.gc.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("failed to delete policy %s: %s", id, err))
			continue
		}
		log.Debugf(fmt.Sprintf("sucessfully deleted policy %s", id))
	}

	return nil
}

func (s *ServiceClient) UploadPolicies(policies []policy.Policy) error {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	errChan := make(chan error, len(policies))
	for _, p := range policies {
		go s.uploadPolicy(p, &wg, errChan)
	}
	wg.Wait()

	var errs error
	for err := range errChan {
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (s *ServiceClient) uploadPolicy(p policy.Policy, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	hc := &http.Client{}
	defer hc.CloseIdleConnections()

	ep := fmt.Sprintf("https://graph.microsoft.com/beta/trustFramework/policies/%s/$value", p.Id())
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(p.Byte()))
	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %s", p.Id(), err)
		return
	}

	t, err := s.Token()
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	resp, err := hc.Do(req)

	if err != nil {
		errChan <- fmt.Errorf("failed to upload policy %s: %s", p.Id(), err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		errChan <- err
		return
	}

	if resp.StatusCode >= 400 {
		if err != nil {
			errChan <- fmt.Errorf("failed to upload policy %s: %s", p.Id(), string(body))
			return
		}
	}

	log.Debugf(fmt.Sprintf("sucessfully uploaded policy %s", p.Id()))

	errChan <- nil
}
