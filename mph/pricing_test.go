package mph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuralIndicatorUnmarshalJSON(t *testing.T) {
	t.Parallel()
	var r RuralIndicator
	err := r.UnmarshalJSON([]byte(`"invalid value"`))
	require.Error(t, err)

	err = r.UnmarshalJSON([]byte(`"R"`))
	require.NoError(t, err)
	require.Equal(t, RuralIndicatorRural, r)

	err = r.UnmarshalJSON([]byte(`66`))
	require.NoError(t, err)
	require.Equal(t, RuralIndicatorSuperRural, r)

	err = r.UnmarshalJSON([]byte(`82`))
	require.NoError(t, err)
	require.Equal(t, RuralIndicatorRural, r)

	err = r.UnmarshalJSON([]byte(`0`))
	require.NoError(t, err)
	require.Equal(t, RuralIndicatorUrban, r)
}

func TestGetClaimRepricingNote(t *testing.T) {
	t.Parallel()
	p := Pricing{MedicareRepricingNote: "foo"}
	assert.Equal(t, "foo", p.GetRepricingNote())
	p = Pricing{EditDetail: &ClaimEdits{ClaimRejectionReasons: []string{"bar"}}}
	assert.Equal(t, "bar", p.GetRepricingNote())
}

func TestServiceGetRepricingNote(t *testing.T) {
	t.Parallel()
	s := PricedService{
		MedicareRepricingNote: "test",
		EditDetail:            &LineEdits{},
	}
	assert.Equal(t, "test", s.GetRepricingNote())
	s.AllowedRepricingNote = "test2"
	assert.Equal(t, "test2", s.GetRepricingNote())
	s.EditDetail.ProcedureEdits = []string{"test3"}
	assert.Equal(t, "test2. test3", s.GetRepricingNote())
}

func TestInpatientPriceDetailIsEmpty(t *testing.T) {
	t.Parallel()
	p := InpatientPriceDetail{}
	assert.True(t, p.IsEmpty())
	p = InpatientPriceDetail{
		DRG: "test",
	}
	assert.False(t, p.IsEmpty())
}

func TestOutpatientPriceDetailIsEmpty(t *testing.T) {
	t.Parallel()
	p := OutpatientPriceDetail{}
	assert.True(t, p.IsEmpty())
	p = OutpatientPriceDetail{
		OutlierAmount: 12.34,
	}
	assert.False(t, p.IsEmpty())
}

func TestProviderDetailIsEmpty(t *testing.T) {
	t.Parallel()
	p := ProviderDetail{}
	assert.True(t, p.IsEmpty())
	p = ProviderDetail{
		CCN: "test",
	}
	assert.False(t, p.IsEmpty())
}

func TestClaimEditsIsEmpty(t *testing.T) {
	t.Parallel()
	c := ClaimEdits{}
	assert.True(t, c.IsEmpty())
	c = ClaimEdits{
		ClaimDenialReasons: []string{"test"},
	}
	assert.False(t, c.IsEmpty())
}

func TestClaimEditsGetMessage(t *testing.T) {
	t.Parallel()
	var c *ClaimEdits
	assert.Empty(t, c.GetMessage())

	c = &ClaimEdits{
		ClaimDenialReasons: []string{"test1", "test2"},
	}
	assert.Equal(t, "test1|test2", c.GetMessage())
}

func TestLineEditsIsEmpty(t *testing.T) {
	t.Parallel()
	var l *LineEdits
	assert.True(t, l.IsEmpty())

	l = &LineEdits{
		ProcedureEdits: []string{"test"},
	}
	assert.False(t, l.IsEmpty())
}

func TestLineEditsGetMessage(t *testing.T) {
	t.Parallel()
	var l *LineEdits
	assert.Empty(t, l.GetMessage())

	l = &LineEdits{
		ProcedureEdits: []string{"test1", "test2"},
	}
	assert.Equal(t, "test1|test2", l.GetMessage())
}
