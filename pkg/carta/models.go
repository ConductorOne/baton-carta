package carta

type BaseResource struct {
	Id string `json:"id"`
}

type Issuer struct {
	BaseResource
	Name    string `json:"legalName"`
	Website string `json:"website"`
}

type Portfolio struct {
	Id      string `json:"portfolioId"`
	Name    string `json:"legalName"`
	Issuers []Issuer
}

type InvestorFirm struct {
	BaseResource
	Name string `json:"name"`
}

type PaginationData struct {
	Next string `json:"nextPageToken"`
}
