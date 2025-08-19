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
)

var (
	SexTypeUnknown SexType = 0
	SexTypeMale    SexType = 1
	SexTypeFemale  SexType = 2
)

type Claim struct {
	Provider                            // Provider information for the claim
	ClaimID            string           `json:"claimID,omitzero"`            // Unique identifier for the claim (from REF D9)
	PlanCode           string           `json:"planCode,omitzero"`           // Identifies the subscriber's plan (from SBR03)
	PatientSex         SexType          `json:"patientSex,omitzero"`         // Biological sex of the patient for clinical purposes (from DMG02). 0:Unknown, 1:Male, 2:Female
	PatientDateOfBirth *Date            `json:"patientDateOfBirth,omitzero"` // Patient date of birth (from DMG03)
	PatientHeightInCM  float64          `json:"patientHeightInCM,omitzero"`  // Patient height in centimeters (from HI value A9, MEA value HT)
	PatientWeightInKG  float64          `json:"patientWeightInKG,omitzero"`  // Patient weight in kilograms (from HI value A8, PAT08, CR102 [ambulance only])
	AmbulancePickupZIP string           `json:"ambulancePickupZIP,omitzero"` // Location where patient was picked up in ambulance (from HI with HIxx_01=BE and HIxx_02=A0 or NM1 loop with NM1 PW)
	FormType           FormType         `json:"formType,omitzero"`           // Type of form used to submit the claim. Can be HCFA or UB-04 (from CLM05_02)
	BillTypeOrPOS      string           `json:"billTypeOrPOS,omitzero"`      // Describes type of facility where services were rendered (from CLM05_01)
	BillTypeSequence   BillTypeSequence `json:"billTypeSequence,omitzero"`   // Where the claim is at in its billing lifecycle (e.g. 0: Non-Pay, 1: Admit Through Discharge, 7: Replacement, etc.) (from CLM05_03)
	BilledAmount       float64          `json:"billedAmount,omitzero"`       // Billed amount from provider (from CLM02)
	AllowedAmount      float64          `json:"allowedAmount,omitzero"`      // Amount allowed by the plan for payment. Both member and plan responsibility (non-EDI)
	PaidAmount         float64          `json:"paidAmount,omitzero"`         // Amount paid by the plan for the claim (non-EDI)
	DateFrom           Date             `json:"dateFrom,omitzero"`           // Earliest service date among services, or statement date if not found
	DateThrough        Date             `json:"dateThrough,omitzero"`        // Latest service date among services, or statement date if not found
	DischargeStatus    string           `json:"dischargeStatus,omitzero"`    // Status of the patient at time of discharge (from CL103)
	AdmitDiagnosis     string           `json:"admitDiagnosis,omitzero"`     // ICD diagnosis at the time the patient was admitted (from HI ABJ or BJ)
	PrincipalDiagnosis *Diagnosis       `json:"principalDiagnosis,omitzero"` // Principal ICD diagnosis for the patient (from HI ABK or BK)
	OtherDiagnoses     []Diagnosis      `json:"otherDiagnoses,omitempty"`    // Other ICD diagnoses that apply to the patient (from HI ABF or BF)
	PrincipalProcedure string           `json:"principalProcedure,omitzero"` // Principal ICD procedure for the patient (from HI BBR or BR)
	OtherProcedures    []string         `json:"otherProcedures,omitempty"`   // Other ICD procedures that apply to the patient (from HI BBQ or BQ)
	ConditionCodes     []string         `json:"conditionCodes,omitempty"`    // Special conditions that may affect payment or other processing (from HI BG)
	ValueCodes         []ValueCode      `json:"valueCodes,omitempty"`        // Numeric values related to the patient or claim (HI BE)
	OccurrenceCodes    []string         `json:"occurrenceCodes,omitempty"`   // Date related occurrences related to the patient or claim (from HI BH)
	DRG                string           `json:"drg,omitzero"`                // Diagnosis Related Group for inpatient services (from HI DR)
	Services           []Service        `json:"services,omitempty"`          // One or more services provided to the patient (from LX loop)
}

// Provider represents the service provider that rendered healthcare services on behalf of the patient.
// This can be found in Loop 2000A and/or Loop 2310 NM101-77 at the claim level, and may also be overridden
// at the service level in the 2400 loop.
type Provider struct {
	NPI                      string   `json:"npi,omitzero"`                      // National Provider Identifier of the provider (from NM109, required)
	CCN                      string   `json:"ccn,omitzero"`                      // CMS Certification Number (optional)
	ProviderTaxID            string   `json:"providerTaxID,omitzero"`            // Tax ID of the provider (from REF highly recommended)
	ProviderPhones           []string `json:"providerPhones,omitzero"`           // Phone numbers of the provider (from PER, optional)
	ProviderFaxes            []string `json:"providerFaxes,omitzero"`            // Fax numbers of the provider (from PER, optional)
	ProviderEmails           []string `json:"providerEmails,omitzero"`           // Email addresses of the provider (from PER, optional)
	ProviderLicenseNumber    string   `json:"providerLicenseNumber,omitzero"`    // State license number of the provider (from REF 0B, optional)
	ProviderCommercialNumber string   `json:"providerCommercialNumber,omitzero"` // Commercial number of the provider used by some payers (from REF G2, optional)
	ProviderTaxonomy         string   `json:"providerTaxonomy,omitzero"`         // Taxonomy code of the provider (from PRV03, highly recommended)
	ProviderFirstName        string   `json:"providerFirstName,omitzero"`        // First name of the provider (NM104, highly recommended)
	ProviderLastName         string   `json:"providerLastName,omitzero"`         // Last name of the provider (from NM103, highly recommended)
	ProviderOrgName          string   `json:"providerOrgName,omitzero"`          // Organization name of the provider (from NM103, highly recommended)
	ProviderAddress1         string   `json:"providerAddress1,omitzero"`         // Address line 1 of the provider (from N301, highly recommended)
	ProviderAddress2         string   `json:"providerAddress2,omitzero"`         // Address line 2 of the provider (from N302, optional)
	ProviderCity             string   `json:"providerCity,omitzero"`             // City of the provider (from N401, highly recommended)
	ProviderState            string   `json:"providerState,omitzero"`            // State of the provider (from N402, highly recommended)
	ProviderZIP              string   `json:"providerZIP,omitzero"`              // ZIP code of the provider (from N403, required)
}

type Diagnosis struct { // Principal, Other Diagnosis, Admitting Diagnosis, External Cause of Injury
	Code               string `json:"code,omitzero"`               // ICD-10 diagnosis code (from HIxx_02)
	PresentOnAdmission string `json:"presentOnAdmission,omitzero"` // Flag indicates whether diagnosis was present at the time of admission (from HIxx_09)
}

type ValueCode struct {
	Code   string          `json:"code,omitzero"`   // Code indicating the type of value provided (from HIxx_02)
	Amount decimal.Decimal `json:"amount,omitzero"` // Amount associated with the value code (from HIxx_05)
}

type Service struct {
	Provider                    // Additional provider information specific to this service item
	LineNumber         string   `json:"lineNumber,omitzero"`         // Unique line number for the service item (from LX01)
	RevCode            string   `json:"revCode,omitzero"`            // Revenue code (from SV2_01)
	ProcedureCode      string   `json:"procedureCode,omitzero"`      // Procedure code (from SV101_02 / SV202_02)
	ProcedureModifiers []string `json:"procedureModifiers,omitzero"` // Procedure modifiers (from SV101_03, 4, 5, 6 / SV202_03, 4, 5, 6)
	DrugCode           string   `json:"drugCode,omitzero"`           // National Drug Code (from LIN03)
	DateFrom           Date     `json:"dateFrom,omitzero"`           // Begin date of service (from DTP 472)
	DateThrough        Date     `json:"dateThrough,omitzero"`        // End date of service (from DTP 472)
	BilledAmount       float64  `json:"billedAmount,omitzero"`       // Billed charge for the service (from SV102 / SV203)
	AllowedAmount      float64  `json:"allowedAmount,omitzero"`      // Plan allowed amount for the service (non-EDI)
	PaidAmount         float64  `json:"paidAmount,omitzero"`         // Plan paid amount for the service (non-EDI)
	Quantity           float64  `json:"quantity"`                    // Quantity of the service (from SV104 / SV205)
	Units              string   `json:"units,omitzero"`              // Units connected to the quantity given (from SV103 / SV204)
	PlaceOfService     string   `json:"placeOfService,omitzero"`     // Place of service code (from SV105)
	AmbulancePickupZIP string   `json:"ambulancePickupZIP,omitzero"` // ZIP code where ambulance picked up patient. Supplied if different than claim-level value (from NM1 PW)
}

type RateSheet struct {
	NPI               string             `json:"npi,omitzero"`               // National Provider Identifier of the provider (from NM109, required)
	ProviderFirstName string             `json:"providerFirstName,omitzero"` // First name of the provider (NM104, highly recommended)
	ProviderLastName  string             `json:"providerLastName,omitzero"`  // Last name of the provider (from NM103, highly recommended)
	ProviderOrgName   string             `json:"providerOrgName,omitzero"`   // Organization name of the provider (from NM103, highly recommended)
	ProviderAddress   string             `json:"providerAddress,omitzero"`   // Address of the provider (from N301, highly recommended)
	ProviderCity      string             `json:"providerCity,omitzero"`      // City of the provider (from N401, highly recommended)
	ProviderState     string             `json:"providerState"`              // State of the provider (from N402, highly recommended)
	ProviderZip       string             `json:"providerZIP"`                // ZIP code of the provider (from N403, required)
	FormType          FormType           `json:"formType"`                   // Type of form used to submit the claim. Can be HCFA or UB-04 (from CLM05_02)
	BillTypeOrPOS     string             `json:"billTypeOrPOS"`              // Describes type of facility where services were rendered (from CLM05_01)
	DRG               string             `json:"drg,omitzero"`               // Diagnosis Related Group for inpatient services (from HI DR)
	BilledAmount      float64            `json:"billedAmount,omitzero"`      // Billed amount from provider (from CLM02)
	AllowedAmount     float64            `json:"allowedAmount,omitzero"`     // Amount allowed by the plan for payment. Both member and plan responsibility (non-EDI)
	PaidAmount        float64            `json:"paidAmount,omitzero"`        // Amount paid by the plan for the claim (non-EDI)
	Services          []RateSheetService `json:"services,omitzero"`          // One or more services provided to the patient (from LX loop)
}

type RateSheetService struct {
	ProcedureCode      string   `json:"procedureCode,omitzero"`      // Procedure code (from SV101_02 / SV202_02)
	ProcedureModifiers []string `json:"procedureModifiers,omitzero"` // Procedure modifiers (from SV101_03, 4, 5, 6 / SV202_03, 4, 5, 6)
	BilledAmount       float64  `json:"billedAmount,omitzero"`       // Billed charge for the service (from SV102 / SV203)
	AllowedAmount      float64  `json:"allowedAmount,omitzero"`      // Plan allowed amount for the service (non-EDI)
}
