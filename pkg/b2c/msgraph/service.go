package msgraph

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type ServiceClient struct {
	gc *msgraph.GraphServiceClient
	ac azcore.TokenCredential
	s  []string
}

var scopes = []string{"https://graph.microsoft.com/.default"}

func NewClientWithCredential(cred azcore.TokenCredential) (*ServiceClient, error) {
	g, err := msgraph.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		return nil, err
	}

	return &ServiceClient{
		s:  scopes,
		ac: cred,
		gc: g,
	}, nil
}

func (s *ServiceClient) Token() (*azcore.AccessToken, error) {
	t, err := s.ac.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: s.s,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get token: %s", err.Error())
	}

	return &t, nil
}
