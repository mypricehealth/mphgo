package mph

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseUnmarshalJSON(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var r Response[int64]
		err := json.Unmarshal([]byte(`{"result": 123, "status": 200}`), &r)
		require.NoError(t, err)

		assert.Equal(t, Response[int64]{Result: 123, StatusCode: 200}, r)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var r Response[int64]
		err := json.Unmarshal([]byte(`{"error": {"title": "title", "detail": "detail"}, "status": 500}`), &r)
		require.NoError(t, err)

		assert.Equal(t, Response[int64]{Error: &ResponseError{Title: "title", Detail: "detail"}, StatusCode: 500}, r)
	})
}

func TestResponsesUnmarshalJSON(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		var r ErrorAndResultResponses[testStruct]
		err := json.Unmarshal([]byte(`{"results": [{"result":{"strVal":"1"}},{"result":{"intVal":2}},{"error":{"title":"lorem","detail":"ipsum"}}], "status": 200}`), &r)
		require.NoError(t, err)

		assert.Equal(t, ErrorAndResultResponses[testStruct]{Results: []ErrorAndResult[testStruct]{{Result: testStruct{StrVal: "1"}}, {Result: testStruct{IntVal: 2}}, {Error: &ResponseError{Title: "lorem", Detail: "ipsum"}}}, StatusCode: 200}, r)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var r ErrorAndResultResponses[int64]
		err := json.Unmarshal([]byte(`{"error": {"title": "title", "detail": "detail"}, "status": 500}`), &r)
		require.NoError(t, err)

		assert.Equal(t, ErrorAndResultResponses[int64]{Error: &ResponseError{Title: "title", Detail: "detail"}, StatusCode: 500}, r)
	})
}

type testStruct struct {
	StrVal string `json:"strVal"`
	IntVal int    `json:"intVal"`
}

func TestErrorAndResultResponsesGetError(t *testing.T) {
	t.Parallel()
	v := ErrorAndResultResponses[testStruct]{}
	assert.Nil(t, v.GetError())

	v = ErrorAndResultResponses[testStruct]{
		StatusCode: http.StatusBadRequest,
		Error: &ResponseError{
			Title:  "foo",
			Detail: "bar",
		},
	}

	err := v.GetError()
	require.NotNil(t, err)
	assert.Equal(t, "foo", err.Title)
	assert.Equal(t, "bar", err.Detail.Error())
	assert.Equal(t, http.StatusBadRequest, err.ErrorCode)
}

func TestErrorAndResultResponsesUnwrap(t *testing.T) {
	t.Parallel()
	err := &ResponseError{Title: "foo", Detail: "bar"}
	v := ErrorAndResultResponses[testStruct]{Error: err}
	gotResult, gotErr := v.Unwrap()
	require.NotNil(t, gotErr)
	assert.Equal(t, "foo", gotErr.Title)
	assert.Equal(t, "bar", gotErr.Detail.Error())
	assert.Empty(t, gotResult)

	results := []ErrorAndResult[testStruct]{{}}
	v = ErrorAndResultResponses[testStruct]{Results: results}
	gotResult, gotErr = v.Unwrap()
	assert.Nil(t, gotErr)
	assert.Equal(t, results, gotResult)
}

func TestErrorAndResultMarshal(t *testing.T) {
	t.Parallel()
	v := ErrorAndResult[testStruct]{
		Error: &ResponseError{
			Title:  "foo",
			Detail: "bar",
		},
		Result: testStruct{
			StrVal: "baz",
			IntVal: 42,
		},
	}
	data, err := json.Marshal(v)
	require.NoError(t, err)
	assert.JSONEq(t, `{"error":{"title":"foo","detail":"bar"},"result":{"strVal":"baz","intVal":42}}`, string(data))
}

func TestErrorAndResultUnmarshal(t *testing.T) {
	t.Parallel()
	data := []byte(`{"error":{"title":"foo","detail":"bar"},"result":{"strVal":"baz","intVal":42}}`)
	var v ErrorAndResult[testStruct]
	err := json.Unmarshal(data, &v)
	require.NoError(t, err)
	assert.Equal(t, "foo", v.Error.Title)
	assert.Equal(t, "bar", v.Error.Detail)
	assert.Equal(t, "baz", v.Result.StrVal)
	assert.Equal(t, 42, v.Result.IntVal)
}

func TestErrorAndResultUnwrap(t *testing.T) {
	t.Parallel()

	res := testStruct{StrVal: "baz", IntVal: 42}
	expected := &ResponseError{"foo", "bar"}
	v := ErrorAndResult[testStruct]{Result: res, Error: expected}
	result, err := v.Unwrap()
	assert.Equal(t, res, result)
	assert.Equal(t, expected, err)
}
