package check

import (
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/query/tickets/resolve/allow"
	pathMock = "/data/tickets/resolve/allow"
	tenant   = "acmecorp"
	subject  = "alice"
)

type imap map[string]interface{}

func check() test.Test {
	apiRequest := imap{
		"input": imap{
			"tenant":  tenant,
			"subject": subject,
		},
	}

	apiResponse := imap{
		"result": true,
	}

	mockResponse := apiResponse

	settings := &test.Settings{
		Name: "check",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodPost,
				Headers: test.DefaultContentTypeHeader,
				Cookies: test.AuthzCookie(tenant, subject),
				Body:    apiRequest,
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			pathMock: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPost),
				},
				Response: test.DefaultResponse(200, mockResponse),
			},
		},
	}

	return test.New(settings)
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			check(),
		}
	}
}
