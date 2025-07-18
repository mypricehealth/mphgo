package mph

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"braces.dev/errtrace"
	"github.com/mypricehealth/sling"
)

// Pricer is the interface that wraps the Price and PriceBatch methods.
type Pricer interface {
	Price(ctx context.Context, config PriceConfig, input Claim) Response[Pricing]
	PriceBatch(ctx context.Context, config PriceConfig, inputs ...Claim) Responses[Pricing]
	EstimateClaims(ctx context.Context, inputs ...Claim) Responses[Pricing]
	EstimateRateSheet(ctx context.Context, inputs ...RateSheet) Responses[Pricing]
}

// PriceConfig is used to configure the behavior of the pricing API.
type PriceConfig struct {
	ContractRuleset                           string  `json:"contractRuleset,omitzero"`                           // set to the name of the ruleset to use for contract pricing
	PriceZeroBilled                           bool    `json:"priceZeroBilled,omitzero"`                           // set to true to price claims with zero billed amounts (default is false)
	IsCommercial                              bool    `json:"isCommercial,omitzero"`                              // set to true to use commercial code crosswalks
	DisableCostBasedReimbursement             bool    `json:"disableCostBasedReimbursement,omitzero"`             // by default, the API will use cost-based reimbursement for MAC priced line-items. This is the best estimate we have for this proprietary pricing
	UseCommercialSyntheticForNotAllowed       bool    `json:"useCommercialSyntheticForNotAllowed,omitzero"`       // set to true to use a synthetic Medicare price for line-items that are not allowed by Medicare
	UseDRGFromGrouper                         bool    `json:"useDRGFromGrouper,omitzero"`                         // set to true to always use the DRG from the inpatient grouper
	UseBestDRGPrice                           bool    `json:"useBestDRGPrice,omitzero"`                           // set to true to use the best DRG price between the price on the claim and the price from the grouper
	OverrideThreshold                         float64 `json:"overrideThreshold,omitzero"`                         // set to a value greater than 0 to allow the pricer flexibility to override NCCI edits and other overridable errors and return a price
	IncludeEdits                              bool    `json:"includeEdits,omitzero"`                              // set to true to include edit details in the response
	ContinueOnEditFail                        bool    `json:"continueOnEditFail,omitzero"`                        // set to true to continue to price the claim even if there are edit failures
	ContinueOnProviderMatchFail               bool    `json:"continueOnProviderMatchFail,omitzero"`               // set to true to continue with a average provider for the geographic area if the provider cannot be matched
	DisableMachineLearningEstimates           bool    `json:"disableMachineLearningEstimates,omitzero"`           // set to true to disable machine learning estimates (applies to estimates only)
	AssumeImpossibleAnesthesiaUnitsAreMinutes bool    `json:"assumeImpossibleAnesthesiaUnitsAreMinutes,omitzero"` // set to true to divide impossible anesthesia units by 15 (max of 96 anesthesia units per day) (default is false)
	FallbackToMaxAnesthesiaUnitsPerDay        bool    `json:"fallbackToMaxAnesthesiaUnitsPerDay,omitzero"`        // set to true to fallback to the maximum anesthesia units per day (default is false which will error if there are more than 96 anesthesia units per day)
	AllowPartialResults                       bool    `json:"allowPartialResults,omitzero"`                       // set to true to return partially repriced claims. This can be useful to get pricing on non-erroring line items, but should be used with caution
}

// Client is used to interact with the My Price Health API.
type Client struct {
	sling *sling.Sling
}

var _ Pricer = &Client{}

// NewClient is used to create a new API client for the My Price Health API. In most cases
// it is simpler to use NewDefaultClient to create a client with the default settings.
func NewClient(doer sling.Doer, isTest bool, apiKey string) *Client {
	url := "https://api.myprice.health"
	if isTest {
		url = "https://api-test.myprice.health"
	}
	client := &Client{sling.New().Doer(doer).Base(url).Set("x-api-key", apiKey)}
	return client
}

// NewDefaultClient is used to create a new API client for the My Price Health API with the default settings.
func NewDefaultClient(apiKey string) *Client {
	return NewClient(http.DefaultClient, false, apiKey)
}

func (c Client) receiveResponse(ctx context.Context, s *sling.Sling, path string) Response[Pricing] {
	var response Response[Pricing]
	res, err := s.Path(path).ReceiveWithContext(ctx, &response, &response)
	if err != nil {
		response.Error = &ResponseError{Title: fmt.Sprintf("fatal error calling %s", path), Detail: err.Error()}
		response.StatusCode = res.StatusCode
	}
	return response
}

func (c *Client) receiveResponses(ctx context.Context, s *sling.Sling, path string, count int) Responses[Pricing] {
	var responses Responses[Pricing]
	res, err := s.Path(path).ReceiveWithContext(ctx, &responses, &responses)
	if err != nil {
		responses.Error = &ResponseError{Title: fmt.Sprintf("fatal error calling %s", path), Detail: err.Error()}
		responses.ErrorCount = count
		responses.StatusCode = res.StatusCode
	}
	return responses
}

// EstimateRateSheet is used to get the estimated Medicare reimbursement of a single claim.
func (c *Client) EstimateRateSheet(ctx context.Context, inputs ...RateSheet) Responses[Pricing] {
	return c.receiveResponses(ctx, c.sling.BodyJSON(inputs).Method("POST"), "/v1/medicare/estimate/rate-sheet", len(inputs))
}

// EstimateClaims is used to get the estimated Medicare reimbursement of multiple claims.
func (c *Client) EstimateClaims(ctx context.Context, inputs ...Claim) Responses[Pricing] {
	return c.receiveResponses(ctx, c.sling.BodyJSON(inputs).Method("POST"), "/v1/medicare/estimate/claims", len(inputs))
}

// Price is used to get the Medicare reimbursement of a single claim.
func (c *Client) Price(ctx context.Context, config PriceConfig, input Claim) Response[Pricing] {
	return c.receiveResponse(ctx, c.sling.BodyJSON(input).AddHeaders(GetHeaders(config)).Method("POST"), "/v1/medicare/price/claim")
}

// PriceBatch is used to get the Medicare reimbursement of multiple claims.
func (c *Client) PriceBatch(ctx context.Context, config PriceConfig, inputs ...Claim) Responses[Pricing] {
	return c.receiveResponses(ctx, c.sling.BodyJSON(inputs).AddHeaders(GetHeaders(config)).Method("POST"), "/v1/medicare/price/claims", len(inputs))
}

func GetHeaders(config PriceConfig) http.Header {
	headers := http.Header{}
	if config.PriceZeroBilled {
		headers.Add("price-zero-billed", "true")
	}
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
	if config.ContinueOnEditFail {
		headers.Add("continue-on-edit-fail", "true")
	}
	if config.ContinueOnProviderMatchFail {
		headers.Add("continue-on-provider-match-fail", "true")
	}
	if config.UseDRGFromGrouper {
		headers.Add("use-drg-from-grouper", "true")
	}
	if config.UseBestDRGPrice {
		headers.Add("use-best-drg-price", "true")
	}
	if config.AssumeImpossibleAnesthesiaUnitsAreMinutes {
		headers.Add("assume-impossible-anesthesia-units-are-minutes", "true")
	}
	if config.FallbackToMaxAnesthesiaUnitsPerDay {
		headers.Add("fallback-to-max-anesthesia-units-per-day", "true")
	}
	if config.DisableMachineLearningEstimates {
		headers.Add("disable-machine-learning-estimates", "true")
	}
	if config.AllowPartialResults {
		headers.Add("allow-partial-results", "true")
	}
	if config.ContractRuleset != "" {
		headers.Add("contract-ruleset", config.ContractRuleset)
	}
	return headers
}

func ParseHeaders(r *http.Request) (PriceConfig, error) {
	overrideThreshold, _ := strconv.ParseFloat(r.Header.Get("override-threshold"), 64)
	config := PriceConfig{
		PriceZeroBilled:                           r.Header.Get("price-zero-billed") == "true",
		IsCommercial:                              r.Header.Get("is-commercial") == "true",
		DisableCostBasedReimbursement:             r.Header.Get("disable-cost-based-reimbursement") == "true",
		UseCommercialSyntheticForNotAllowed:       r.Header.Get("use-commercial-synthetic-for-not-allowed") == "true",
		OverrideThreshold:                         overrideThreshold,
		IncludeEdits:                              r.Header.Get("include-edits") == "true",
		UseDRGFromGrouper:                         r.Header.Get("use-drg-from-grouper") == "true",
		UseBestDRGPrice:                           r.Header.Get("use-best-drg-price") == "true",
		ContinueOnEditFail:                        r.Header.Get("continue-on-edit-fail") == "true",
		ContinueOnProviderMatchFail:               r.Header.Get("continue-on-provider-match-fail") == "true",
		AssumeImpossibleAnesthesiaUnitsAreMinutes: r.Header.Get("assume-impossible-anesthesia-units-are-minutes") == "true",
		FallbackToMaxAnesthesiaUnitsPerDay:        r.Header.Get("fallback-to-max-anesthesia-units-per-day") == "true",
		AllowPartialResults:                       r.Header.Get("allow-partial-results") == "true",
		ContractRuleset:                           r.Header.Get("contract-ruleset"),
	}
	if config.UseDRGFromGrouper && config.UseBestDRGPrice {
		return PriceConfig{}, errtrace.Errorf("use-drg-from-grouper and use-best-drg-price are mutually exclusive")
	}
	return config, nil
}
