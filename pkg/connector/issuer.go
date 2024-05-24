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

type issuerResourceType struct {
	resourceType *v2.ResourceType
	client       *carta.Client
}

func (o *issuerResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return o.resourceType
}

// Create a new connector resource for an Carta Issuer (Company to invest in).
func issuerResource(ctx context.Context, issuer *carta.Issuer, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"issuer_legal_name": issuer.Name,
		"issuer_id":         issuer.Id,
	}

	issuerTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(v2.UserTrait_Status_STATUS_UNSPECIFIED),
	}

	resource, err := rs.NewUserResource(
		issuer.Name,
		resourceTypeIssuer,
		issuer.Id,
		issuerTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (o *issuerResourceType) List(ctx context.Context, parentId *v2.ResourceId, token *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	bag, err := parsePageToken(token.Token, &v2.ResourceId{ResourceType: resourceTypeIssuer.Id})
	if err != nil {
		return nil, "", nil, err
	}

	issuers, nextToken, err := o.client.GetIssuers(
		ctx,
		carta.PaginationParams{Size: ResourcesPageSize, After: bag.PageToken()},
	)
	if err != nil {
		return nil, "", nil, fmt.Errorf("carta-connector: failed to list issuers: %w", err)
	}

	pageToken, err := bag.NextToken(nextToken)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Resource
	for _, issuer := range issuers {
		issuerCopy := issuer
		ir, err := issuerResource(ctx, &issuerCopy, parentId)

		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, ir)
	}

	return rv, pageToken, nil, nil
}

func (o *issuerResourceType) Entitlements(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *issuerResourceType) Grants(ctx context.Context, resource *v2.Resource, token *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func issuerBuilder(client *carta.Client) *issuerResourceType {
	return &issuerResourceType{
		resourceType: resourceTypeIssuer,
		client:       client,
	}
}
