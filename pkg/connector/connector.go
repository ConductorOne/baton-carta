package connector

import (
	"context"

	"github.com/ConductorOne/baton-carta/pkg/carta"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

var (
	resourceTypeIssuer = &v2.ResourceType{
		Id:          "issuer",
		DisplayName: "Issuer",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
	resourceTypeInvestor = &v2.ResourceType{
		Id:          "investor",
		DisplayName: "Investor",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
)

type Carta struct {
	client *carta.Client
}

func (c *Carta) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		issuerBuilder(c.client),
		investorBuilder(c.client),
	}
}

func (c *Carta) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Carta",
	}, nil
}

func (c *Carta) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns the Carta connector.
func New(ctx context.Context, accessToken string) (*Carta, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))

	if err != nil {
		return nil, err
	}

	return &Carta{
		client: carta.NewClient(accessToken, httpClient),
	}, nil
}
