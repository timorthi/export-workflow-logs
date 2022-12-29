package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponseBodyByURL(t *testing.T) {
	testFileContents := strings.Repeat("A", 1024*1024) // 1MB

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testFileContents))
	}))
	defer ts.Close()

	buf, err := getResponseBodyByURL(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(buf.String()), testFileContents)
}

func TestGetRequiredEnv(t *testing.T) {
	testCases := []struct {
		desc             string
		envVarNameToTest string
		shouldSucceed    bool
		want             string
	}{
		{
			desc:             "Returns env vars that are set",
			envVarNameToTest: "foo",
			shouldSucceed:    false,
			want:             "env var 'foo' does not exist",
		},
		{
			desc:             "Errors when an env var is not set",
			envVarNameToTest: "bar",
			shouldSucceed:    true,
			want:             "barValue",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.shouldSucceed {
				t.Setenv(tC.envVarNameToTest, tC.want)
			}

			val, err := getRequiredEnv(tC.envVarNameToTest)

			if tC.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, val, tC.want)
			} else {
				assert.ErrorContains(t, err, tC.want)
			}
		})
	}
}
