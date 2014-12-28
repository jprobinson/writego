package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetWeight(t *testing.T) {
	tests := []struct {
		givenResp  string
		wantErr    bool
		wantWeight int64
	}{
		{
			"",
			true,
			0,
		},
		{
			"10",
			false,
			int64(10),
		},
	}

	// setup a test server that will respond with whats been passed to our resps channel
	resps := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := <-resps
		fmt.Fprint(w, resp)
	}))
	defer srv.Close()

	for _, test := range tests {
		// put our resp on the channel
		resps <- test.givenResp
		// make the call against our server and capture the result
		gotWeight, gotErr := getWeight(srv.URL, nil)

		if test.wantErr {
			if gotErr == nil {
				t.Error("weighCat(_,_) expected err but did not get one")
			}
			continue
		}

		if gotErr != nil {
			t.Errorf("weighCat(_,_) expected no error but got one: %s", gotErr)
			continue
		}

		if gotWeight != test.wantWeight {
			t.Errorf("weighCat(_,_) expected %q; got %q", test.wantWeight, gotWeight)
		}
	}
}

func TestCalcAvgWeight(t *testing.T) {
	// create a var to count the number of requests we get
	gotReqs := 0
	// start a test server that will always return 10
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqs++
		fmt.Fprint(w, "10")
	}))
	defer srv.Close()

	// set up 10 lines of data
	testData := "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n"

	// override global vars for testing
	host = srv.URL
	procs = 1
	gotWeight := calcAverageWeight(strings.NewReader(testData))

	expWeight := float32(10)
	if gotWeight != expWeight {
		t.Errorf("calcAverageWeight(_) expected weight %q; got %q", expWeight, gotWeight)
	}
	expReqs := 10
	if gotReqs != expReqs {
		t.Errorf("calcAverageWeight(_) expected %d requests; got %d", expReqs, gotReqs)
	}
}
