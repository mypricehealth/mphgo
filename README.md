# mphgo

[![Go Reference](https://pkg.go.dev/badge/github.com/mypricehealth/mphgo.svg)](https://pkg.go.dev/github.com/mypricehealth/mphgo)

<img src="https://raw.githubusercontent.com/mypricehealth/mphgo/main/gopherdoc.png">

A Go client library to make it easy to connect to the My Price Health API. It can be accessed at https://api.myprice.health. This client library makes it easy to connect to the API for Go developers.

## The My Price Health API

The My Price Health API can be used to get pricing and estimated pricing for Medicare reimbursement. Medicare benchmarks are a great tool to advance the quadruple aim in healthcare.

## Usage

See also the examples folder for additional examples

```go

import (
	"context"
	"fmt"

	"github.com/mypricehealth/mphgo/mph"
)

func main() {
	var config = mph.PriceConfig{
		IsCommercial:                        true,  // uses commercial code crosswalks
		DisableCostBasedReimbursement:       false, // use cost-based reimbursement for MAC priced line-items
		UseCommercialSyntheticForNotAllowed: true,  // use synthetic Medicare for line items not allowed by Medicare, but which may still be paid by commercial plans
		UseDRGFromGrouper:                   false, // always use the DRG from the inpatient grouper (not applicable with UseBestDRGPrice set to true)
		UseBestDRGPrice:                     true,  // price both using the DRG supplied in the claim and the DRG from the grouper and return the lowest price
		OverrideThreshold:                   300,   // for claims which fail NCCI or other edit rules, override the errors up to this amount to get a price
		IncludeEdits:                        true,  // get detailed information from the code editor about why a claim failed
	}

	c := mph.NewDefaultClient("apiKey") // replace this with your API key
	result := c.Price(context.Background(), config, inpatientClaim)
	if result.Error != nil {
		fmt.Println(result.Error)
	} else {
		fmt.Println(result.Result.MedicareAmount)
	}
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
	BilledAmount:       47224,
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
```

## Why Medicare Pricing?

It is possible and practical to achieve the quadruple aim in healthcare. With Medicare pricing for all your claims data, you’ll have the tools you need to:

- Lower cost through better provider negotiation, better plan design and more engaged members
- Better outcomes through an effective health plan prioritizing quality care, low cost and better health management
- Happier plan members through education and a better health plan
- More engaged providers who act as partners, not adversaries

For more information about how Medicare pricing can advance the quadruple aim in healthcare, see our [Advancing the quadruple aim with Medicare pricing white paper](https://myprice.health/Advancing%20the%20quadruple%20aim%20with%20Medicare%20pricing%20-%20v2.pdf).
