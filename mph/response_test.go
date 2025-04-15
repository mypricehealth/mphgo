package mph

import (
	"encoding/json"
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
		var r Responses[int64]
		err := json.Unmarshal([]byte(`{"results": [1, 2, 3], "status": 200}`), &r)
		require.NoError(t, err)

		assert.Equal(t, Responses[int64]{Results: []int64{1, 2, 3}, StatusCode: 200}, r)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		var r Responses[int64]
		err := json.Unmarshal([]byte(`{"error": {"title": "title", "detail": "detail"}, "status": 500}`), &r)
		require.NoError(t, err)

		assert.Equal(t, Responses[int64]{Error: &ResponseError{Title: "title", Detail: "detail"}, StatusCode: 500}, r)
	})
}
