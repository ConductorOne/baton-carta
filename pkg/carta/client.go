package carta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const BaseURL = "https://mock-api.carta.com/v1alpha1/"
const InvestorsBaseURL = BaseURL + "investors/firms"
const IssuersBaseURL = BaseURL + "issuers"
const PortfolioBaseURL = BaseURL + "portfolios"

type Client struct {
	httpClient  *http.Client
	accessToken string
}

type IssuerResponse struct {
	Issuers []Issuer `json:"issuers"`
	PaginationData
}

type InvestorResponse struct {
	Firms []InvestorFirm `json:"firms"`
	PaginationData
}

type PaginationParams struct {
	Size  int    `json:"pageSize"`
	After string `json:"pageToken"`
}

func NewClient(accessToken string, httpClient *http.Client) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  httpClient,
	}
}

func setupPaginationQuery(query url.Values, size int, after string) url.Values {
	// add size
	if size != 0 {
		query.Add("pageSize", strconv.Itoa(size))
	}

	// add page reference
	if after != "" {
		query.Add("pageToken", after)
	}

	return query
}

// GetIssuers returns all issuers (companies to invest in) accessible to the user or investor.
func (c *Client) GetIssuers(ctx context.Context, getIssuerVars PaginationParams) ([]Issuer, string, error) {
	queryParams := setupPaginationQuery(url.Values{}, getIssuerVars.Size, getIssuerVars.After)
	var userResponse IssuerResponse

	err := c.doRequest(
		ctx,
		IssuersBaseURL,
		&userResponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getIssuerVars.After != userResponse.Next && userResponse.Next != "" {
		return userResponse.Issuers, userResponse.Next, nil
	}

	return userResponse.Issuers, "", nil
}

// GetInvestors returns all investor firms accessible to the user.
func (c *Client) GetInvestors(ctx context.Context, getInvestorVars PaginationParams) ([]InvestorFirm, string, error) {
	queryParams := setupPaginationQuery(url.Values{}, getInvestorVars.Size, getInvestorVars.After)
	var userResponse InvestorResponse

	err := c.doRequest(
		ctx,
		InvestorsBaseURL,
		&userResponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getInvestorVars.After != userResponse.Next && userResponse.Next != "" {
		return userResponse.Firms, userResponse.Next, nil
	}

	return userResponse.Firms, "", nil
}

func (c *Client) doRequest(ctx context.Context, url string, resourceResponse interface{}, queryParams url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if queryParams != nil {
		req.URL.RawQuery = queryParams.Encode()
	}

	req.Header.Add("authorization", fmt.Sprint("Bearer ", c.accessToken))
	req.Header.Add("accept", "application/json")

	rawResponse, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	if err := json.NewDecoder(rawResponse.Body).Decode(&resourceResponse); err != nil {
		return err
	}

	return nil
}
