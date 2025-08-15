package mph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditMessages(t *testing.T) {
	t.Parallel()
	p := Pricing{
		EditDetail: &ClaimEdits{
			ClaimRejectionReasons:        []string{"reason1", "reason2"},
			ClaimReturnToProviderReasons: []string{"reason2", "reason1"},
		},
	}
	assert.Equal(t, []string{"reason1", "reason2"}, p.GetEditMessages())
}

func TestUnmarshalRuralIndicator(t *testing.T) {
	t.Parallel()
	unmarshalRuralIndicator(t, []byte(`"A"`), "", true)
	unmarshalRuralIndicator(t, []byte(`""`), RuralIndicatorUrban, false)
	unmarshalRuralIndicator(t, []byte(`"B"`), RuralIndicatorSuperRural, false)
	unmarshalRuralIndicator(t, []byte(`"R"`), RuralIndicatorRural, false)

	unmarshalRuralIndicator(t, []byte(`0`), RuralIndicatorUrban, false)
	unmarshalRuralIndicator(t, []byte(`66`), RuralIndicatorSuperRural, false)
	unmarshalRuralIndicator(t, []byte(`82`), RuralIndicatorRural, false)
}

func unmarshalRuralIndicator(t *testing.T, data []byte, expected RuralIndicator, isError bool) {
	t.Helper()
	var r RuralIndicator
	err := r.UnmarshalJSON(data)
	assert.Equal(t, isError, err != nil, "Want error: %v, Got: %s", isError, err)
	assert.Equal(t, expected, r)
}
