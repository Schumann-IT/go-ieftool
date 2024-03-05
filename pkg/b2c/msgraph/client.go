package msgraph

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type Client struct {
	gc *sdk.GraphServiceClient
	ac azcore.TokenCredential
	s  []string
}

func NewClient(tid, cid, s string) (*Client, error) {
	log.Debugf("creating client for: %s", tid)

	g := &Client{
		s: []string{"https://graph.microsoft.com/.default"},
	}

	cr, err := azidentity.NewClientSecretCredential(
		tid,
		cid,
		s,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create client credentials: %s", err.Error())
	}
	g.ac = cr

	c, err := sdk.NewGraphServiceClientWithCredentials(cr, g.s)
	if err != nil {
		return nil, err
	}
	g.gc = c

	return g, nil
}

func NewClientFromEnvironment() (*Client, error) {
	g := &Client{
		s: []string{"https://graph.microsoft.com/.default"},
	}

	cr, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create client credentials: %s", err.Error())
	}
	g.ac = cr

	c, err := sdk.NewGraphServiceClientWithCredentials(cr, g.s)
	if err != nil {
		return nil, err
	}
	g.gc = c

	return g, nil
}

func (c *Client) Token() (*azcore.AccessToken, error) {
	t, err := c.ac.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: c.s,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get token: %s", err.Error())
	}

	return &t, nil
}
