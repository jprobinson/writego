package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWeighCats(t *testing.T) {
	tests := []struct {
		given      string
		wantWeight string
		wantCode   int
	}{
		{
			"",
			badCatErr + "\n",
			400,
		},
		{
			`{"breed":"Maine Coon","name":"George"}`,
			"16",
			200,
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
		got, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("weighCats(_,%q) encountered unexpected error reading response: %s", test.given, err)
			continue
		}

		if err != nil {
			t.Errorf("weighCats(_,%q) encountered unexpected error parsing response: %s", test.given, err)
			continue
		}

		if string(got) != test.wantWeight {
			t.Errorf("weighCats(_,%q) expected a response of %q; got %q", test.given, test.wantWeight, got)
			continue
		}

	}
}
