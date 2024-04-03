package mph

type ClaimRepricingCode string
type LineRepricingCode string
type HospitalType string

const (
	// claim-level repricing codes
	ClaimRepricingCodeMedicare            ClaimRepricingCode = "MED"
	ClaimRepricingCodeContractPricing     ClaimRepricingCode = "CON"
	ClaimRepricingCodeRBPPricing          ClaimRepricingCode = "RBP"
	ClaimRepricingCodeSingleCaseAgreement ClaimRepricingCode = "SCA"
	ClaimRepricingCodeNeedsMoreInfo       ClaimRepricingCode = "IFO"

	// line-level Medicare repricing codes
	LineRepricingCodeMedicare          LineRepricingCode = "MED"
	LineRepricingCodeSyntheticMedicare LineRepricingCode = "SYN"
	LineRepricingCodeCostPercent       LineRepricingCode = "CST"
	LineRepricingCodeMedicarePercent   LineRepricingCode = "MPT"
	LineRepricingCodeMedicareNoOutlier LineRepricingCode = "MNO"
	LineRepricingCodeBilledPercent     LineRepricingCode = "BIL"
	LineRepricingCodeFeeSchedule       LineRepricingCode = "FSC"
	LineRepricingCodePerDiem           LineRepricingCode = "PDM"
	LineRepricingCodeFlatRate          LineRepricingCode = "FLT"
	LineRepricingCodeLimitedToBilled   LineRepricingCode = "LTB"

	// line-level zero dollar repricing explanations
	LineRepricingCodeNotAllowedByMedicare  LineRepricingCode = "NAM"
	LineRepricingCodePackaged              LineRepricingCode = "PKG"
	LineRepricingCodeNeedsMoreInfo         LineRepricingCode = "IFO"
	LineRepricingCodeProcedureCodeProblem  LineRepricingCode = "CPB"
	LineRepricingCodeNotRepricedPerRequest LineRepricingCode = "NRP"

	// Hospital types
	AcuteCareHospitalType      HospitalType = "Acute Care Hospitals"
	CriticalAccessHospitalType HospitalType = "Critical Access Hospitals"
	ChildrensHospitalType      HospitalType = "Childrens"
	PsychiatricHospitalType    HospitalType = "Psychiatric"
	AcuteCareDODHospitalType   HospitalType = "Acute Care - Department of Defense"
)

// Pricing contains the results of a pricing request
type Pricing struct {
	ClaimID                 string                 `json:"claimID,omitempty"`                 // The unique identifier for the claim (copied from input)
	MedicareAmount          float64                `json:"medicareAmount,omitempty"`          // The amount Medicare would pay for the service
	AllowedAmount           float64                `json:"allowedAmount,omitempty"`           // The allowed amount based on a contract or RBP pricing
	AllowedCalculationError string                 `json:"allowedCalculationError,omitempty"` // The reason the allowed amount was not calculated
	MedicareRepricingCode   ClaimRepricingCode     `json:"medicareRepricingCode,omitempty"`   // Explains the methodology used to calculate Medicare (MED or IFO)
	MedicareRepricingNote   string                 `json:"medicareRepricingNote,omitempty"`   // Note explaining approach for pricing or reason for error
	AllowedRepricingCode    ClaimRepricingCode     `json:"allowedRepricingCode,omitempty"`    // Explains the methodology used to calculate allowed amount (CON, RBP, SCA, or IFO)
	AllowedRepricingNote    string                 `json:"allowedRepricingNote,omitempty"`    // Note explaining approach for pricing or reason for error
	MedicareStdDev          float64                `json:"medicareStdDev,omitempty"`          // The standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource          string                 `json:"medicareSource,omitempty"`          // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	InpatientPriceDetail    *InpatientPriceDetail  `json:"inpatientPriceDetail,omitempty"`    // Details about the inpatient pricing
	OutpatientPriceDetail   *OutpatientPriceDetail `json:"outpatientPriceDetail,omitempty"`   // Details about the outpatient pricing
	ProviderDetail          *ProviderDetail        `json:"providerDetail,omitempty"`          // The provider details used when pricing the claim
	EditDetail              *ClaimEdits            `json:"editDetail,omitempty"`              // Errors which cause the claim to be denied, rejected, suspended, or returned to the provider
	PricerResult            string                 `json:"pricerResult,omitempty"`            // Pricer return details
	Services                []Service              `json:"services,omitempty"`                // Pricing for each service line on the claim
	EditError               *ResponseError         `json:"error,omitempty"`                   // An error that occurred during some step of the pricing process
}

// PricedService contains the results of a pricing request for a single service line
type PricedService struct {
	LineNumber                  string            `json:"lineNumber,omitempty"`              // Number of the service line item (copied from input)
	ProviderDetail              *ProviderDetail   `json:"providerDetail,omitempty"`          // Provider Details used when pricing the service if different than the claim
	MedicareAmount              float64           `json:"medicareAmount,omitempty"`          // Amount Medicare would pay for the service
	AllowedAmount               float64           `json:"allowedAmount,omitempty"`           // Allowed amount based on a contract or RBP pricing
	AllowedCalculationError     string            `json:"allowedCalculationError,omitempty"` // Reason the allowed amount was not calculated
	RepricingCode               LineRepricingCode `json:"repricingCode,omitempty"`           // Explains the methodology used to calculate Medicare
	RepricingNote               string            `json:"repricingNote,omitempty"`           // Note explaining approach for pricing or reason for error
	TechnicalComponentAmount    float64           `json:"tcAmount,omitempty"`                // Amount Medicare would pay for the technical component
	ProfessionalComponentAmount float64           `json:"pcAmount,omitempty"`                // Amount Medicare would pay for the professional component
	MedicareStdDev              float64           `json:"medicareStdDev,omitempty"`          // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource              string            `json:"medicareSource,omitempty"`          // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	PricerResult                string            `json:"pricerResult,omitempty"`            // Pricing service return details
	StatusIndicator             string            `json:"statusIndicator,omitempty"`         // Code which gives more detail about how Medicare pays for the service
	PaymentIndicator            string            `json:"paymentIndicator,omitempty"`        // Text which explains the type of payment for Medicare
	PaymentAPC                  string            `json:"paymentAPC,omitempty"`              // Ambulatory Payment Classification
	EditDetail                  *LineEdits        `json:"editDetail,omitempty"`              // Errors which cause the line item to be unable to be priced
}

// InpatientPriceDetail contains pricing details for an inpatient claim
type InpatientPriceDetail struct {
	DRG                            string  `json:"drg,omitempty"`                            // Diagnosis Related Group (DRG) code used to price the claim
	DRGAmount                      float64 `json:"drgAmount,omitempty"`                      // Amount Medicare would pay for the DRG
	PassthroughAmount              float64 `json:"passthroughAmount,omitempty"`              // Per diem amount to cover capital-related costs, direct medical education, and other costs
	OutlierAmount                  float64 `json:"outlierAmount,omitempty"`                  // Additional amount paid for high cost cases
	IndirectMedicalEducationAmount float64 `json:"indirectMedicalEducationAmount,omitempty"` // Additional amount paid for teaching hospitals
	DisproportionateShareAmount    float64 `json:"disproportionateShareAmount,omitempty"`    // Additional amount paid for hospitals with a high number of low-income patients
	UncompensatedCareAmount        float64 `json:"uncompensatedCareAmount,omitempty"`        // Additional amount paid for patients who are unable to pay for their care
	ReadmissionAdjustmentAmount    float64 `json:"readmissionAdjustmentAmount,omitempty"`    // Adjustment amount for hospitals with high readmission rates
	ValueBasedPurchasingAmount     float64 `json:"valueBasedPurchasingAmount,omitempty"`     // Adjustment for hospitals based on quality measures
}

// OutpatientPriceDetail contains pricing details for an outpatient claim
type OutpatientPriceDetail struct {
	OutlierAmount                         float64 `json:"outlierAmount,omitempty"`                         // Additional amount paid for high cost cases
	FirstPassthroughDrugOffsetAmount      float64 `json:"firstPassthroughDrugOffsetAmount,omitempty"`      // Amount built into the APC payment for certain drugs
	SecondPassthroughDrugOffsetAmount     float64 `json:"secondPassthroughDrugOffsetAmount,omitempty"`     // Amount built into the APC payment for certain drugs
	ThirdPassthroughDrugOffsetAmount      float64 `json:"thirdPassthroughDrugOffsetAmount,omitempty"`      // Amount built into the APC payment for certain drugs
	FirstDeviceOffsetAmount               float64 `json:"firstDeviceOffsetAmount,omitempty"`               // Amount built into the APC payment for certain devices
	SecondDeviceOffsetAmount              float64 `json:"secondDeviceOffsetAmount,omitempty"`              // Amount built into the APC payment for certain devices
	FullOrPartialDeviceCreditOffsetAmount float64 `json:"fullOrPartialDeviceCreditOffsetAmount,omitempty"` // Credit for devices that are supplied for free or at a reduced cost
	TerminatedDeviceProcedureOffsetAmount float64 `json:"terminatedDeviceProcedureOffsetAmount,omitempty"` // Credit for devices that are not used due to a terminated procedure
}

// ProviderDetail contains basic information about the provider and/or locality used for pricing.
// Not all fields are returned with every pricing request. For example, the CMS Certification
// Number (CCN) is only returned for facilities which have a CCN such as hospitals.
type ProviderDetail struct {
	CCN            string       `json:"ccn,omitempty"`           // CMS Certification Number for the facility
	MAC            uint16       `json:"mac"`                     // Medicare Administrative Contractor number
	Locality       uint8        `json:"locality"`                // Geographic locality number used for pricing
	RuralIndicator byte         `json:"ruralIndicator"`          // Indicates whether provider is Rural (R), Super Rural (B), or Urban (blank)
	SpecialtyType  string       `json:"specialtyType,omitempty"` // Medicare provider specialty type
	HospitalType   HospitalType `json:"hospitalType,omitempty"`  // Type of hospital
}

// ClaimEdits contains errors which cause the claim to be denied, rejected, suspended, or returned to the provider.
type ClaimEdits struct {
	ClaimOverallDisposition          string   `json:"claimOverallDisposition,omitempty"`
	ClaimRejectionDisposition        string   `json:"claimRejectionDisposition,omitempty"`
	ClaimDenialDisposition           string   `json:"claimDenialDisposition,omitempty"`
	ClaimReturnToProviderDisposition string   `json:"claimReturnToProviderDisposition,omitempty"`
	ClaimSuspensionDisposition       string   `json:"claimSuspensionDisposition,omitempty"`
	LineItemRejectionDisposition     string   `json:"lineItemRejectionDisposition,omitempty"`
	LineItemDenialDisposition        string   `json:"lineItemDenialDisposition,omitempty"`
	ClaimRejectionReasons            []string `json:"claimRejectionReasons,omitempty"`
	ClaimDenialReasons               []string `json:"claimDenialReasons,omitempty"`
	ClaimReturnToProviderReasons     []string `json:"claimReturnToProviderReasons,omitempty"`
	ClaimSuspensionReasons           []string `json:"claimSuspensionReasons,omitempty"`
	LineItemRejectionReasons         []string `json:"lineItemRejectionReasons,omitempty"`
	LineItemDenialReasons            []string `json:"lineItemDenialReasons,omitempty"`
}

// LineEdits contains errors which cause the line item to be unable to be priced.
type LineEdits struct {
	DenialOrRejectionText string   `json:"denialOrRejectionText,omitempty"`
	ProcedureEdits        []string `json:"procedureEdits,omitempty"`
	Modifier1Edits        []string `json:"modifier1Edits,omitempty"`
	Modifier2Edits        []string `json:"modifier2Edits,omitempty"`
	Modifier3Edits        []string `json:"modifier3Edits,omitempty"`
	Modifier4Edits        []string `json:"modifier4Edits,omitempty"`
	Modifier5Edits        []string `json:"modifier5Edits,omitempty"`
	DataEdits             []string `json:"dataEdits,omitempty"`
	RevenueEdits          []string `json:"revenueEdits,omitempty"`
	ProfessionalEdits     []string `json:"professionalEdits,omitempty"`
}
