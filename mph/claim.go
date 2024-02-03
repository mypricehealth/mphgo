package mph

import (
	"github.com/shopspring/decimal"
)

type FormType string         // Type of form used to submit the claim. Can be HCFA or UB-04
type BillTypeSequence string // The location where the claim is at in its billing lifecycle (e.g. 0: Non-Pay, 1: Admit Through Discharge, 7: Replacement, etc.)
type SexType uint8           // Biological sex of the patient for clinical purposes

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
	Provider                            // Provider information for the claim
	ClaimID            string           `json:"claimID,omitempty"`            // Unique identifier for the claim (from REF D9)
	PlanCode           string           `json:"planCode,omitempty"`           // Identifies the subscriber's plan (from SBR03)
	PatientSex         SexType          `json:"patientSex,omitempty"`         // Biological sex of the patient for clinical purposes (from DMG02). 0:Unknown, 1:Male, 2:Female
	PatientDateOfBirth *Date            `json:"patientDateOfBirth,omitempty"` // Patient date of birth (from DMG03)
	PatientHeightInCM  float64          `json:"patientHeightInCM,omitempty"`  // Patient height in centimeters (from HI value A9, MEA value HT)
	PatientWeightInKG  float64          `json:"patientWeightInKG,omitempty"`  // Patient weight in kilograms (from HI value A8, PAT08, CR102 [ambulance only])
	AmbulancePickupZIP string           `json:"ambulancePickupZIP,omitempty"` // Location where patient was picked up in ambulance (from HI with HIxx_01=BE and HIxx_02=A0 or NM1 loop with NM1 PW)
	FormType           FormType         `json:"formType,omitempty"`           // Type of form used to submit the claim. Can be HCFA or UB-04 (from CLM05_02)
	BillTypeOrPOS      string           `json:"billTypeOrPOS,omitempty"`      // Describes type of facility where services were rendered (from CLM05_01)
	BillTypeSequence   BillTypeSequence `json:"billTypeSequence,omitempty"`   // Where the claim is at in its billing lifecycle (e.g. 0: Non-Pay, 1: Admit Through Discharge, 7: Replacement, etc.) (from CLM05_03)
	BilledAmount       float64          `json:"billedAmount,omitempty"`       // Billed amount from provider (from CLM02)
	AllowedAmount      float64          `json:"allowedAmount,omitempty"`      // Amount allowed by the plan for payment. Both member and plan responsibility (non-EDI)
	PaidAmount         float64          `json:"paidAmount,omitempty"`         // Amount paid by the plan for the claim (non-EDI)
	DateFrom           Date             `json:"dateFrom,omitempty"`           // Earliest service date among services, or statement date if not found
	DateThrough        Date             `json:"dateThrough,omitempty"`        // Latest service date among services, or statement date if not found
	DischargeStatus    string           `json:"dischargeStatus,omitempty"`    // Status of the patient at time of discharge (from CL103)
	AdmitDiagnosis     string           `json:"admitDiagnosis,omitempty"`     // ICD diagnosis at the time the patient was admitted (from HI ABJ or BJ)
	PrincipalDiagnosis *Diagnosis       `json:"principalDiagnosis,omitempty"` // Principal ICD diagnosis for the patient (from HI ABK or BK)
	OtherDiagnoses     []Diagnosis      `json:"otherDiagnoses,omitempty"`     // Other ICD diagnoses that apply to the patient (from HI ABF or BF)
	PrincipalProcedure string           `json:"principalProcedure,omitempty"` // Principal ICD procedure for the patient (from HI BBR or BR)
	OtherProcedures    []string         `json:"otherProcedures,omitempty"`    // Other ICD procedures that apply to the patient (from HI BBQ or BQ)
	ConditionCodes     []string         `json:"conditionCodes,omitempty"`     // Special conditions that may affect payment or other processing (from HI BG)
	ValueCodes         []ValueCode      `json:"valueCodes,omitempty"`         // Numeric values related to the patient or claim (HI BE)
	OccurrenceCodes    []string         `json:"occurrenceCodes,omitempty"`    // Date related occurrences related to the patient or claim (from HI BH)
	DRG                string           `json:"drg,omitempty"`                // Diagnosis Related Group for inpatient services (from HI DR)
	Services           []Service        `json:"services,omitempty"`           // One or more services provided to the patient (from LX loop)
}

type Provider struct {
	NPI                      string   `json:"npi,omitempty"`                      // National Provider Identifier of the provider (from NM109, required)
	ProviderTaxID            string   `json:"providerTaxID,omitempty"`            // Tax ID of the provider (from REF highly recommended)
	ProviderPhones           []string `json:"providerPhones,omitempty"`           // Phone numbers of the provider (from PER, optional)
	ProviderFaxes            []string `json:"providerFaxes,omitempty"`            // Fax numbers of the provider (from PER, optional)
	ProviderEmails           []string `json:"providerEmails,omitempty"`           // Email addresses of the provider (from PER, optional)
	ProviderLicenseNumber    string   `json:"providerLicenseNumber,omitempty"`    // State license number of the provider (from REF 0B, optional)
	ProviderCommercialNumber string   `json:"providerCommercialNumber,omitempty"` // Commercial number of the provider used by some payers (from REF G2, optional)
	ProviderTaxonomy         string   `json:"providerTaxonomy,omitempty"`         // Taxonomy code of the provider (from PRV03, highly recommended)
	ProviderFirstName        string   `json:"providerFirstName,omitempty"`        // First name of the provider (NM104, highly recommended)
	ProviderLastName         string   `json:"providerLastName,omitempty"`         // Last name of the provider (from NM103, highly recommended)
	ProviderOrgName          string   `json:"providerOrgName,omitempty"`          // Organization name of the provider (from NM103, highly recommended)
	ProviderAddress1         string   `json:"providerAddress1,omitempty"`         // Address line 1 of the provider (from N301, highly recommended)
	ProviderAddress2         string   `json:"providerAddress2,omitempty"`         // Address line 2 of the provider (from N302, optional)
	ProviderCity             string   `json:"providerCity,omitempty"`             // City of the provider (from N401, highly recommended)
	ProviderState            string   `json:"providerState,omitempty"`            // State of the provider (from N402, highly recommended)
	ProviderZIP              string   `json:"providerZIP,omitempty"`              // ZIP code of the provider (from N403, required)
}

type Diagnosis struct { // Principal, Other Diagnosis, Admitting Diagnosis, External Cause of Injury
	Code               string `json:"code,omitempty"               fixed:"1,9"`   // ICD-10 diagnosis code (from HIxx_02)
	PresentOnAdmission string `json:"presentOnAdmission,omitempty" fixed:"10,13"` // Flag indicates whether diagnosis was present at the time of admission (from HIxx_09)
}

type ValueCode struct {
	Code   string          `json:"code,omitempty"   fixed:"1,2"`        // Code indicating the type of value provided (from HIxx_02)
	Amount decimal.Decimal `json:"amount,omitempty" fixed:"3,11,right"` // Amount associated with the value code (from HIxx_05)
}

type Service struct {
	Provider                    // Additional provider information specific to this service item
	LineNumber         string   `json:"lineNumber,omitempty"`         // Unique line number for the service item (from LX01)
	RevCode            string   `json:"revCode,omitempty"`            // Revenue code (from SV2_01)
	ProcedureCode      string   `json:"procedureCode,omitempty"`      // Procedure code (from SV101_02 / SV202_02)
	ProcedureModifiers []string `json:"procedureModifiers,omitempty"` // Procedure modifiers (from SV101_03, 4, 5, 6 / SV202_03, 4, 5, 6)
	DrugCode           string   `json:"drugCode,omitempty"`           // National Drug Code (from LIN03)
	DateFrom           Date     `json:"dateFrom,omitempty"`           // Begin date of service (from DTP 472)
	DateThrough        Date     `json:"dateThrough,omitempty"`        // End date of service (from DTP 472)
	BilledAmount       float64  `json:"billedAmount,omitempty"`       // Billed charge for the service (from SV102 / SV203)
	AllowedAmount      float64  `json:"allowedAmount,omitempty"`      // Plan allowed amount for the service (non-EDI)
	PaidAmount         float64  `json:"paidAmount,omitempty"`         // Plan paid amount for the service (non-EDI)
	Quantity           float64  `json:"quantity"`                     // Quantity of the service (from SV104 / SV205)
	Units              string   `json:"units,omitempty"`              // Units connected to the quantity given (from SV103 / SV204)
	PlaceOfService     string   `json:"placeOfService,omitempty"`     // Place of service code (from SV105)
	DiagnosisPointers  []int8   `json:"diagnosisPointers,omitempty"`  // Diagnosis pointers (from SV107)
	AmbulancePickupZIP string   `json:"ambulancePickupZIP,omitempty"` // ZIP code where ambulance picked up patient. Supplied if different than claim-level value (from NM1 PW)
}
