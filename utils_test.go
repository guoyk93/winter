package winter

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondInternal(t *testing.T) {
	rw := httptest.NewRecorder()
	internalRespond(rw, "OK", http.StatusTeapot)
	require.Equal(t, rw.Code, http.StatusTeapot)
	require.Equal(t, rw.Body.String(), "OK")
}

func TestFlattenSimpleSlice(t *testing.T) {
	require.Equal(t, "a", flattenSingleSlice([]string{"a"}))
	require.Equal(t, []int{1, 2}, flattenSingleSlice([]int{1, 2}))
}

func TestExtractRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "https://example.com/get?aaa=bbb", nil)

	m := map[string]any{}
	err := extractRequest(m, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`{"hello":"world"}`)))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "application/json;charset=utf-8", "hello": "world", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "application/x-www-form-urlencoded;charset=utf-8", "hello": "world", "query_aaa": "bbb"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "text/plain;charset=utf-8")

	m = map[string]any{}
	err = extractRequest(m, req)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"aaa": "bbb", "header_content_type": "text/plain;charset=utf-8", "query_aaa": "bbb", "text": "hello=world"}, m)

	req = httptest.NewRequest("POST", "https://example.com/post?aaa=bbb", bytes.NewReader([]byte(`hello=world`)))
	req.Header.Set("Content-Type", "application/x-custom")

	m = map[string]any{}
	err = extractRequest(m, req)
	require.Error(t, err)
}
