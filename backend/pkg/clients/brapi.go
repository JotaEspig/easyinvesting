package clients

import (
	"easyinvesting/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	BrApiErrNoResults = fmt.Errorf("No results found for asset code")
)

type BrApi interface {
	GetQuote(code string) (*BrApiQuote, error)
}

type brApi struct {
	httpClient *http.Client
}

type BrApiQuote struct {
	Symbol             string  `json:"symbol"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type BrApiResponse struct {
	Results []*BrApiQuote `json:"results"`
}

func NewBrApi(httpClient *http.Client) BrApi {
	return &brApi{
		httpClient: httpClient,
	}
}

func (b *brApi) GetQuote(code string) (*BrApiQuote, error) {
	url := "https://brapi.dev/api/quote/" + code
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", code, err)
	}

	req.Header.Set("Authorization", "Bearer "+config.BRAPI_TOKEN)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request for %s: %w", code, err)
	}
	defer resp.Body.Close()

	var data BrApiResponse
	if resp.StatusCode != http.StatusOK {
		errorData := struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}{}

		if err := json.NewDecoder(resp.Body).Decode(&errorData); err != nil {
			return nil, fmt.Errorf("error decoding error response for %s: %w", code, err)
		}
		return nil, errors.New(errorData.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding response for %s: %w", code, err)
	}

	if len(data.Results) == 0 {
		return nil, fmt.Errorf("no results found for asset code: %s", code)
	}

	quote := data.Results[0]
	return quote, nil
}
