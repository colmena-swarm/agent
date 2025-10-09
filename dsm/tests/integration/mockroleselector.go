package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

type Expectation struct {
	// What we expect from the client:
	Method    string
	Path      string                       // exact path, e.g. "/v1/widgets"
	Query     url.Values                   // expected query key->[]values (exact match if non-nil)
	Headers   http.Header                  // expected headers (subset match)
	BodyRaw   []byte                       // exact raw body (optional)
	BodyJSON  any                          // compare as JSON structure (optional)
	BodyCheck func(t *testing.T, b []byte) // custom check (optional)

	// What we send back:
	ResponseStatus int
	ResponseHeader http.Header
	ResponseBody   []byte
}

type MockRoleSelector struct {
	t            *testing.T
	srv          *httptest.Server
	expectations []Expectation
	idx          int
}

func NewMockRoleSelector(t *testing.T, expectations ...Expectation) *MockRoleSelector {
	ms := &MockRoleSelector{t: t, expectations: append([]Expectation(nil), expectations...)}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		if ms.idx >= len(ms.expectations) {
			t.Fatalf("unexpected request #%d: %s %s", ms.idx+1, r.Method, r.URL.String())
		}
		exp := ms.expectations[ms.idx]
		ms.idx++

		// Assert method & path
		if exp.Method != "" && r.Method != exp.Method {
			t.Fatalf("expect Method %q, got %q", exp.Method, r.Method)
		}
		if exp.Path != "" && r.URL.Path != exp.Path {
			t.Fatalf("expect Path %q, got %q", exp.Path, r.URL.Path)
		}

		// Assert query (exact match if provided)
		if exp.Query != nil {
			gotQ := r.URL.Query()
			if !valuesEqual(exp.Query, gotQ) {
				t.Fatalf("query mismatch\nexpect: %v\ngot:    %v", exp.Query, gotQ)
			}
		}

		// Assert headers (subset)
		for k, wantVals := range exp.Headers {
			gotVals := r.Header.Values(k)
			if !slicesEq(wantVals, gotVals) {
				t.Fatalf("header %q mismatch\nexpect: %v\ngot:    %v", k, wantVals, gotVals)
			}
		}

		// Read body once
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r.Body)
		_ = r.Body.Close()
		body := buf.Bytes()

		// Assert body
		if exp.BodyRaw != nil && !bytes.Equal(exp.BodyRaw, body) {
			t.Fatalf("raw body mismatch\nexpect: %q\ngot:    %q", string(exp.BodyRaw), string(body))
		}
		if exp.BodyJSON != nil {
			var got any
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("body not valid JSON: %v; got: %q", err, string(body))
			}
			want := normalizeJSON(t, exp.BodyJSON)
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("JSON body mismatch\nexpect: %#v\ngot:    %#v", want, got)
			}
		}
		if exp.BodyCheck != nil {
			exp.BodyCheck(t, body)
		}

		// Send response
		for k, vals := range exp.ResponseHeader {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		if exp.ResponseStatus == 0 {
			exp.ResponseStatus = http.StatusOK
		}
		w.WriteHeader(exp.ResponseStatus)
		if len(exp.ResponseBody) > 0 {
			_, _ = w.Write(exp.ResponseBody)
		}
	})

	ms.srv = httptest.NewServer(h)
	return ms
}

func valuesEqual(a, b url.Values) bool {
	if len(a) != len(b) {
		return false
	}
	for k, va := range a {
		vb, ok := b[k]
		if !ok || !slicesEq(va, vb) {
			return false
		}
	}
	return true
}

func slicesEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// normalizeJSON marshals+unmarshals to get map[string]any with stable types
func normalizeJSON(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("normalizeJSON marshal: %v", err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("normalizeJSON unmarshal: %v", err)
	}
	return out
}

func (m *MockRoleSelector) AddExpectation(e Expectation) {
	m.expectations = append(m.expectations, e)
}

func (m *MockRoleSelector) URL() string { return m.srv.URL }
func (m *MockRoleSelector) Close()      { m.srv.Close() }

// Verify all expectations were hit (call at the end of the test)
func (m *MockRoleSelector) Verify() {
	m.t.Helper()
	if m.idx != len(m.expectations) {
		m.t.Fatalf("not all expectations met: served %d of %d", m.idx, len(m.expectations))
	}
}
