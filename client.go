package mph

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mypricehealth/sling"
)

// Pricer is the interface that wraps the Price and PriceBatch methods
type Pricer interface {
	Price(ctx context.Context, config PriceConfig, input Claim) Response[Pricing]
	PriceBatch(ctx context.Context, config PriceConfig, inputs ...Claim) Responses[Pricing]
}

// PriceConfig is used to configure the behavior of the pricing API
type PriceConfig struct {
	IsCommercial                        bool    // set to true to use commercial code crosswalks
	DisableCostBasedReimbursement       bool    // by default, the API will use cost-based reimbursement for MAC priced line-items. This is the best estimate we have for this proprietary pricing
	UseCommercialSyntheticForNotAllowed bool    // set to true to use a synthetic Medicare price for line-items that are not allowed by Medicare
	UseDRGFromGrouper                   bool    // set to true to always use the DRG from the inpatient grouper
	UseBestDRGPrice                     bool    // set to true to use the best DRG price between the price on the claim and the price from the grouper
	OverrideThreshold                   float64 // set to a value greater than 0 to allow the pricer flexibility to override NCCI edits and other overridable errors and return a price
	IncludeEdits                        bool    // set to true to include edit details in the response
}

// Client is used to interact with the My Price Health API
type Client struct {
	*sling.Sling
}

var _ Pricer = &Client{}

// NewClient is used to create a new API client for the My Price Health API. In most cases
// it is simpler to use NewDefaultClient to create a client with the default settings.
func NewClient(doer sling.Doer, isTest bool, apiKey string) *Client {
	url := "https://api.mypricehealth.com"
	if isTest {
		url = "https://api-test.mypricehealth.com"
	}
	client := &Client{sling.New().Doer(doer).Base(url).Set("x-api-key", apiKey)}
	return client
}

// NewDefaultClient is used to create a new API client for the My Price Health API with the default settings.
func NewDefaultClient(apiKey string) *Client {
	return NewClient(http.DefaultClient, false, apiKey)
}

// Price is used to get the Medicare price of a single claim
func (c *Client) Price(ctx context.Context, config PriceConfig, input Claim) Response[Pricing] {
	const path = "/v1/medicare/price/claim"

	var response Response[Pricing]
	res, err := c.Sling.BodyJSON(input).AddHeaders(getHeaders(config)).Post(path).ReceiveWithContext(ctx, &response, &response)
	if err != nil {
		response.Error = &ResponseError{Title: fmt.Sprintf("fatal error calling %s", path), Detail: err.Error()}
		response.StatusCode = res.StatusCode
	}

	return response
}

// PriceBatch is used to get the Medicare price of multiple claims
func (c *Client) PriceBatch(ctx context.Context, config PriceConfig, inputs ...Claim) Responses[Pricing] {
	const path = "/v1/medicare/price/claims"

	var responses Responses[Pricing]
	res, err := c.Sling.BodyJSON(inputs).AddHeaders(getHeaders(config)).Post(path).ReceiveWithContext(ctx, &responses, &responses)
	if err != nil {
		responses.Error = &ResponseError{Title: fmt.Sprintf("fatal error calling %s", path), Detail: err.Error()}
		responses.ErrorCount = len(inputs)
		responses.StatusCode = res.StatusCode
	}

	return responses
}

func getHeaders(config PriceConfig) http.Header {
	headers := http.Header{}
	if config.IsCommercial {
		headers.Add("is-commercial", "true")
	}
	if config.DisableCostBasedReimbursement {
		headers.Add("disable-cost-based-reimbursement", "true")
	}
	if config.UseCommercialSyntheticForNotAllowed {
		headers.Add("use-commercial-synthetic-for-not-allowed", "true")
	}
	if config.OverrideThreshold > 0 {
		headers.Add("override-threshold", strconv.FormatFloat(config.OverrideThreshold, 'f', -1, 64))
	}
	if config.IncludeEdits {
		headers.Add("include-edits", "true")
	}
	if config.UseDRGFromGrouper {
		headers.Add("use-drg-from-grouper", "true")
	}
	if config.UseBestDRGPrice {
		headers.Add("use-best-drg-price", "true")
	}
	return headers
}
