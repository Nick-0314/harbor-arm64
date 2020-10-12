package notification

import (
	"github.com/goharbor/harbor/src/jobservice/job/impl"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMaxFails(t *testing.T) {
	rep := &WebhookJob{}
	// test default max fails
	assert.Equal(t, uint(10), rep.MaxFails())

	// test user defined max fails
	_ = os.Setenv(maxFails, "15")
	assert.Equal(t, uint(15), rep.MaxFails())

	// test user defined wrong max fails
	_ = os.Setenv(maxFails, "abc")
	assert.Equal(t, uint(10), rep.MaxFails())
}

func TestShouldRetry(t *testing.T) {
	rep := &WebhookJob{}
	assert.True(t, rep.ShouldRetry())
}

func TestValidate(t *testing.T) {
	rep := &WebhookJob{}
	assert.Nil(t, rep.Validate(nil))
}

func TestRun(t *testing.T) {
	rep := &WebhookJob{}

	// test webhook request
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := ioutil.ReadAll(r.Body)

			// test request method
			assert.Equal(t, http.MethodPost, r.Method)
			// test request header
			assert.Equal(t, "auth_test", r.Header.Get("Authorization"))
			// test request body
			assert.Equal(t, string(body), `{"key": "value"}`)
		}))
	defer ts.Close()
	params := map[string]interface{}{
		"skip_cert_verify": true,
		"payload":          `{"key": "value"}`,
		"address":          ts.URL,
		"auth_header":      "auth_test",
	}
	// test correct webhook response
	assert.Nil(t, rep.Run(&impl.Context{}, params))

	tsWrong := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
	defer tsWrong.Close()
	paramsWrong := map[string]interface{}{
		"skip_cert_verify": true,
		"payload":          `{"key": "value"}`,
		"address":          tsWrong.URL,
		"auth_header":      "auth_test",
	}
	// test incorrect webhook response
	assert.NotNil(t, rep.Run(&impl.Context{}, paramsWrong))
}
