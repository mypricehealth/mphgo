package mph

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatePtr(t *testing.T) {
	d := NewDatePtr(2001, 02, 03)
	assert.Equal(t, &Date{Time: time.Date(2001, 02, 03, 0, 0, 0, 0, time.UTC)}, d)
}

func TestMarshalJSON(t *testing.T) {
	d := NewDate(2020, 1, 1)
	data, err := d.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`"20200101"`), data)

	d = Date{}
	data, err = d.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestUnmarshalJSON(t *testing.T) {
	d := Date{}
	err := d.UnmarshalJSON([]byte(`"20101106"`))
	assert.NoError(t, err)
	assert.Equal(t, NewDate(2010, 11, 6), d)

	d = Date{}
	err = d.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
	assert.Equal(t, Date{}, d)
}
