package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestWeighCats(t *testing.T) {
	tests := []struct {
		given      string
		wantWeight int64
		wantCode   int
		wantErr    bool
	}{
		{
			"",
			0,
			400,
			true,
		},
		{
			`{"breed":"Maine Coon","name":"George"}`,
			16,
			200,
			false,
		},
	}

	for _, test := range tests {
		// build a request with our given body
		req, err := http.NewRequest("GET", "", strings.NewReader(test.given))
		if err != nil {
			t.Errorf("weighCats(_,%q) is an invalid test: %s", test.given, err)
			continue
		}

		// create our new ResponseRecorder for this test
		w := httptest.NewRecorder()

		// run the test!
		weighCat(w, req)

		// inspect our response status code
		if w.Code != test.wantCode {
			t.Errorf("weighCats(_,%q) expected status code %d; got %d", test.given, test.wantCode, w.Code)
		}

		// read the entire response body
		gotBody, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("weighCats(_,%q) encountered unexpected error reading response: %s", test.given, err)
			continue
		}

		// attempt to parse into our expected integer
		got, err := strconv.ParseInt(string(gotBody), 10, 64)
		if test.wantErr {
			if err == nil {
				t.Errorf("weighCats(_,%q) expected error but got none", test.given)
			}
			continue
		}

		if err != nil {
			t.Errorf("weighCats(_,%q) encountered unexpected error parsing response: %s", test.given, err)
			continue
		}

		if got != test.wantWeight {
			t.Errorf("weighCats(_,%q) expected a response of %d; got %d", test.given, test.wantWeight, got)
			continue
		}

	}
}
