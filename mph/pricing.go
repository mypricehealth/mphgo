package mph

import (
	"encoding/json"
	"strings"

	"braces.dev/errtrace"
	"github.com/mypricehealth/mphgo/set"
)

type ClaimRepricingCode string
type LineRepricingCode string
type HospitalType string
type RuralIndicator string
type MedicareSource string

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

// Pricing contains the results of a pricing request
type Pricing struct {
	ClaimID               string                `json:"claimID,omitzero"`               // The unique identifier for the claim (copied from input)
	MedicareAmount        float64               `json:"medicareAmount,omitzero"`        // The amount Medicare would pay for the service
	AllowedAmount         float64               `json:"allowedAmount,omitzero"`         // The allowed amount based on a contract or RBP pricing
	MedicareRepricingCode ClaimRepricingCode    `json:"medicareRepricingCode,omitzero"` // Explains the methodology used to calculate Medicare (MED or IFO)
	MedicareRepricingNote string                `json:"medicareRepricingNote,omitzero"` // Note explaining approach for pricing or reason for error
	NetworkCode           string                `json:"networkCode,omitzero"`           // The network code used for pricing (is placed into HCP04)
	AllowedRepricingCode  ClaimRepricingCode    `json:"allowedRepricingCode,omitzero"`  // Explains the methodology used to calculate allowed amount (CON, RBP, SCA, or IFO)
	AllowedRepricingNote  string                `json:"allowedRepricingNote,omitzero"`  // Note explaining approach for pricing or reason for error
	MedicareStdDev        float64               `json:"medicareStdDev,omitzero"`        // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource        MedicareSource        `json:"medicareSource,omitzero"`        // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	InpatientPriceDetail  InpatientPriceDetail  `json:"inpatientPriceDetail,omitzero"`  // Details about the inpatient pricing
	OutpatientPriceDetail OutpatientPriceDetail `json:"outpatientPriceDetail,omitzero"` // Details about the outpatient pricing
	ProviderDetail        ProviderDetail        `json:"providerDetail,omitzero"`        // The provider details used when pricing the claim
	EditDetail            ClaimEdits            `json:"editDetail,omitzero"`            // Errors which cause the claim to be denied, rejected, suspended, or returned to the provider
	PricerResult          string                `json:"pricerResult,omitzero"`          // Pricer return details
	Services              []PricedService       `json:"services,omitzero"`              // Pricing for each service line on the claim
	EditError             *ResponseError        `json:"error,omitzero"`                 // An error that occurred during some step of the pricing process
}

// PricedService contains the results of a pricing request for a single service line.
type PricedService struct {
	LineNumber                    string                  `json:"lineNumber,omitzero"                    db:"-"`                                  // Number of the service line item (copied from input)
	ProviderDetail                ProviderDetail          `json:"providerDetail,omitzero"                 db:",inline"`                           // Provider Details used when pricing the service if different than the claim
	MedicareAmount                float64                 `json:"medicareAmount,omitzero"                db:"medicare_amount"`                    // Amount Medicare would pay for the service
	AllowedAmount                 float64                 `json:"allowedAmount,omitzero"                 db:"allowed_amount"`                     // Allowed amount based on a contract or RBP pricing
	MedicareRepricingCode         LineRepricingCode       `json:"medicareRepricingCode,omitzero"         db:"medicare_repricing_code"`            // Explains the methodology used to calculate Medicare
	MedicareRepricingNote         string                  `json:"medicareRepricingNote,omitzero"         db:"medicare_repricing_note"`            // Note explaining approach for pricing or reason for error
	NetworkCode                   string                  `json:"networkCode,omitzero"                   db:"network_code"`                       // The network code used for pricing (is placed into HCP04)
	AllowedRepricingCode          LineRepricingCode       `json:"allowedRepricingCode,omitzero"          db:"allowed_repricing_code"`             // Explains the methodology used to calculate allowed amount
	AllowedRepricingNote          string                  `json:"allowedRepricingNote,omitzero"          db:"allowed_repricing_note"`             // Note explaining approach for pricing or reason for error
	AllowedRepricingFormula       AllowedRepricingFormula `json:"allowedRepricingFormula,omitzero"        db:",inline"`                           // Formula used to calculate the allowed amount
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
	EditDetail                    LineEdits               `json:"editDetail,omitzero"                     db:",inline"`                           // Errors which cause the line item to be unable to be priced
}

type AllowedRepricingFormula struct {
	MedicarePercent float64 `json:"medicarePercent,omitzero" db:"allowed_repricing_formula_medicare_percent"` // Percentage of the Medicare amount used to calculate the allowed amount
	BilledPercent   float64 `json:"billedPercent,omitzero"   db:"allowed_repricing_formula_billed_percent"`   // Percentage of the billed amount used to calculate the allowed amount
	FixedAmount     float64 `json:"fixedAmount,omitzero"     db:"allowed_repricing_formula_fixed_amount"`     // Fixed amount used as the allowed amount
}

func (s PricedService) GetRepricingNote() string {
	var buf strings.Builder
	if s.AllowedRepricingNote != "" {
		buf.WriteString(s.AllowedRepricingNote)
	} else if s.MedicareRepricingNote != "" {
		buf.WriteString(s.MedicareRepricingNote)
	}
	if edit := s.EditDetail; !edit.IsEmpty() {
		if buf.Len() > 0 {
			buf.WriteString(". ")
		}
		buf.WriteString(edit.GetMessage())
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
	CCN            string         `json:"ccn,omitzero"`            // CMS Certification Number for the facility
	MAC            uint16         `json:"mac"`                     // Medicare Administrative Contractor number
	Locality       uint8          `json:"locality"`                // Geographic locality number used for pricing
	GeographicCBSA uint32         `json:"geographicCBSA,omitzero"` // Core-Based Statistical Area (CBSA) number for provider ZIP
	StateCBSA      uint8          `json:"stateCBSA,omitzero"`      // State Core-Based Statistical Area (CBSA) number
	RuralIndicator RuralIndicator `json:"ruralIndicator,omitzero"` // Indicates whether provider is Rural (R), Super Rural (B), or Urban (blank)
	SpecialtyType  string         `json:"specialtyType,omitzero"`  // Medicare provider specialty type
	HospitalType   HospitalType   `json:"hospitalType,omitzero"`   // Type of hospital
}

func (p ProviderDetail) IsEmpty() bool {
	return p == ProviderDetail{}
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
