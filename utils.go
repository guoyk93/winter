package winter

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type checkResult struct {
	failed bool
	sb     *strings.Builder
}

func (c *checkResult) Collect(name string, err error) {
	if c.sb == nil {
		c.sb = &strings.Builder{}
	}
	if c.sb.Len() > 0 {
		c.sb.WriteString("\n")
	}
	c.sb.WriteString(name)
	if err == nil {
		c.sb.WriteString(": OK")
	} else {
		c.failed = true
		c.sb.WriteString(": ")
		c.sb.WriteString(err.Error())
	}
}

func (c *checkResult) Result() (s string, failed bool) {
	if c.sb.Len() == 0 {
		s = "OK"
	} else {
		s = c.sb.String()
	}
	failed = c.failed
	return
}

func flattenSingleSlice[T any](s []T) any {
	if len(s) == 1 {
		return s[0]
	}
	return s
}

func extractRequest(m map[string]any, f map[string][]*multipart.FileHeader, req *http.Request) (err error) {
	// header
	for k, vs := range req.Header {
		k = "header_" + strings.ToLower(strings.ReplaceAll(k, "-", "_"))
		m[k] = flattenSingleSlice(vs)
	}

	// query
	for k, vs := range req.URL.Query() {
		v := flattenSingleSlice(vs)
		m[k] = v
		m["query_"+k] = v
	}

	// body
	var buf []byte

	contentType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))

	if contentType != ContentTypeMultipart {
		if buf, err = io.ReadAll(req.Body); err != nil {
			return
		}

		if len(buf) == 0 {
			return
		}
	}

	switch contentType {
	case ContentTypeTextPlain:
		m["body"] = string(buf)
	case ContentTypeApplicationJSON:
		var j map[string]any
		if err = json.Unmarshal(buf, &j); err != nil {
			return
		}
		for k, v := range j {
			m[k] = v
		}
	case ContentTypeFormURLEncoded:
		var q url.Values
		if q, err = url.ParseQuery(string(buf)); err != nil {
			return
		}
		for k, vs := range q {
			m[k] = flattenSingleSlice(vs)
		}
	case ContentTypeMultipart:
		if err = req.ParseMultipartForm(1024 * 1024 * 10); err != nil {
			return
		}
		for k, v := range req.MultipartForm.Value {
			m[k] = flattenSingleSlice(v)
		}
		for k, v := range req.MultipartForm.File {
			f[k] = v
		}
	default:
		m["body"] = buf
		return
	}

	return
}
