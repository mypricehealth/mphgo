package mph

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"braces.dev/errtrace"
	"github.com/mypricehealth/sling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultClient(t *testing.T) {
	t.Parallel()
	client := NewDefaultClient("test")
	assert.NotNil(t, client.sling)
}

func TestClient(t *testing.T) {
	t.Parallel()
	doSuccess := &fakeDoer{Response: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}}
	doFail := &fakeDoer{Response: &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader(""))}, Error: errtrace.Errorf("error")}

	clientSuccess := NewClient(doSuccess, false, "test")
	assert.NotNil(t, clientSuccess.sling)
	clientFail := NewClient(doFail, false, "test")
	assert.NotNil(t, clientFail.sling)
	clientTestFail := NewClient(doFail, true, "test")
	assert.NotNil(t, clientTestFail.sling)

	expectedRequestHeader := http.Header{}
	expectedRequestHeader.Set("Content-Type", "application/json")
	expectedRequestHeader.Set("x-api-key", "test")

	expectedAllOptionsHeader := http.Header{}
	expectedAllOptionsHeader.Set("Content-Type", "application/json")
	expectedAllOptionsHeader.Set("x-api-key", "test")
	expectedAllOptionsHeader.Set("contract-ruleset", "abcd")
	expectedAllOptionsHeader.Set("price-zero-billed", "true")
	expectedAllOptionsHeader.Set("is-commercial", "true")
	expectedAllOptionsHeader.Set("disable-cost-based-reimbursement", "true")
	expectedAllOptionsHeader.Set("use-commercial-synthetic-for-not-allowed", "true")
	expectedAllOptionsHeader.Set("use-drg-from-grouper", "true")
	expectedAllOptionsHeader.Set("use-best-drg-price", "true")
	expectedAllOptionsHeader.Set("override-threshold", "300")
	expectedAllOptionsHeader.Set("include-edits", "true")
	expectedAllOptionsHeader.Set("continue-on-edit-fail", "true")
	expectedAllOptionsHeader.Set("continue-on-provider-match-fail", "true")
	expectedAllOptionsHeader.Set("disable-machine-learning-estimates", "true")
	expectedAllOptionsHeader.Set("assume-impossible-anesthesia-units-are-minutes", "true")
	expectedAllOptionsHeader.Set("fallback-to-max-anesthesia-units-per-day", "true")
	expectedAllOptionsHeader.Set("allow-partial-results", "true")

	// Price TEST environment fail
	expectedRequest := newRequest("POST", "https://api-test.myprice.health/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientTestFail.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[0])

	// Price TEST environment fail
	expectedRequest = newRequest("POST", "https://api-test.myprice.health/v1/medicare/price/claim", Claim{}, expectedAllOptionsHeader)
	clientTestFail.Price(context.Background(), PriceConfig{
		ContractRuleset:                           "abcd",
		PriceZeroBilled:                           true,
		IsCommercial:                              true,
		DisableCostBasedReimbursement:             true,
		UseCommercialSyntheticForNotAllowed:       true,
		UseDRGFromGrouper:                         true,
		UseBestDRGPrice:                           true,
		OverrideThreshold:                         300,
		IncludeEdits:                              true,
		ContinueOnEditFail:                        true,
		ContinueOnProviderMatchFail:               true,
		DisableMachineLearningEstimates:           true,
		AssumeImpossibleAnesthesiaUnitsAreMinutes: true,
		FallbackToMaxAnesthesiaUnitsPerDay:        true,
		AllowPartialResults:                       true,
	}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[1])

	// Price success
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientSuccess.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[0])

	// Price fail
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientFail.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[2])

	// PriceBatch success
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/price/claims", []Claim{{}}, expectedRequestHeader)
	clientSuccess.PriceBatch(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[1])

	// PriceBatch fail
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/price/claims", []Claim{{}}, expectedRequestHeader)
	clientFail.PriceBatch(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[3])

	// Estimate success
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/estimate/rate-sheet", []RateSheet{{}}, expectedRequestHeader)
	clientSuccess.EstimateRateSheet(context.Background(), RateSheet{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[2])

	// EstimateBatch success
	expectedRequest = newRequest("POST", "https://api.myprice.health/v1/medicare/estimate/claims", []Claim{{}}, expectedRequestHeader)
	clientSuccess.EstimateClaims(context.Background(), Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[3])
}

func newRequest(method, url string, bodyStruct interface{}, headers http.Header) *http.Request {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(bodyStruct)
	body := io.NopCloser(&buf)
	req, _ := http.NewRequest(method, url, body)
	req.Header = headers
	return req
}

func assertRequests(t *testing.T, expected, actual *http.Request) {
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Method, actual.Method)
	assert.Equal(t, expected.Header, actual.Header)
	assertReaders(t, expected.Body, actual.Body)
}

func assertReaders(t *testing.T, expected, actual io.Reader) {
	b1, err := io.ReadAll(expected)
	require.NoError(t, err)
	b2, err := io.ReadAll(actual)
	require.NoError(t, err)
	assert.Equal(t, string(b1), string(b2))
}

type fakeDoer struct {
	Response     *http.Response
	Error        error
	RequestsMade []*http.Request
}

func (f *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	f.RequestsMade = append(f.RequestsMade, req)
	return f.Response, errtrace.Wrap(f.Error)
}

var _ sling.Doer = &fakeDoer{}

func TestGetHeaders(t *testing.T) {
	t.Parallel()
	input := PriceConfig{
		PriceZeroBilled:                     true,
		IsCommercial:                        true,
		DisableCostBasedReimbursement:       true,
		UseCommercialSyntheticForNotAllowed: true,
		OverrideThreshold:                   4,
		IncludeEdits:                        true,
		UseDRGFromGrouper:                   true,
		ContinueOnEditFail:                  true,
		ContinueOnProviderMatchFail:         true,
	}
	assert.Equal(t, fakeHeader(), GetHeaders(input))
}

func fakeHeader() http.Header {
	h := http.Header{}
	h.Add("price-zero-billed", "true")
	h.Add("is-commercial", "true")
	h.Add("disable-cost-based-reimbursement", "true")
	h.Add("use-commercial-synthetic-for-not-allowed", "true")
	h.Add("override-threshold", "4")
	h.Add("include-edits", "true")
	h.Add("use-drg-from-grouper", "true")
	h.Add("continue-on-edit-fail", "true")
	h.Add("continue-on-provider-match-fail", "true")
	return h
}

func TestParseHeaders(t *testing.T) {
	t.Parallel()
	input := &http.Request{Header: fakeHeader()}
	expected := PriceConfig{
		PriceZeroBilled:                     true,
		IsCommercial:                        true,
		DisableCostBasedReimbursement:       true,
		UseCommercialSyntheticForNotAllowed: true,
		OverrideThreshold:                   4,
		IncludeEdits:                        true,
		UseDRGFromGrouper:                   true,
		ContinueOnEditFail:                  true,
		ContinueOnProviderMatchFail:         true,
	}
	config, err := ParseHeaders(input)
	require.NoError(t, err)
	assert.Equal(t, expected, config)
}
