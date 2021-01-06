package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApsec_ListPolicyApiEndpoints(t *testing.T) {

	result := GetPolicyApiEndpointsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestPolicyApiEndpoints/PolicyApiEndpoints.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetPolicyApiEndpointsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPolicyApiEndpointsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetPolicyApiEndpointsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-endpoints",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPolicyApiEndpointsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching PolicyApiEndpoints",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-endpoints",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching PolicyApiEndpoints",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPolicyApiEndpoints(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test ApiEndpoints
func TestAppSec_GetPolicyApiEndpoints(t *testing.T) {

	result := GetPolicyApiEndpointsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestApiEndpoints/ApiEndpoints.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetPolicyApiEndpointsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetPolicyApiEndpointsResponse
		withError        error
	}{
		"200 OK": {
			params: GetPolicyApiEndpointsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-endpoints",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetPolicyApiEndpointsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/api-endpoints",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetPolicyApiEndpoints(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
