package msgraph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"com.schumann-it.go-ieftool/pkg/b2c/policy"
)

func (c *Client) ListPolicies() ([]string, error) {
	r, err := c.gc.TrustFramework().Policies().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var i []string
	for _, p := range r.GetValue() {
		id := p.GetId()
		i = append(i, *id)
	}

	return i, nil
}

func (c *Client) DeletePolicies() error {
	ps, err := c.ListPolicies()
	if err != nil {
		return err
	}

	for _, id := range ps {
		err = c.gc.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			log.Errorf("Failed to delete Policy %s: %s", id, err)
			continue
		}
		log.Debugf(fmt.Sprintf("Policy %s deleted", id))
	}

	return nil
}

func (c *Client) UploadPolicies(policies []policy.Policy) {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	for _, p := range policies {
		go c.uploadPolicy(p, &wg)
	}
	wg.Wait()
}

func (c *Client) uploadPolicy(p policy.Policy, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	defer client.CloseIdleConnections()

	ep := fmt.Sprintf("https://graph.microsoft.com/beta/trustFramework/policies/%s/$value", p.Id())
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(p.Byte()))
	if err != nil {
		log.Fatal(err)
	}

	t, err := c.Token()
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		log.Fatalf("Upload failed for Policy %s \n%s\n", p.Id(), string(body))
	}

	log.Debugf(fmt.Sprintf("Policy %s uploaded", p.Id()))
}
