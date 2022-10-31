package put_user_binding

import (
	"fmt"
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/user_bindings/%s"
	pathMock = "/data/rbac/user_bindings/%s/%s"
	tenant   = "acmecorp"
	subject  = "alice"
	user     = "cesar"
)

type imap map[string]interface{}
type ilist []interface{}

func putUserBinding() test.Test {
	apiRequest := ilist{
		"ADMIN",
	}

	apiResponse := imap{}

	mockRequest := ilist{
		"ADMIN",
	}

	mockResponse := imap{}

	settings := &test.Settings{
		Name: "put-user-binding",
		Api: &test.Api{
			Request: &test.Request{
				Path:    fmt.Sprintf(pathSdk, user),
				Method:  http.MethodPut,
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
			test.AuthzPath: test.AuthzMock(tenant, subject, true),
			fmt.Sprintf(pathMock, tenant, user): {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPut),
					test.CheckRequestContentType(test.ApplicationJson),
					test.CheckRequestBody(mockRequest),
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
			putUserBinding(),
		}
	}
}
