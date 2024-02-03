package mph

import (
	"github.com/shopspring/decimal"
)

// Here are the structs we use to represent the input of our API
// By converting into these structs, you should be able to marshal
// into JSON that our API will accept
type FormType string
type BillTypeSequence string
type PaymentType string
type SexType uint8

var (
	UBFormType   FormType = "UB-04"
	HCFAFormType FormType = "HCFA"
)

var (
	NonPayBillTypeSequence                 BillTypeSequence = "0"
	AdmitThroughDischargeBillTypeSequence  BillTypeSequence = "1"
	FirstInterimBillTypeSequence           BillTypeSequence = "2"
	ContinuingInterimBillTypeSequence      BillTypeSequence = "3"
	LastInterimBillTypeSequence            BillTypeSequence = "4"
	LateChargeBillTypeSequence             BillTypeSequence = "5"
	FirstInterimBillTypeSequenceDeprecated BillTypeSequence = "6"
	ReplacementBillTypeSequence            BillTypeSequence = "7"
	VoidOrCancelBillTypeSequence           BillTypeSequence = "8"
	FinalClaimBillTypeSequence             BillTypeSequence = "9"
	CWFAdjustmentBillTypeSequence          BillTypeSequence = "G"
	CMSAdjustmentBillTypeSequence          BillTypeSequence = "H"
	IntermediaryAdjustmentBillTypeSequence BillTypeSequence = "I"
	OtherAdjustmentBillTypeSequence        BillTypeSequence = "J"
	OIGAdjustmentBillTypeSequence          BillTypeSequence = "K"
	MSPAdjustmentBillTypeSequence          BillTypeSequence = "M"
	QIOAdjustmentBillTypeSequence          BillTypeSequence = "P"
	ProviderAdjustmentBillTypeSequence     BillTypeSequence = "Q"
	SexTypeUnknown                         SexType          = 0
	SexTypeMale                            SexType          = 1
	SexTypeFemale                          SexType          = 2
)

type Claim struct {
	Provider
	ClaimID            string           `json:"claimID,omitempty"`
	PlanCode           string           `json:"planCode,omitempty"`
	PatientSex         SexType          `json:"patientSex,omitempty"`         // DMG02 (0:Unknown, 1:Male, 2:Female)
	PatientDateOfBirth *Date            `json:"patientDateOfBirth,omitempty"` // DMG03
	PatientHeightInCM  float64          `json:"patientHeightInCM,omitempty"`  // HI value A9, MEA value HT
	PatientWeightInKG  float64          `json:"patientWeightInKG,omitempty"`  // HI value A8, PAT08, CR102 (ambulance only)
	AmbulancePickupZIP string           `json:"ambulancePickupZIP,omitempty"` // HI with HIxx_01=BE and HIxx_02=A0 or NM1 loop with NM1 PW
	FormType           FormType         `json:"formType,omitempty"`           // CLM05_02
	BillTypeOrPOS      string           `json:"billTypeOrPOS,omitempty"`      // CLM05_01
	BillTypeSequence   BillTypeSequence `json:"billTypeSequence,omitempty"`   // CLM05_03
	BilledAmount       float64          `json:"billedAmount,omitempty"`       // CLM02
	AllowedAmount      float64          `json:"allowedAmount,omitempty"`      // plan allowed
	PaidAmount         float64          `json:"paidAmount,omitempty"`         // plan paid
	DateFrom           Date             `json:"dateFrom,omitempty"`           // Earliest service date among services, or statement date if not found
	DateThrough        Date             `json:"dateThrough,omitempty"`        // Latest service date among services, or statement date if not found
	DischargeStatus    string           `json:"dischargeStatus,omitempty"`    // CL103
	AdmitDiagnosis     string           `json:"admitDiagnosis,omitempty"`     // HI segment
	PrincipalDiagnosis *Diagnosis       `json:"principalDiagnosis,omitempty"` // HI segment
	OtherDiagnoses     []Diagnosis      `json:"otherDiagnoses,omitempty"`     // HI segment
	PrincipalProcedure string           `json:"principalProcedure,omitempty"` // HI segment
	OtherProcedures    []string         `json:"otherProcedures,omitempty"`    // HI segment
	ConditionCodes     []string         `json:"conditionCodes,omitempty"`     // HI segment
	ValueCodes         []ValueCode      `json:"valueCodes,omitempty"`         // HI segment
	OccurrenceCodes    []string         `json:"occurrenceCodes,omitempty"`    // HI segment
	DRG                string           `json:"drg,omitempty"`                // HI segment
	Services           []Service        `json:"services,omitempty"`
}

type Provider struct {
	NPI                      string   `json:"npi,omitempty"`
	ProviderTaxID            string   `json:"providerTaxID,omitempty"`
	ProviderPhones           []string `json:"providerPhones,omitempty"`
	ProviderFaxes            []string `json:"providerFaxes,omitempty"`
	ProviderEmails           []string `json:"providerEmails,omitempty"`
	ProviderLicenseNumber    string   `json:"providerLicenseNumber,omitempty"`
	ProviderCommercialNumber string   `json:"providerCommercialNumber,omitempty"`
	ProviderTaxonomy         string   `json:"providerTaxonomy,omitempty"`
	ProviderFirstName        string   `json:"providerFirstName,omitempty"`
	ProviderLastName         string   `json:"providerLastName,omitempty"`
	ProviderOrgName          string   `json:"providerOrgName,omitempty"`
	ProviderAddress1         string   `json:"providerAddress1,omitempty"`
	ProviderAddress2         string   `json:"providerAddress2,omitempty"`
	ProviderCity             string   `json:"providerCity,omitempty"`
	ProviderState            string   `json:"providerState,omitempty"`
	ProviderZIP              string   `json:"providerZIP,omitempty"`
}

type Diagnosis struct { // Principal, Other Diagnosis, Admitting Diagnosis, External Cause of Injury
	Code               string `json:"code,omitempty"               fixed:"1,9"`   // HI01_02
	PresentOnAdmission string `json:"presentOnAdmission,omitempty" fixed:"10,13"` // HI01_09
}

type ValueCode struct {
	Code   string          `json:"code,omitempty"   fixed:"1,2"`        // HIxx_02
	Amount decimal.Decimal `json:"amount,omitempty" fixed:"3,11,right"` // HIxx_05
}

type Service struct {
	Provider
	LineNumber         string   `json:"lineNumber,omitempty"`
	RevCode            string   `json:"revCode,omitempty"`
	ProcedureCode      string   `json:"procedureCode,omitempty"`
	ProcedureModifiers []string `json:"procedureModifiers,omitempty"`
	DrugCode           string   `json:"drugCode,omitempty"`
	DateFrom           Date     `json:"dateFrom,omitempty"`
	DateThrough        Date     `json:"dateThrough,omitempty"`
	BilledAmount       float64  `json:"billedAmount,omitempty"`
	AllowedAmount      float64  `json:"allowedAmount,omitempty"`
	PaidAmount         float64  `json:"paidAmount,omitempty"`
	Quantity           float64  `json:"quantity"`
	Units              string   `json:"units,omitempty"`
	PlaceOfService     string   `json:"placeOfService,omitempty"`
	DiagnosisPointers  []int8   `json:"diagnosisPointers,omitempty"`
	AmbulancePickupZIP string   `json:"ambulancePickupZIP,omitempty"` // may override the claim-level value
}
