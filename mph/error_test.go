package mph

import (
	"testing"

	"braces.dev/errtrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Parallel()

	e := &Error{Title: "title", Detail: errtrace.Errorf("detail")}
	assert.Equal(t, "title: detail", e.Error())

	e = nil
	assert.Empty(t, e.Error())
}

func TestErrorMarshalJSON(t *testing.T) {
	t.Parallel()

	e := &Error{Title: "title", Detail: errtrace.Errorf("detail"), ErrorCode: 500}
	got, err := e.MarshalJSON()
	want := `{"title":"title","errorCode":500}`
	require.NoError(t, err)
	assert.Equal(t, want, string(got))

	e = nil
	_, err = e.MarshalJSON()
	require.NoError(t, err)

	e = &Error{Title: "title", Detail: errtrace.Errorf("detail"), ErrorCode: 400}
	assert.Panics(t, func() {
		_, _ = e.MarshalJSON()
	})
}

func TestErrorUnmarshalJSON(t *testing.T) {
	t.Parallel()

	e := &Error{}
	data := []byte(`{"title":"title","errorCode":300}`)
	err := e.UnmarshalJSON(data)
	require.NoError(t, err)
	expected := &Error{Title: "title", ErrorCode: 300}
	assert.Equal(t, expected, e)
}

func TestParseResponseError(t *testing.T) {
	msg := "Title: detail"
	got := ParseResponseError(msg)
	want := &ResponseError{Title: "Title", Detail: "detail"}
	assert.Equal(t, want, got)

	msg = "No colon here"
	got = ParseResponseError(msg)
	want = &ResponseError{Title: "Error", Detail: "No colon here"}
	assert.Equal(t, want, got)

	msg = ""
	assert.Nil(t, ParseResponseError(msg))
}
