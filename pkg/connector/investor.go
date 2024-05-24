package connector

import (
	"context"
	"fmt"

	"github.com/ConductorOne/baton-carta/pkg/carta"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type investorResourceType struct {
	resourceType *v2.ResourceType
	client       *carta.Client
}

func (o *investorResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for an Carta Investor.
func investorResource(ctx context.Context, investor *carta.InvestorFirm, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"login":       investor.Name,
		"investor_id": investor.Id,
	}

	investorTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(v2.UserTrait_Status_STATUS_UNSPECIFIED),
	}

	resource, err := rs.NewUserResource(
		investor.Name,
		resourceTypeInvestor,
		investor.Id,
		investorTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *investorResourceType) List(ctx context.Context, parentId *v2.ResourceId, token *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, err := parsePageToken(token.Token, &v2.ResourceId{ResourceType: resourceTypeInvestor.Id})
	if err != nil {
		return nil, "", nil, err
	}

	investors, nextToken, err := o.client.GetInvestors(
		ctx,
		carta.PaginationParams{Size: ResourcesPageSize, After: bag.PageToken()},
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("carta-connector: failed to list investors: %w", err)
	}

	pageToken, err := bag.NextToken(nextToken)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Resource
	for _, investor := range investors {
		investorCopy := investor
		ir, err := investorResource(ctx, &investorCopy, parentId)

		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, ir)
	}

	return rv, pageToken, nil, nil
}

func (o *investorResourceType) Entitlements(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *investorResourceType) Grants(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func investorBuilder(client *carta.Client) *investorResourceType {
	return &investorResourceType{
		resourceType: resourceTypeInvestor,
		client:       client,
	}
}
