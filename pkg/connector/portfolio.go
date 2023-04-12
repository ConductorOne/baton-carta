package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/ConductorOne/baton-carta/pkg/carta"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type portfolioResourceType struct {
	resourceType *v2.ResourceType
	client       *carta.Client
}

func (o *portfolioResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for an Carta Portfolio (Grouping entity of issuers).
func portfolioResource(ctx context.Context, portfolio *carta.Portfolio, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"portfolio_legal_name": portfolio.Name,
		"portfolio_id":         portfolio.Id,
		"portfolio_issuer_ids": strings.Join(mapIssuerIds(portfolio.Issuers), ","),
	}

	portfolioTraitOptions := []rs.GroupTraitOption{
		rs.WithGroupProfile(profile),
	}

	resource, err := rs.NewGroupResource(
		portfolio.Name,
		resourceTypePortfolio,
		portfolio.Id,
		portfolioTraitOptions,
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *portfolioResourceType) List(ctx context.Context, parentId *v2.ResourceId, token *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, err := parsePageToken(token.Token, &v2.ResourceId{ResourceType: resourceTypePortfolio.Id})
	if err != nil {
		return nil, "", nil, err
	}

	portfolios, nextToken, err := o.client.GetPortfolios(
		ctx,
		carta.PaginationParams{Size: ResourcesPageSize, After: bag.PageToken()},
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("carta-connector: failed to list portfolios: %w", err)
	}

	pageToken, err := bag.NextToken(nextToken)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Resource
	for _, portfolio := range portfolios {
		portfolioCopy := portfolio
		pr, err := portfolioResource(ctx, &portfolioCopy, parentId)

		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, pr)
	}

	return rv, pageToken, nil, nil
}

func (o *portfolioResourceType) Entitlements(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *portfolioResourceType) Grants(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func portfolioBuilder(client *carta.Client) *portfolioResourceType {
	return &portfolioResourceType{
		resourceType: resourceTypePortfolio,
		client:       client,
	}
}
