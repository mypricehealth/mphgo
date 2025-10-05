package mph

import (
	"encoding/json"
	"fmt"
	"strings"

	"braces.dev/errtrace"
	"github.com/mypricehealth/mphgo/set"
)

const (
	editErrorTitle        = "claim edits failed"
	fatalEditErrorTitle   = "fatal edit error"
	editErrorDetail       = "see editDetail for more information"
	PriceErrorTitle       = "pricing not available"
	SyntheticPricerResult = "Processed via synthetic Medicare"
)

var (
	ErrorEditSeeDetail = &ResponseError{Title: editErrorTitle, Detail: editErrorDetail}
	ErrorEditFatal     = &ResponseError{Title: fatalEditErrorTitle, Detail: "claim must be returned to provider for resolution"}
)

type ClaimRepricingCode string
type LineRepricingCode string
type HospitalType string
type RuralIndicator string
type MedicareSource string
type Step string
type Status string

var _ json.Unmarshaler = new(RuralIndicator)

func (r *RuralIndicator) UnmarshalJSON(data []byte) error { // old code was sending the int value of the Rural indicator, so this handles both
	var strData string
	err := json.Unmarshal(data, &strData)
	if err == nil && (strData == "R" || strData == "B" || strData == "") {
		*r = RuralIndicator(strData)
		return nil
	}

	var intData int
	err = json.Unmarshal(data, &intData)
	if err == nil {
		switch intData {
		case 66:
			*r = RuralIndicator("B")
			return nil
		case 82:
			*r = RuralIndicator("R")
			return nil
		case 0:
			*r = RuralIndicator("")
			return nil
		}
	}
	return errtrace.Errorf("invalid RuralIndicator value: %s", strData)
}

const (
	// claim-level repricing codes.

	ClaimRepricingCodeMedicare            ClaimRepricingCode = "MED"
	ClaimRepricingCodeContractPricing     ClaimRepricingCode = "CON"
	ClaimRepricingCodeRBPPricing          ClaimRepricingCode = "RBP"
	ClaimRepricingCodeCoralRBPPricing     ClaimRepricingCode = "CRBP"
	ClaimRepricingCodeSingleCaseAgreement ClaimRepricingCode = "SCA"
	ClaimRepricingCodeNeedsMoreInfo       ClaimRepricingCode = "IFO"
	ClaimRepricingCodeOutOfNetwork        ClaimRepricingCode = "OON"

	// line-level Medicare repricing codes.

	LineRepricingCodeMedicare          LineRepricingCode = "MED"
	LineRepricingCodeMedicarePercent   LineRepricingCode = "MPT"
	LineRepricingCodeMedicareNoOutlier LineRepricingCode = "MNO"
	LineRepricingCodeSyntheticMedicare LineRepricingCode = "SYN"
	LineRepricingCodeBilledPercent     LineRepricingCode = "BIL"
	LineRepricingCodeFeeSchedule       LineRepricingCode = "FSC"
	LineRepricingCodePerDiem           LineRepricingCode = "PDM"
	LineRepricingCodeFlatRate          LineRepricingCode = "FLT"
	LineRepricingCodeCostPercent       LineRepricingCode = "CST"
	LineRepricingCodeLimitedToBilled   LineRepricingCode = "LTB"

	// line-level zero dollar repricing explanations.

	LineRepricingCodeNotRepricedPerRequest LineRepricingCode = "NRP"
	LineRepricingCodeNotAllowedByMedicare  LineRepricingCode = "NAM"
	LineRepricingCodePackaged              LineRepricingCode = "PKG"
	LineRepricingCodeNeedsMoreInfo         LineRepricingCode = "IFO"
	LineRepricingCodeProcedureCodeProblem  LineRepricingCode = "CPB"
	LineRepricingCodeOutOfNetwork          LineRepricingCode = "OON"

	// Hospital types.

	AcuteCareHospitalType      HospitalType = "Acute Care Hospitals"
	CriticalAccessHospitalType HospitalType = "Critical Access Hospitals"
	ChildrensHospitalType      HospitalType = "Childrens"
	PsychiatricHospitalType    HospitalType = "Psychiatric"
	AcuteCareDODHospitalType   HospitalType = "Acute Care - Department of Defense"
)

var (
	ClaimRepricingCodes = map[ClaimRepricingCode]struct{}{
		ClaimRepricingCodeMedicare:            {},
		ClaimRepricingCodeContractPricing:     {},
		ClaimRepricingCodeRBPPricing:          {},
		ClaimRepricingCodeCoralRBPPricing:     {},
		ClaimRepricingCodeSingleCaseAgreement: {},
		ClaimRepricingCodeNeedsMoreInfo:       {},
		ClaimRepricingCodeOutOfNetwork:        {},
	}
	LineRepricingCodes = map[LineRepricingCode]struct{}{
		LineRepricingCodeMedicare:              {},
		LineRepricingCodeMedicarePercent:       {},
		LineRepricingCodeMedicareNoOutlier:     {},
		LineRepricingCodeSyntheticMedicare:     {},
		LineRepricingCodeBilledPercent:         {},
		LineRepricingCodeFeeSchedule:           {},
		LineRepricingCodePerDiem:               {},
		LineRepricingCodeFlatRate:              {},
		LineRepricingCodeCostPercent:           {},
		LineRepricingCodeLimitedToBilled:       {},
		LineRepricingCodeNotRepricedPerRequest: {},
		LineRepricingCodeNotAllowedByMedicare:  {},
		LineRepricingCodePackaged:              {},
		LineRepricingCodeNeedsMoreInfo:         {},
		LineRepricingCodeProcedureCodeProblem:  {},
		LineRepricingCodeOutOfNetwork:          {},
	}
)

const (
	// Rural indicators.
	RuralIndicatorRural      RuralIndicator = "R"
	RuralIndicatorSuperRural RuralIndicator = "B"
	RuralIndicatorUrban      RuralIndicator = ""
)

const (
	MedicareSourceAmbulance              MedicareSource = "AmbulanceFS"
	MedicareSourceAnesthesia             MedicareSource = "AnesthesiaFS"
	MedicareSourceASC                    MedicareSource = "ASC pricer"
	MedicareSourceCriticalAccessHospital MedicareSource = "CAH pricer"
	MedicareSourceDME                    MedicareSource = "DMEFS"
	MedicareSourceDrugs                  MedicareSource = "DrugsFS"
	MedicareSourceEditError              MedicareSource = "Claim editor"
	MedicareSourceEstimateByCodeOnly     MedicareSource = "CodeOnly"
	MedicareSourceEstimateByLocalityCode MedicareSource = "LocalityCode"
	MedicareSourceEstimateByLocalityOnly MedicareSource = "LocalityOnly"
	MedicareSourceEstimateByNational     MedicareSource = "National"
	MedicareSourceEstimateByStateCode    MedicareSource = "StateCode"
	MedicareSourceEstimateByStateOnly    MedicareSource = "StateOnly"
	MedicareSourceEstimateByUnknown      MedicareSource = "Unknown"
	MedicareSourceInpatient              MedicareSource = "IPPS"
	MedicareSourceLabs                   MedicareSource = "LabsFS"
	MedicareSourceMPFS                   MedicareSource = "MPFS"
	MedicareSourceOutpatient             MedicareSource = "Outpatient pricer"
	MedicareSourceManualPricing          MedicareSource = "Manual Pricing"
	MedicareSourceSNF                    MedicareSource = "SNF PPS"
	MedicareSourceSynthetic              MedicareSource = "Synthetic Medicare"
)

const (
	stepNew                  Step = "New"
	stepReceived             Step = "Received"
	StepPending              Step = "Pending"
	stepHeld                 Step = "Held"
	stepError                Step = "Error"
	stepInputValidated       Step = "Input Validated"
	stepProviderMatched      Step = "Provider Matched"
	stepEditComplete         Step = "Edit Complete"
	stepMedicarePriced       Step = "Medicare Priced"
	stepPrimaryAllowedPriced Step = "Primary Allowed Priced"
	stepNetworkAllowedPriced Step = "Network Allowed Priced"
	stepOutOfNetwork         Step = "Out of Network"
	stepRequestMoreInfo      Step = "Request More Info"
	stepPriced               Step = "Priced"
	stepReturned             Step = "Returned"
)

const (
	statusPendingClaimInputValidation        Status = "Claim Input Validation"
	statusPendingClaimEditReview             Status = "Claim Edit Review"
	statusPendingProviderMatching            Status = "Provider Matching"
	statusPendingMedicareReview              Status = "Medicare Review"
	statusPendingMedicareCalculation         Status = "Medicare Calculation"
	statusPendingPrimaryAllowedReview        Status = "Primary Allowed Review"
	statusPendingNetworkAllowedReview        Status = "Network Allowed Review"
	statusPendingPrimaryAllowedDetermination Status = "Primary Allowed Determination"
	statusPendingNetworkAllowedDetermination Status = "Network Allowed Determination"
)

var (
	StatusNew                                = ClaimStatus{Step: stepNew}                                                       // created by TPA. We use the transaction date as a proxy for this date
	StatusReceived                           = ClaimStatus{Step: stepReceived}                                                  // received and ready for processing. This is modified date of the file we get from SFTP
	StatusHeld                               = ClaimStatus{Step: stepHeld}                                                      // held for various reasons
	StatusError                              = ClaimStatus{Step: stepError}                                                     // claim encountered an error during processing
	StatusInputValidated                     = ClaimStatus{Step: stepInputValidated}                                            // claim input has been validated
	StatusProviderMatched                    = ClaimStatus{Step: stepProviderMatched}                                           // providers in the claim have been matched to the provider system of record
	StatusEditComplete                       = ClaimStatus{Step: stepEditComplete}                                              // claim has been edited and is ready for pricing
	StatusMedicarePriced                     = ClaimStatus{Step: stepMedicarePriced}                                            // claim has been priced according to Medicare
	StatusPrimaryAllowedPriced               = ClaimStatus{Step: stepPrimaryAllowedPriced}                                      // claim has been priced according to the primary allowed amount (e.g. contract, RBP, etc.)
	StatusNetworkAllowedPriced               = ClaimStatus{Step: stepNetworkAllowedPriced}                                      // claim has been priced according to the allowed amount of the network
	StatusOutOfNetwork                       = ClaimStatus{Step: stepOutOfNetwork}                                              // is out of network
	StatusRequestMoreInfo                    = ClaimStatus{Step: stepRequestMoreInfo}                                           // return claim to trading partner for more information to enable correct processing
	StatusPriced                             = ClaimStatus{Step: stepPriced}                                                    // done pricing
	StatusReturned                           = ClaimStatus{Step: stepReturned}                                                  // returned to TPA
	StatusPendingClaimInputValidation        = ClaimStatus{Step: StepPending, Status: statusPendingClaimInputValidation}        // waiting for claim input validation
	StatusPendingClaimEditReview             = ClaimStatus{Step: StepPending, Status: statusPendingClaimEditReview}             // waiting for claim edit review
	StatusPendingProviderMatching            = ClaimStatus{Step: StepPending, Status: statusPendingProviderMatching}            // waiting for provider matching
	StatusPendingMedicareReview              = ClaimStatus{Step: StepPending, Status: statusPendingMedicareReview}              // waiting for Medicare amount review
	StatusPendingMedicareCalculation         = ClaimStatus{Step: StepPending, Status: statusPendingMedicareCalculation}         // waiting for Medicare amount calculation
	StatusPendingPrimaryAllowedReview        = ClaimStatus{Step: StepPending, Status: statusPendingPrimaryAllowedReview}        // waiting for primary allowed amount review
	StatusPendingNetworkAllowedReview        = ClaimStatus{Step: StepPending, Status: statusPendingNetworkAllowedReview}        // waiting for network allowed amount review
	StatusPendingPrimaryAllowedDetermination = ClaimStatus{Step: StepPending, Status: statusPendingPrimaryAllowedDetermination} // waiting for the primary allowed amount (e.g. contract, RBP rate, etc.) to be determined
	StatusPendingNetworkAllowedDetermination = ClaimStatus{Step: StepPending, Status: statusPendingNetworkAllowedDetermination} // waiting for allowed amount from the network
)

type ClaimStatus struct {
	Step   Step   `json:"step,omitempty"`
	Status Status `json:"status,omitempty"`
}

func (s ClaimStatus) IsEmpty() bool {
	return s.Step == "" && s.Status == ""
}

func (s ClaimStatus) String() string {
	if s.Status != "" {
		return fmt.Sprintf("%s %s", s.Step, s.Status)
	}
	return string(s.Step)
}

// Pricing contains the results of a pricing request
type Pricing struct {
	ClaimID               string                `json:"claimID,omitzero"               db:"-"`                       // The unique identifier for the claim (copied from input)
	MedicareAmount        float64               `json:"medicareAmount,omitzero"        db:"medicare_amount"`         // The amount Medicare would pay for the service
	AllowedAmount         float64               `json:"allowedAmount,omitzero"         db:"allowed_amount"`          // The allowed amount based on a contract or RBP pricing
	MedicareRepricingCode ClaimRepricingCode    `json:"medicareRepricingCode,omitzero" db:"medicare_repricing_code"` // Explains the methodology used to calculate Medicare (MED or IFO)
	MedicareRepricingNote string                `json:"medicareRepricingNote,omitzero" db:"medicare_repricing_note"` // Note explaining approach for pricing or reason for error
	NetworkCode           string                `json:"networkCode,omitzero"           db:"network_code"`            // Code describing the network used for allowed amount pricing
	AllowedRepricingCode  ClaimRepricingCode    `json:"allowedRepricingCode,omitzero"  db:"allowed_repricing_code"`  // Explains the methodology used to calculate allowed amount (CON, RBP, SCA, or IFO)
	AllowedRepricingNote  string                `json:"allowedRepricingNote,omitzero"  db:"allowed_repricing_note"`  // Note explaining approach for pricing or reason for error
	MedicareStdDev        float64               `json:"medicareStdDev,omitzero"        db:"medicare_std_dev"`        // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource        MedicareSource        `json:"medicareSource,omitzero"        db:"medicare_source"`         // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	InpatientPriceDetail  InpatientPriceDetail  `json:"inpatientPriceDetail,omitzero"  db:",inline"`                 // Details about the inpatient pricing
	OutpatientPriceDetail OutpatientPriceDetail `json:"outpatientPriceDetail,omitzero" db:",inline"`                 // Details about the outpatient pricing
	ProviderDetail        ProviderDetail        `json:"providerDetail,omitzero"        db:",inline"`                 // The provider details used when pricing the claim
	EditDetail            *ClaimEdits           `json:"editDetail,omitzero"            db:",inline"`                 // Errors which cause the claim to be denied, rejected, suspended, or returned to the provider
	PricerResult          string                `json:"pricerResult,omitzero"          db:"pricer_result"`           // Pricer return details
	PriceConfig           PriceConfig           `json:"priceConfig,omitzero"           db:",inline"`                 // The configuration used for pricing the claim
	Services              []PricedService       `json:"services,omitzero,omitempty"    db:"services"`                // Pricing for each service line on the claim
	EditError             *ResponseError        `json:"editError,omitzero"             db:"edit_error"`              // An error that occurred during some step of the pricing process
}

func (p Pricing) IsEmpty() bool {
	return p.ClaimID == "" && p.MedicareAmount == 0 && p.AllowedAmount == 0 && p.MedicareRepricingCode == "" && p.MedicareRepricingNote == "" && p.NetworkCode == "" &&
		p.AllowedRepricingCode == "" && p.AllowedRepricingNote == "" && p.MedicareStdDev == 0 && p.MedicareSource == "" && p.InpatientPriceDetail.IsEmpty() &&
		p.OutpatientPriceDetail.IsEmpty() && p.ProviderDetail.IsEmpty() && p.EditDetail.IsEmpty() &&
		p.PricerResult == "" && len(p.Services) == 0 && p.EditError == nil
}

func (p Pricing) GetRepricingNote() string {
	var buf strings.Builder
	if p.EditError.HasSpecificMessage() {
		buf.WriteString(p.EditError.Detail)
	}
	addSeparatedMessage(&buf, p.EditDetail.GetMessage())
	if p.AllowedRepricingNote != "" {
		addSeparatedMessage(&buf, p.AllowedRepricingNote)
	} else if p.MedicareRepricingNote != "" {
		addSeparatedMessage(&buf, p.MedicareRepricingNote)
	}
	return buf.String()
}

func addSeparatedMessage(buf *strings.Builder, message string) {
	if buf.Len() > 0 {
		buf.WriteString(". ")
	}
	buf.WriteString(message)
}

func (p Pricing) GetEditMessages() []string {
	if p.EditDetail == nil {
		return nil
	}
	messages := set.NewOrderedSet[string]()
	e := p.EditDetail
	messages.AddSlices(e.ClaimRejectionReasons, e.ClaimDenialReasons, e.ClaimReturnToProviderReasons, e.ClaimSuspensionReasons, e.LineItemRejectionReasons, e.LineItemDenialReasons)
	return messages.Items()
}

func (p Pricing) HasFatalError() bool {
	return p.EditError != nil && p.EditError.Title == fatalEditErrorTitle
}

// PricedService contains the results of a pricing request for a single service line.
type PricedService struct {
	LineNumber                    string                  `json:"lineNumber,omitzero"                    db:"-"`                                  // Number of the service line item (copied from input)
	ProviderDetail                *ProviderDetail         `json:"providerDetail,omitzero"                db:",inline"`                            // Provider Details used when pricing the service if different than the claim
	MedicareAmount                float64                 `json:"medicareAmount,omitzero"                db:"medicare_amount"`                    // Amount Medicare would pay for the service
	AllowedAmount                 float64                 `json:"allowedAmount,omitzero"                 db:"allowed_amount"`                     // Allowed amount based on a contract or RBP pricing
	MedicareRepricingCode         LineRepricingCode       `json:"medicareRepricingCode,omitzero"         db:"medicare_repricing_code"`            // Explains the methodology used to calculate Medicare
	MedicareRepricingNote         string                  `json:"medicareRepricingNote,omitzero"         db:"medicare_repricing_note"`            // Note explaining approach for pricing or reason for error
	NetworkCode                   string                  `json:"networkCode,omitzero"                   db:"network_code"`                       // Code describing the network used for allowed amount pricing
	AllowedRepricingCode          LineRepricingCode       `json:"allowedRepricingCode,omitzero"          db:"allowed_repricing_code"`             // Explains the methodology used to calculate allowed amount
	AllowedRepricingNote          string                  `json:"allowedRepricingNote,omitzero"          db:"allowed_repricing_note"`             // Note explaining approach for pricing or reason for error
	AllowedRepricingFormula       AllowedRepricingFormula `json:"allowedRepricingFormula,omitzero"       db:",inline"`                            // Formula used to calculate the allowed amount
	TechnicalComponentAmount      float64                 `json:"tcAmount,omitzero"                      db:"technical_component_amount"`         // Amount Medicare would pay for the technical component
	ProfessionalComponentAmount   float64                 `json:"pcAmount,omitzero"                      db:"professional_component_amount"`      // Amount Medicare would pay for the professional component
	MedicareStdDev                float64                 `json:"medicareStdDev,omitzero"                db:"medicare_std_dev"`                   // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource                MedicareSource          `json:"medicareSource,omitzero"                db:"medicare_source"`                    // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	PricerResult                  string                  `json:"pricerResult,omitzero"                  db:"pricer_result"`                      // Pricing service return details
	StatusIndicator               string                  `json:"statusIndicator,omitzero"               db:"status_indicator"`                   // Code which gives more detail about how Medicare pays for the service (outpatient + professional)
	PaymentIndicator              string                  `json:"paymentIndicator,omitzero"              db:"payment_indicator"`                  // Text which explains the type of payment for Medicare (outpatient only)
	DiscountFormula               string                  `json:"discountFormula,omitzero"               db:"discount_formula"`                   // The multi-procedure discount formula used to calculate the allowed amount (outpatient only)
	LineItemDenialOrRejectionFlag string                  `json:"lineItemDenialOrRejectionFlag,omitzero" db:"line_item_denial_or_rejection_flag"` // Identifies how a line item was denied or rejected and how the rejection can be overridden (outpatient only)
	PackagingFlag                 string                  `json:"packagingFlag,omitzero"                 db:"packaging_flag"`                     // Indicates if the service is packaged and the reason for packaging (outpatient only)
	PaymentAdjustmentFlag         string                  `json:"paymentAdjustmentFlag,omitzero"         db:"payment_adjustment_flag"`            // Identifies special adjustments made to the payment (outpatient only)
	PaymentAdjustmentFlag2        string                  `json:"paymentAdjustmentFlag2,omitzero"        db:"payment_adjustment_flag2"`           // Identifies special adjustments made to the payment (outpatient only)
	PaymentMethodFlag             string                  `json:"paymentMethodFlag,omitzero"             db:"payment_method_flag"`                // The method used to calculate the allowed amount (outpatient only)
	CompositeAdjustmentFlag       string                  `json:"compositeAdjustmentFlag,omitzero"       db:"composite_adjustment_flag"`          // Assists in composite APC determination (outpatient only)
	HCPCSAPC                      string                  `json:"hcpcsAPC,omitzero"                      db:"hcpcs_apc"`                          // Ambulatory Payment Classification code of the line item HCPCS (outpatient only)
	PaymentAPC                    string                  `json:"paymentAPC,omitzero"                    db:"payment_apc"`                        // Ambulatory Payment Classification code used for payment (outpatient only)
	EditDetail                    *LineEdits              `json:"editDetail,omitzero"                    db:",inline"`                            // Errors which cause the line item to be unable to be priced
}

type AllowedRepricingFormula struct {
	MedicarePercent float64 `json:"medicarePercent,omitzero" db:"allowed_repricing_formula_medicare_percent"` // Percentage of the Medicare amount used to calculate the allowed amount
	BilledPercent   float64 `json:"billedPercent,omitzero"   db:"allowed_repricing_formula_billed_percent"`   // Percentage of the billed amount used to calculate the allowed amount
	FeeSchedule     float64 `json:"feeSchedule,omitzero"     db:"allowed_repricing_formula_fee_schedule"`     // Fee schedule amount used as the allowed amount
	FixedAmount     float64 `json:"fixedAmount,omitzero"     db:"allowed_repricing_formula_fixed_amount"`     // Fixed amount used as the allowed amount
	PerDiem         float64 `json:"perDiem,omitzero"         db:"allowed_repricing_formula_per_diem"`         // Per diem rate used to calculate the allowed amount
}

func (a AllowedRepricingFormula) IsEmpty() bool {
	return a == AllowedRepricingFormula{}
}

func (s PricedService) GetRepricingNote() string {
	var buf strings.Builder
	buf.WriteString(s.EditDetail.GetMessage())
	if s.AllowedRepricingNote != "" {
		addSeparatedMessage(&buf, s.AllowedRepricingNote)
	} else if s.MedicareRepricingNote != "" {
		addSeparatedMessage(&buf, s.MedicareRepricingNote)
	}
	return buf.String()
}

// InpatientPriceDetail contains pricing details for an inpatient claim.
type InpatientPriceDetail struct {
	DRG                            string  `json:"drg,omitzero"                            db:"inpatient_drg"`                               // Diagnosis Related Group (DRG) code used to price the claim
	DRGAmount                      float64 `json:"drgAmount,omitzero"                      db:"inpatient_drg_amount"`                        // Amount Medicare would pay for the DRG
	PassthroughAmount              float64 `json:"passthroughAmount,omitzero"              db:"inpatient_passthrough_amount"`                // Per diem amount to cover capital-related costs, direct medical education, and other costs
	OutlierAmount                  float64 `json:"outlierAmount,omitzero"                  db:"inpatient_outlier_amount"`                    // Additional amount paid for high cost cases
	IndirectMedicalEducationAmount float64 `json:"indirectMedicalEducationAmount,omitzero" db:"inpatient_indirect_medical_education_amount"` // Additional amount paid for teaching hospitals
	DisproportionateShareAmount    float64 `json:"disproportionateShareAmount,omitzero"    db:"inpatient_disproportionate_share_amount"`     // Additional amount paid for hospitals with a high number of low-income patients
	UncompensatedCareAmount        float64 `json:"uncompensatedCareAmount,omitzero"        db:"inpatient_uncompensated_care_amount"`         // Additional amount paid for patients who are unable to pay for their care
	ReadmissionAdjustmentAmount    float64 `json:"readmissionAdjustmentAmount,omitzero"    db:"inpatient_readmission_adjustment_amount"`     // Adjustment amount for hospitals with high readmission rates
	ValueBasedPurchasingAmount     float64 `json:"valueBasedPurchasingAmount,omitzero"     db:"inpatient_value_based_purchasing_amount"`     // Adjustment for hospitals based on quality measures
	WageIndex                      float64 `json:"wageIndex,omitzero"                      db:"inpatient_wage_index"`                        // Wage index used for geographic adjustment
}

func (i InpatientPriceDetail) IsEmpty() bool {
	return i == InpatientPriceDetail{}
}

// OutpatientPriceDetail contains pricing details for an outpatient claim.
type OutpatientPriceDetail struct {
	OutlierAmount                         float64 `json:"outlierAmount,omitzero"                         db:"outpatient_outlier_amount"`                              // Additional amount paid for high cost cases
	FirstPassthroughDrugOffsetAmount      float64 `json:"firstPassthroughDrugOffsetAmount,omitzero"      db:"outpatient_first_passthrough_drug_offset_amount"`        // Amount built into the APC payment for certain drugs
	SecondPassthroughDrugOffsetAmount     float64 `json:"secondPassthroughDrugOffsetAmount,omitzero"     db:"outpatient_second_passthrough_drug_offset_amount"`       // Amount built into the APC payment for certain drugs
	ThirdPassthroughDrugOffsetAmount      float64 `json:"thirdPassthroughDrugOffsetAmount,omitzero"      db:"outpatient_third_passthrough_drug_offset_amount"`        // Amount built into the APC payment for certain drugs
	FirstDeviceOffsetAmount               float64 `json:"firstDeviceOffsetAmount,omitzero"               db:"outpatient_first_device_offset_amount"`                  // Amount built into the APC payment for certain devices
	SecondDeviceOffsetAmount              float64 `json:"secondDeviceOffsetAmount,omitzero"              db:"outpatient_second_device_offset_amount"`                 // Amount built into the APC payment for certain devices
	FullOrPartialDeviceCreditOffsetAmount float64 `json:"fullOrPartialDeviceCreditOffsetAmount,omitzero" db:"outpatient_full_or_partial_device_credit_offset_amount"` // Credit for devices that are supplied for free or at a reduced cost
	TerminatedDeviceProcedureOffsetAmount float64 `json:"terminatedDeviceProcedureOffsetAmount,omitzero" db:"outpatient_terminated_device_procedure_offset_amount"`   // Credit for devices that are not used due to a terminated procedure
	WageIndex                             float64 `json:"wageIndex,omitzero"                             db:"outpatient_wage_index"`                                  // Wage index used for geographic adjustment
}

func (o OutpatientPriceDetail) IsEmpty() bool {
	return o == OutpatientPriceDetail{}
}

// ProviderDetail contains basic information about the provider and/or locality used for pricing
// Not all fields are returned with every pricing request. For example, the CMS Certification
// Number (CCN) is only returned for facilities which have a CCN such as hospitals.
type ProviderDetail struct {
	CCN                   string         `json:"ccn,omitzero"            db:"provider_ccn"`             // CMS Certification Number for the facility
	MAC                   uint16         `json:"mac"                     db:"provider_mac"`             // Medicare Administrative Contractor number
	Locality              uint8          `json:"locality"                db:"provider_locality"`        // Geographic locality number used for pricing
	GeographicCBSA        uint32         `json:"geographicCBSA,omitzero" db:"provider_geographic_cbsa"` // Core-Based Statistical Area (CBSA) number for provider ZIP
	StateCBSA             uint8          `json:"stateCBSA,omitzero"      db:"provider_state_cbsa"`      // State Core-Based Statistical Area (CBSA) number
	RuralIndicator        RuralIndicator `json:"ruralIndicator,omitzero" db:"provider_rural_indicator"` // Indicates whether provider is Rural (R), Super Rural (B), or Urban (blank)
	SpecialtyType         string         `json:"specialtyType,omitzero"  db:"provider_specialty_type"`  // Medicare provider specialty type
	HospitalType          HospitalType   `json:"hospitalType,omitzero"   db:"provider_hospital_type"`   // Type of hospital
	BilledToMedicareRatio float64        `json:"-"                       db:"-"`                        // used for synthetic Medicare. Internal use only
}

var empty ProviderDetail

func (p *ProviderDetail) IsEmpty() bool {
	return p == nil || *p == empty
}

// ClaimEdits contains errors which cause the claim to be denied, rejected, suspended, or returned to the provider.
type ClaimEdits struct {
	HCP13DenyCode                    string   `json:"hcpDenyCode,omitzero"                      db:"hcp_deny_code"`                             // The deny code that will be placed into the HCP13 data element for EDI 837 claims
	ClaimOverallDisposition          string   `json:"claimOverallDisposition,omitzero"          db:"claim_edit_overall_disposition"`            // Overall explanation of why the claim edit failed
	ClaimRejectionDisposition        string   `json:"claimRejectionDisposition,omitzero"        db:"claim_edit_rejection_disposition"`          // Explanation of why the claim was rejected
	ClaimDenialDisposition           string   `json:"claimDenialDisposition,omitzero"           db:"claim_edit_denial_disposition"`             // Explanation of why the claim was denied
	ClaimReturnToProviderDisposition string   `json:"claimReturnToProviderDisposition,omitzero" db:"claim_edit_return_to_provider_disposition"` // Explanation of why the claim should be returned to provider
	ClaimSuspensionDisposition       string   `json:"claimSuspensionDisposition,omitzero"       db:"claim_edit_suspension_disposition"`         // Explanation of why the claim was suspended
	LineItemRejectionDisposition     string   `json:"lineItemRejectionDisposition,omitzero"     db:"line_item_edit_rejection_disposition"`      // Explanation of why the line item was rejected
	LineItemDenialDisposition        string   `json:"lineItemDenialDisposition,omitzero"        db:"line_item_edit_denial_disposition"`         // Explanation of why the line item was denied
	ClaimRejectionReasons            []string `json:"claimRejectionReasons,omitzero"            db:"claim_edit_rejection_reasons"`              // Detailed reason(s) describing why the claim was rejected
	ClaimDenialReasons               []string `json:"claimDenialReasons,omitzero"               db:"claim_edit_denial_reasons"`                 // Detailed reason(s) describing why the claim was denied
	ClaimReturnToProviderReasons     []string `json:"claimReturnToProviderReasons,omitzero"     db:"claim_edit_return_to_provider_reasons"`     // Detailed reason(s) describing why the claim should be returned to provider
	ClaimSuspensionReasons           []string `json:"claimSuspensionReasons,omitzero"           db:"claim_edit_suspension_reasons"`             // Detailed reason(s) describing why the claim was suspended
	LineItemRejectionReasons         []string `json:"lineItemRejectionReasons,omitzero"         db:"line_item_edit_rejection_reasons"`          // Detailed reason(s) describing why the line item was rejected
	LineItemDenialReasons            []string `json:"lineItemDenialReasons,omitzero"            db:"line_item_edit_denial_reasons"`             // Detailed reason(s) describing why the line item was denied
}

func (e *ClaimEdits) IsEmpty() bool {
	return e == nil || e.HCP13DenyCode == "" && e.ClaimOverallDisposition == "" && e.ClaimRejectionDisposition == "" && e.ClaimDenialDisposition == "" &&
		e.ClaimReturnToProviderDisposition == "" && e.ClaimSuspensionDisposition == "" && e.LineItemRejectionDisposition == "" && e.LineItemDenialDisposition == "" &&
		len(e.ClaimRejectionReasons) == 0 && len(e.ClaimDenialReasons) == 0 && len(e.ClaimReturnToProviderReasons) == 0 && len(e.ClaimSuspensionReasons) == 0 &&
		len(e.LineItemRejectionReasons) == 0 && len(e.LineItemDenialReasons) == 0
}

func (e *ClaimEdits) GetMessage() string {
	if e == nil {
		return ""
	}
	edits := set.NewOrderedSet[string]()
	edits.AddSlices(e.ClaimRejectionReasons, e.ClaimDenialReasons, e.ClaimReturnToProviderReasons, e.ClaimSuspensionReasons, e.LineItemRejectionReasons, e.LineItemDenialReasons)
	return strings.Join(edits.Items(), "|")
}

// LineEdits contains errors which cause the line item to be unable to be priced.
type LineEdits struct {
	ProcedureEdits []string `json:"procedureEdits,omitzero" db:"procedure_edits"` // Detailed description of each procedure code edit error (from outpatient editor)
	Modifier1Edits []string `json:"modifier1Edits,omitzero" db:"modifier1_edits"` // Detailed description of each edit error for the first procedure code modifier (from outpatient editor)
	Modifier2Edits []string `json:"modifier2Edits,omitzero" db:"modifier2_edits"` // Detailed description of each edit error for the second procedure code modifier (from outpatient editor)
	Modifier3Edits []string `json:"modifier3Edits,omitzero" db:"modifier3_edits"` // Detailed description of each edit error for the third procedure code modifier (from outpatient editor)
	Modifier4Edits []string `json:"modifier4Edits,omitzero" db:"modifier4_edits"` // Detailed description of each edit error for the fourth procedure code modifier (from outpatient editor)
	Modifier5Edits []string `json:"modifier5Edits,omitzero" db:"modifier5_edits"` // Detailed description of each edit error for the fifth procedure code modifier (from outpatient editor)
	DataEdits      []string `json:"dataEdits,omitzero"      db:"data_edits"`      // Detailed description of each data edit error (from outpatient editor)
	RevenueEdits   []string `json:"revenueEdits,omitzero"   db:"revenue_edits"`   // Detailed description of each revenue code edit error (from outpatient editor)
}

func (e *LineEdits) IsEmpty() bool {
	return e == nil || len(e.ProcedureEdits) == 0 && len(e.Modifier1Edits) == 0 && len(e.Modifier2Edits) == 0 &&
		len(e.Modifier3Edits) == 0 && len(e.Modifier4Edits) == 0 && len(e.Modifier5Edits) == 0 &&
		len(e.DataEdits) == 0 && len(e.RevenueEdits) == 0
}

func (e *LineEdits) GetMessage() string {
	if e == nil {
		return ""
	}
	edits := append(append(append(append(append(append(append(e.ProcedureEdits, e.RevenueEdits...), e.Modifier1Edits...), e.Modifier2Edits...), e.Modifier3Edits...), e.Modifier4Edits...), e.Modifier5Edits...), e.DataEdits...)
	return strings.Join(edits, "|")
}
