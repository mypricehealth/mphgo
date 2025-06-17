package mph

import (
	"encoding/json"

	"braces.dev/errtrace"
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
	ClaimID               string                `json:"claimID,omitempty"`               // The unique identifier for the claim (copied from input)
	MedicareAmount        float64               `json:"medicareAmount,omitempty"`        // The amount Medicare would pay for the service
	AllowedAmount         float64               `json:"allowedAmount,omitempty"`         // The allowed amount based on a contract or RBP pricing
	MedicareRepricingCode ClaimRepricingCode    `json:"medicareRepricingCode,omitempty"` // Explains the methodology used to calculate Medicare (MED or IFO)
	MedicareRepricingNote string                `json:"medicareRepricingNote,omitempty"` // Note explaining approach for pricing or reason for error
	NetworkCode           string                `json:"networkCode,omitempty"`           // The network code used for pricing (is placed into HCP04)
	AllowedRepricingCode  ClaimRepricingCode    `json:"allowedRepricingCode,omitempty"`  // Explains the methodology used to calculate allowed amount (CON, RBP, SCA, or IFO)
	AllowedRepricingNote  string                `json:"allowedRepricingNote,omitempty"`  // Note explaining approach for pricing or reason for error
	MedicareStdDev        float64               `json:"medicareStdDev,omitempty"`        // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource        MedicareSource        `json:"medicareSource,omitempty"`        // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	InpatientPriceDetail  InpatientPriceDetail  `json:"inpatientPriceDetail,omitzero"`   // Details about the inpatient pricing
	OutpatientPriceDetail OutpatientPriceDetail `json:"outpatientPriceDetail,omitzero"`  // Details about the outpatient pricing
	ProviderDetail        ProviderDetail        `json:"providerDetail,omitzero"`         // The provider details used when pricing the claim
	EditDetail            ClaimEdits            `json:"editDetail,omitzero"`             // Errors which cause the claim to be denied, rejected, suspended, or returned to the provider
	PricerResult          string                `json:"pricerResult,omitempty"`          // Pricer return details
	Services              []PricedService       `json:"services,omitempty"`              // Pricing for each service line on the claim
	EditError             *ResponseError        `json:"error,omitempty"`                 // An error that occurred during some step of the pricing process
}

// PricedService contains the results of a pricing request for a single service line.
type PricedService struct {
	LineNumber                  string                  `json:"lineNumber,omitempty"`             // Number of the service line item (copied from input)
	ProviderDetail              *ProviderDetail         `json:"providerDetail,omitempty"`         // Provider Details used when pricing the service if different than the claim
	MedicareAmount              float64                 `json:"medicareAmount,omitempty"`         // Amount Medicare would pay for the service
	AllowedAmount               float64                 `json:"allowedAmount,omitempty"`          // Allowed amount based on a contract or RBP pricing
	MedicareRepricingCode       LineRepricingCode       `json:"medicareRepricingCode,omitempty"`  // Explains the methodology used to calculate Medicare
	MedicareRepricingNote       string                  `json:"medicareRepricingNote,omitempty"`  // Note explaining approach for pricing or reason for error
	NetworkCode                 string                  `json:"networkCode,omitempty"`            // The network code used for pricing (is placed into HCP04)
	AllowedRepricingCode        LineRepricingCode       `json:"allowedRepricingCode,omitempty"`   // Explains the methodology used to calculate allowed amount
	AllowedRepricingNote        string                  `json:"allowedRepricingNote,omitempty"`   // Note explaining approach for pricing or reason for error
	AllowedRepricingFormula     AllowedRepricingFormula `json:"allowedRepricingFormula,omitzero"` // Formula used to calculate the allowed amount
	TechnicalComponentAmount    float64                 `json:"tcAmount,omitempty"`               // Amount Medicare would pay for the technical component
	ProfessionalComponentAmount float64                 `json:"pcAmount,omitempty"`               // Amount Medicare would pay for the professional component
	MedicareStdDev              float64                 `json:"medicareStdDev,omitempty"`         // Standard deviation of the estimated Medicare amount (estimates service only)
	MedicareSource              MedicareSource          `json:"medicareSource,omitempty"`         // Source of the Medicare amount (e.g. physician fee schedule, OPPS, etc.)
	PricerResult                string                  `json:"pricerResult,omitempty"`           // Pricing service return details
	StatusIndicator             string                  `json:"statusIndicator,omitempty"`        // (outpatient + professional) Code which gives more detail about how Medicare pays for the service
	PaymentIndicator            string                  `json:"paymentIndicator,omitempty"`       // (outpatient only) Text which explains the type of payment for Medicare
	DiscountFormula             string                  `json:"discountFormula,omitempty"`        // (outpatient only) The multi-procedure discount formula used to calculate the allowed amount
	PackagingFlag               string                  `json:"packagingFlag,omitempty"`          // (outpatient only) Indicates if the service is packaged and the reason for packaging
	PaymentMethod               string                  `json:"paymentMethod,omitempty"`          // (outpatient only) The method used to calculate the allowed amount
	PaymentAPC                  string                  `json:"paymentAPC,omitempty"`             // (outpatient only) Ambulatory Payment Classification code used for payment
	EditDetail                  *LineEdits              `json:"editDetail,omitempty"`             // Errors which cause the line item to be unable to be priced
}

type AllowedRepricingFormula struct {
	MedicarePercent float64 `json:"medicarePercent,omitempty"` // Percentage of the Medicare amount used to calculate the allowed amount
	BilledPercent   float64 `json:"billedPercent,omitempty"`   // Percentage of the billed amount used to calculate the allowed amount
	FixedAmount     float64 `json:"fixedAmount,omitempty"`     // Fixed amount used as the allowed amount
}

// InpatientPriceDetail contains pricing details for an inpatient claim.
type InpatientPriceDetail struct {
	DRG                            string  `json:"drg,omitempty"                            db:"inpatient_drg"`                               // Diagnosis Related Group (DRG) code used to price the claim
	DRGAmount                      float64 `json:"drgAmount,omitempty"                      db:"inpatient_drg_amount"`                        // Amount Medicare would pay for the DRG
	PassthroughAmount              float64 `json:"passthroughAmount,omitempty"              db:"inpatient_passthrough_amount"`                // Per diem amount to cover capital-related costs, direct medical education, and other costs
	OutlierAmount                  float64 `json:"outlierAmount,omitempty"                  db:"inpatient_outlier_amount"`                    // Additional amount paid for high cost cases
	IndirectMedicalEducationAmount float64 `json:"indirectMedicalEducationAmount,omitempty" db:"inpatient_indirect_medical_education_amount"` // Additional amount paid for teaching hospitals
	DisproportionateShareAmount    float64 `json:"disproportionateShareAmount,omitempty"    db:"inpatient_disproportionate_share_amount"`     // Additional amount paid for hospitals with a high number of low-income patients
	UncompensatedCareAmount        float64 `json:"uncompensatedCareAmount,omitempty"        db:"inpatient_uncompensated_care_amount"`         // Additional amount paid for patients who are unable to pay for their care
	ReadmissionAdjustmentAmount    float64 `json:"readmissionAdjustmentAmount,omitempty"    db:"inpatient_readmission_adjustment_amount"`     // Adjustment amount for hospitals with high readmission rates
	ValueBasedPurchasingAmount     float64 `json:"valueBasedPurchasingAmount,omitempty"     db:"inpatient_value_based_purchasing_amount"`     // Adjustment for hospitals based on quality measures
	WageIndex                      float64 `json:"wageIndex,omitempty"                      db:"inpatient_wage_index"`                        // Wage index used for geographic adjustment
}

func (i InpatientPriceDetail) IsEmpty() bool {
	return i.DRG == "" && i.DRGAmount == 0 && i.PassthroughAmount == 0 && i.OutlierAmount == 0 && i.IndirectMedicalEducationAmount == 0 && i.DisproportionateShareAmount == 0 &&
		i.UncompensatedCareAmount == 0 && i.ReadmissionAdjustmentAmount == 0 && i.ValueBasedPurchasingAmount == 0 && i.WageIndex == 0
}

// OutpatientPriceDetail contains pricing details for an outpatient claim.
type OutpatientPriceDetail struct {
	OutlierAmount                         float64 `json:"outlierAmount,omitempty"                         db:"outpatient_outlier_amount"`                              // Additional amount paid for high cost cases
	FirstPassthroughDrugOffsetAmount      float64 `json:"firstPassthroughDrugOffsetAmount,omitempty"      db:"outpatient_first_passthrough_drug_offset_amount"`        // Amount built into the APC payment for certain drugs
	SecondPassthroughDrugOffsetAmount     float64 `json:"secondPassthroughDrugOffsetAmount,omitempty"     db:"outpatient_second_passthrough_drug_offset_amount"`       // Amount built into the APC payment for certain drugs
	ThirdPassthroughDrugOffsetAmount      float64 `json:"thirdPassthroughDrugOffsetAmount,omitempty"      db:"outpatient_third_passthrough_drug_offset_amount"`        // Amount built into the APC payment for certain drugs
	FirstDeviceOffsetAmount               float64 `json:"firstDeviceOffsetAmount,omitempty"               db:"outpatient_first_device_offset_amount"`                  // Amount built into the APC payment for certain devices
	SecondDeviceOffsetAmount              float64 `json:"secondDeviceOffsetAmount,omitempty"              db:"outpatient_second_device_offset_amount"`                 // Amount built into the APC payment for certain devices
	FullOrPartialDeviceCreditOffsetAmount float64 `json:"fullOrPartialDeviceCreditOffsetAmount,omitempty" db:"outpatient_full_or_partial_device_credit_offset_amount"` // Credit for devices that are supplied for free or at a reduced cost
	TerminatedDeviceProcedureOffsetAmount float64 `json:"terminatedDeviceProcedureOffsetAmount,omitempty" db:"outpatient_terminated_device_procedure_offset_amount"`   // Credit for devices that are not used due to a terminated procedure
	WageIndex                             float64 `json:"wageIndex,omitempty"                             db:"outpatient_wage_index"`                                  // Wage index used for geographic adjustment
}

func (o OutpatientPriceDetail) IsEmpty() bool {
	return o.OutlierAmount == 0 && o.FirstPassthroughDrugOffsetAmount == 0 && o.SecondPassthroughDrugOffsetAmount == 0 &&
		o.ThirdPassthroughDrugOffsetAmount == 0 && o.FirstDeviceOffsetAmount == 0 && o.SecondDeviceOffsetAmount == 0 &&
		o.FullOrPartialDeviceCreditOffsetAmount == 0 && o.TerminatedDeviceProcedureOffsetAmount == 0 && o.WageIndex == 0
}

// ProviderDetail contains basic information about the provider and/or locality used for pricing
// Not all fields are returned with every pricing request. For example, the CMS Certification
// Number (CCN) is only returned for facilities which have a CCN such as hospitals.
type ProviderDetail struct {
	CCN            string         `json:"ccn,omitempty"`            // CMS Certification Number for the facility
	MAC            uint16         `json:"mac"`                      // Medicare Administrative Contractor number
	Locality       uint8          `json:"locality"`                 // Geographic locality number used for pricing
	GeographicCBSA uint32         `json:"geographicCBSA,omitempty"` // Core-Based Statistical Area (CBSA) number for provider ZIP
	StateCBSA      uint8          `json:"stateCBSA,omitempty"`      // State Core-Based Statistical Area (CBSA) number
	RuralIndicator RuralIndicator `json:"ruralIndicator,omitempty"` // Indicates whether provider is Rural (R), Super Rural (B), or Urban (blank)
	SpecialtyType  string         `json:"specialtyType,omitempty"`  // Medicare provider specialty type
	HospitalType   HospitalType   `json:"hospitalType,omitempty"`   // Type of hospital
}

func (p ProviderDetail) IsEmpty() bool {
	return p.CCN == "" && p.MAC == 0 && p.Locality == 0 && p.GeographicCBSA == 0 && p.StateCBSA == 0 &&
		p.RuralIndicator == "" && p.SpecialtyType == "" && p.HospitalType == ""
}

// ClaimEdits contains errors which cause the claim to be denied, rejected, suspended, or returned to the provider.
type ClaimEdits struct {
	HCP13DenyCode                    string   `json:"hcpDenyCode,omitempty"`                      // The deny code that will be placed into the HCP13 data element for EDI 837 claims
	ClaimOverallDisposition          string   `json:"claimOverallDisposition,omitempty"`          // Overall explanation of why the claim edit failed
	ClaimRejectionDisposition        string   `json:"claimRejectionDisposition,omitempty"`        // Explanation of why the claim was rejected
	ClaimDenialDisposition           string   `json:"claimDenialDisposition,omitempty"`           // Explanation of why the claim was denied
	ClaimReturnToProviderDisposition string   `json:"claimReturnToProviderDisposition,omitempty"` // Explanation of why the claim should be returned to provider
	ClaimSuspensionDisposition       string   `json:"claimSuspensionDisposition,omitempty"`       // Explanation of why the claim was suspended
	LineItemRejectionDisposition     string   `json:"lineItemRejectionDisposition,omitempty"`     // Explanation of why the line item was rejected
	LineItemDenialDisposition        string   `json:"lineItemDenialDisposition,omitempty"`        // Explanation of why the line item was denied
	ClaimRejectionReasons            []string `json:"claimRejectionReasons,omitempty"`            // Detailed reason(s) describing why the claim was rejected
	ClaimDenialReasons               []string `json:"claimDenialReasons,omitempty"`               // Detailed reason(s) describing why the claim was denied
	ClaimReturnToProviderReasons     []string `json:"claimReturnToProviderReasons,omitempty"`     // Detailed reason(s) describing why the claim should be returned to provider
	ClaimSuspensionReasons           []string `json:"claimSuspensionReasons,omitempty"`           // Detailed reason(s) describing why the claim was suspended
	LineItemRejectionReasons         []string `json:"lineItemRejectionReasons,omitempty"`         // Detailed reason(s) describing why the line item was rejected
	LineItemDenialReasons            []string `json:"lineItemDenialReasons,omitempty"`            // Detailed reason(s) describing why the line item was denied
}

func (c ClaimEdits) IsEmpty() bool {
	return c.HCP13DenyCode == "" && c.ClaimOverallDisposition == "" && c.ClaimRejectionDisposition == "" &&
		c.ClaimDenialDisposition == "" && c.ClaimReturnToProviderDisposition == "" && c.ClaimSuspensionDisposition == "" &&
		c.LineItemRejectionDisposition == "" && c.LineItemDenialDisposition == "" && len(c.ClaimRejectionReasons) == 0 &&
		len(c.ClaimDenialReasons) == 0 && len(c.ClaimReturnToProviderReasons) == 0 && len(c.ClaimSuspensionReasons) == 0 &&
		len(c.LineItemRejectionReasons) == 0 && len(c.LineItemDenialReasons) == 0
}

// LineEdits contains errors which cause the line item to be unable to be priced.
type LineEdits struct {
	DenialOrRejectionText string   `json:"denialOrRejectionText,omitempty"` // The overall explanation for why this line item was denied or rejected by the claim editor
	ProcedureEdits        []string `json:"procedureEdits,omitempty"`        // Detailed description of each procedure code edit error (from outpatient editor)
	Modifier1Edits        []string `json:"modifier1Edits,omitempty"`        // Detailed description of each edit error for the first procedure code modifier (from outpatient editor)
	Modifier2Edits        []string `json:"modifier2Edits,omitempty"`        // Detailed description of each edit error for the second procedure code modifier (from outpatient editor)
	Modifier3Edits        []string `json:"modifier3Edits,omitempty"`        // Detailed description of each edit error for the third procedure code modifier (from outpatient editor)
	Modifier4Edits        []string `json:"modifier4Edits,omitempty"`        // Detailed description of each edit error for the fourth procedure code modifier (from outpatient editor)
	Modifier5Edits        []string `json:"modifier5Edits,omitempty"`        // Detailed description of each edit error for the fifth procedure code modifier (from outpatient editor)
	DataEdits             []string `json:"dataEdits,omitempty"`             // Detailed description of each data edit error (from outpatient editor)
	RevenueEdits          []string `json:"revenueEdits,omitempty"`          // Detailed description of each revenue code edit error (from outpatient editor)
}

func (l LineEdits) IsEmpty() bool {
	return l.DenialOrRejectionText == "" && len(l.ProcedureEdits) == 0 && len(l.Modifier1Edits) == 0 &&
		len(l.Modifier2Edits) == 0 && len(l.Modifier3Edits) == 0 && len(l.Modifier4Edits) == 0 &&
		len(l.Modifier5Edits) == 0 && len(l.DataEdits) == 0 && len(l.RevenueEdits) == 0
}
