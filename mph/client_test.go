package mph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mypricehealth/sling"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient("test")
	assert.NotNil(t, client.sling)
}

func TestClient(t *testing.T) {
	doSuccess := &fakeDoer{Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}}
	doFail := &fakeDoer{Response: &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(""))}, Error: fmt.Errorf("error")}

	clientSuccess := NewClient(doSuccess, false, "test")
	assert.NotNil(t, clientSuccess.sling)
	clientFail := NewClient(doFail, false, "test")
	assert.NotNil(t, clientFail.sling)
	clientTestFail := NewClient(doFail, true, "test")
	assert.NotNil(t, clientTestFail.sling)

	expectedRequestHeader := http.Header{}
	expectedRequestHeader.Set("Content-Type", "application/json")
	expectedRequestHeader.Set("x-api-key", "test")

	expectedAllOptionsHeder := http.Header{}
	expectedAllOptionsHeder.Set("Content-Type", "application/json")
	expectedAllOptionsHeder.Set("x-api-key", "test")
	expectedAllOptionsHeder.Set("is-commercial", "true")
	expectedAllOptionsHeder.Set("disable-cost-based-reimbursement", "true")
	expectedAllOptionsHeder.Set("use-commercial-synthetic-for-not-allowed", "true")
	expectedAllOptionsHeder.Set("override-threshold", "300")
	expectedAllOptionsHeder.Set("include-edits", "true")
	expectedAllOptionsHeder.Set("use-drg-from-grouper", "true")
	expectedAllOptionsHeder.Set("use-best-drg-price", "true")

	// Price TEST environment fail
	expectedRequest := newRequest("POST", "https://api-test.mypricehealth.com/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientTestFail.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[0])

	// Price TEST environment fail
	expectedRequest = newRequest("POST", "https://api-test.mypricehealth.com/v1/medicare/price/claim", Claim{}, expectedAllOptionsHeder)
	clientTestFail.Price(context.Background(), PriceConfig{IsCommercial: true, DisableCostBasedReimbursement: true, UseCommercialSyntheticForNotAllowed: true, UseDRGFromGrouper: true, UseBestDRGPrice: true, OverrideThreshold: 300, IncludeEdits: true}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[1])

	// Price success
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientSuccess.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[0])

	// Price fail
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/price/claim", Claim{}, expectedRequestHeader)
	clientFail.Price(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[2])

	// PriceBatch success
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/price/claims", []Claim{{}}, expectedRequestHeader)
	clientSuccess.PriceBatch(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[1])

	// PriceBatch fail
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/price/claims", []Claim{{}}, expectedRequestHeader)
	clientFail.PriceBatch(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doFail.RequestsMade[3])

	// Estimate success
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/estimate/claim", Claim{}, expectedRequestHeader)
	clientSuccess.Estimate(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[2])

	// EstimateBatch success
	expectedRequest = newRequest("POST", "https://api.mypricehealth.com/v1/medicare/estimate/claims", []Claim{{}}, expectedRequestHeader)
	clientSuccess.EstimateBatch(context.Background(), PriceConfig{}, Claim{})
	assertRequests(t, expectedRequest, doSuccess.RequestsMade[3])
}

func newRequest(method, url string, bodyStruct interface{}, headers http.Header) *http.Request {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(bodyStruct)
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
	assert.Nil(t, err)
	b2, err := io.ReadAll(actual)
	assert.Nil(t, err)
	assert.Equal(t, string(b1), string(b2))
}

type fakeDoer struct {
	Response     *http.Response
	Error        error
	RequestsMade []*http.Request
}

func (f *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	f.RequestsMade = append(f.RequestsMade, req)
	return f.Response, f.Error
}

var _ sling.Doer = &fakeDoer{}
