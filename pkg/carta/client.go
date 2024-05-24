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
const IssuerBaseURL = IssuersBaseURL + "/%s"
const PortfoliosBaseURL = BaseURL + "portfolios"
const PortfoliosIssuersBaseURL = PortfoliosBaseURL + "/%s/issuers"

type Client struct {
	httpClient  *http.Client
	accessToken string
}

type IssuerResponse struct {
	Issuer Issuer `json:"issuer"`
}

type IssuersResponse struct {
	Issuers []Issuer `json:"issuers"`
	PaginationData
}

type PortfoliosResponse struct {
	Portfolios []Portfolio `json:"portfolios"`
	PaginationData
}

type PortfoliosIssuersResponse struct {
	Issuers []Issuer `json:"issuers"`
	PaginationData
}

type InvestorsResponse struct {
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
	var issuersResponse IssuersResponse

	err := c.doRequest(
		ctx,
		IssuersBaseURL,
		&issuersResponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getIssuerVars.After != issuersResponse.Next && issuersResponse.Next != "" {
		return issuersResponse.Issuers, issuersResponse.Next, nil
	}

	return issuersResponse.Issuers, "", nil
}

// GetIssuer returns specific issuer based on provided id, accessible to the user or investor.
func (c *Client) GetIssuer(ctx context.Context, issuerId string) (Issuer, error) {
	var issuerResponse IssuerResponse

	err := c.doRequest(
		ctx,
		fmt.Sprintf(IssuerBaseURL, issuerId),
		&issuerResponse,
		nil,
	)

	if err != nil {
		return Issuer{}, err
	}

	return issuerResponse.Issuer, nil
}

// GetPortfolios returns all portfolios (groupings of issuers) accessible to the user or investor.
func (c *Client) GetPortfolios(ctx context.Context, getPortfolioVars PaginationParams) ([]Portfolio, string, error) {
	queryParams := setupPaginationQuery(url.Values{}, getPortfolioVars.Size, getPortfolioVars.After)
	var portfoliosResponse PortfoliosResponse

	err := c.doRequest(
		ctx,
		PortfoliosBaseURL,
		&portfoliosResponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// get all issuers for each portfolio
	for i, portfolio := range portfoliosResponse.Portfolios {
		var issuers []Issuer
		var next string

		// get issuers for portfolio ( loop until all issuers are retrieved )
		for {
			issuersForPortfolio, nextToken, err := c.GetIssuersForPortfolio(
				ctx,
				portfolio.Id,
				PaginationParams{Size: 100, After: next},
			)

			if err != nil {
				return nil, "", err
			}

			issuers = append(issuers, issuersForPortfolio...)

			if nextToken == "" {
				break
			}

			next = nextToken
		}

		portfoliosResponse.Portfolios[i].Issuers = issuers
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getPortfolioVars.After != portfoliosResponse.Next && portfoliosResponse.Next != "" {
		return portfoliosResponse.Portfolios, portfoliosResponse.Next, nil
	}

	return portfoliosResponse.Portfolios, "", nil
}

// GetIssuersForPortfolio returns all issuers (companies to invest in) under specific portfolio.
func (c *Client) GetIssuersForPortfolio(ctx context.Context, portfolioId string, getIssuerVars PaginationParams) ([]Issuer, string, error) {
	queryParams := setupPaginationQuery(url.Values{}, getIssuerVars.Size, getIssuerVars.After)
	var issuersReponse PortfoliosIssuersResponse

	err := c.doRequest(
		ctx,
		fmt.Sprintf(PortfoliosIssuersBaseURL, portfolioId),
		&issuersReponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getIssuerVars.After != issuersReponse.Next && issuersReponse.Next != "" {
		return issuersReponse.Issuers, issuersReponse.Next, nil
	}

	return issuersReponse.Issuers, "", nil
}

// GetInvestors returns all investor firms accessible to the user.
func (c *Client) GetInvestors(ctx context.Context, getInvestorVars PaginationParams) ([]InvestorFirm, string, error) {
	queryParams := setupPaginationQuery(url.Values{}, getInvestorVars.Size, getInvestorVars.After)
	var investorsResponse InvestorsResponse

	err := c.doRequest(
		ctx,
		InvestorsBaseURL,
		&investorsResponse,
		queryParams,
	)

	if err != nil {
		return nil, "", err
	}

	// check for duplicates to prevent infinite loop (this can happen with mock data)
	if getInvestorVars.After != investorsResponse.Next && investorsResponse.Next != "" {
		return investorsResponse.Firms, investorsResponse.Next, nil
	}

	return investorsResponse.Firms, "", nil
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
