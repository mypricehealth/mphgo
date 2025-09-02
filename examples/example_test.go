package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"braces.dev/errtrace"
	"github.com/mypricehealth/mphgo/mph"
	"github.com/stretchr/testify/assert"
)

var config = mph.PriceConfig{
	IsCommercial:                        true,  // uses commercial code crosswalks
	DisableCostBasedReimbursement:       false, // use cost-based reimbursement for MAC priced line-items
	UseCommercialSyntheticForNotAllowed: true,  // use synthetic Medicare for line items not allowed by Medicare, but which may still be paid by commercial plans
	UseDRGFromGrouper:                   false, // always use the DRG from the inpatient grouper (not applicable with UseBestDRGPrice set to true)
	UseBestDRGPrice:                     true,  // price both using the DRG supplied in the claim and the DRG from the grouper and return the lowest price
	OverrideThreshold:                   300,   // for claims which fail NCCI or other edit rules, override the errors up to this amount to get a price
	IncludeEdits:                        true,  // get detailed information from the code editor about why a claim failed
}

func TestClientWithJSON(t *testing.T) {
	t.SkipNow()

	c := mph.NewDefaultClient("apiKey") // replace this with your API key
	inpatientClaim, err := readJSON("testdata/inpatient.json")
	assert.Nil(t, err)
	outpatientClaim, err := readJSON("testdata/outpatient.json")
	assert.Nil(t, err)
	hcfaClaim, err := readJSON("testdata/hcfa.json")
	assert.Nil(t, err)

	result := c.Price(context.Background(), config, inpatientClaim)
	assert.Nil(t, result.Error)
	fmt.Println(result.Result.MedicareAmount)

	result = c.Price(context.Background(), config, outpatientClaim)
	assert.Nil(t, result.Error)
	fmt.Println(result.Result.MedicareAmount)

	result = c.Price(context.Background(), config, hcfaClaim)
	assert.Nil(t, result.Error)
	fmt.Println(t, result.Result.MedicareAmount)
}

func readJSON(filename string) (mph.Claim, error) {
	var c mph.Claim
	f, err := os.Open(filename)
	if err != nil {
		return c, errtrace.Wrap(err)
	}

	err = json.NewDecoder(f).Decode(&c)
	return c, errtrace.Wrap(err)
}

func TestClientConstructingStructs(t *testing.T) {
	t.SkipNow()

	c := mph.NewDefaultClient("apiKey") // replace this with your API key
	result := c.Price(context.Background(), config, inpatientClaim)
	assert.Nil(t, result.Error)
	fmt.Println(result.Result.MedicareAmount)

	result = c.Price(context.Background(), config, outpatientClaim)
	assert.Nil(t, result.Error)
	fmt.Println(result.Result.MedicareAmount)
}

// fake inpatient claim for testing purposes
var inpatientClaim = mph.Claim{
	Provider: mph.Provider{
		NPI:         "1962999664",
		ProviderZIP: "35960",
	},
	DRG:                "461",
	PatientDateOfBirth: mph.NewDatePtr(1988, 1, 2),
	FormType:           mph.UBFormType,
	BillTypeOrPOS:      "111",
	BilledAmount:       47000,
	DateFrom:           mph.NewDate(2020, 2, 27),
	DateThrough:        mph.NewDate(2020, 2, 28),
	PrincipalDiagnosis: &mph.Diagnosis{Code: "N186"},
	OtherDiagnoses: []mph.Diagnosis{
		{Code: "Z992"},
		{Code: "I120"},
		{Code: "E6601"},
		{Code: "E785"},
		{Code: "Z6832"},
	},
	Services: []mph.Service{
		{LineNumber: "1", RevCode: "320", BilledAmount: 2126, DateFrom: mph.NewDate(2020, 2, 27), DateThrough: mph.NewDate(2020, 2, 27), ProcedureCode: "76000", Quantity: 1},
		{LineNumber: "2", RevCode: "360", BilledAmount: 28684, DateFrom: mph.NewDate(2020, 2, 27), DateThrough: mph.NewDate(2020, 2, 27), ProcedureCode: "36821", Quantity: 1},
		{LineNumber: "3", RevCode: "370", BilledAmount: 16414, DateFrom: mph.NewDate(2020, 2, 27), DateThrough: mph.NewDate(2020, 2, 27), ProcedureCode: "", Quantity: 48},
	},
}

var date = mph.Date{Time: time.Date(2020, 9, 11, 0, 0, 0, 0, time.UTC)}

// fake outpatient claim for testing purposes
var outpatientClaim = mph.Claim{
	Provider: mph.Provider{
		NPI:              "1164403861",
		ProviderOrgName:  "SOUTHEAST ALABAMA MEDICAL CENTER",
		ProviderAddress1: "1108 ROSS CLARK CIRCLE",
		ProviderCity:     "DOTHAN",
		ProviderState:    "AL",
		ProviderZIP:      "36301",
		ProviderTaxonomy: "282N00000X",
	},
	PatientSex:         1,
	PatientDateOfBirth: &mph.Date{Time: time.Date(1926, 11, 11, 0, 0, 0, 0, time.UTC)},
	BilledAmount:       21000,
	FormType:           "UB-04",
	BillTypeOrPOS:      "13",
	BillTypeSequence:   "1",
	DateFrom:           date,
	DateThrough:        date,
	DischargeStatus:    "01",
	PrincipalDiagnosis: &mph.Diagnosis{Code: "Z0001"},
	OtherDiagnoses: []mph.Diagnosis{
		{Code: "Z13220"},
		{Code: "I10"},
	},
	OccurrenceCodes: []string{"A1"},
	Services: []mph.Service{
		{
			LineNumber:    "1",
			RevCode:       "0301",
			ProcedureCode: "80053",
			BilledAmount:  312,
			Units:         "UN",
			Quantity:      1,
			DateFrom:      date,
			DateThrough:   date,
		},
		{
			LineNumber:    "2",
			RevCode:       "0301",
			ProcedureCode: "80061",
			BilledAmount:  298,
			Units:         "UN",
			Quantity:      1,
			DateFrom:      date,
			DateThrough:   date,
		},
		{
			LineNumber:    "3",
			RevCode:       "0301",
			ProcedureCode: "84443",
			BilledAmount:  51,
			Units:         "UN",
			Quantity:      1,
			DateFrom:      date,
			DateThrough:   date,
		},
	},
}
