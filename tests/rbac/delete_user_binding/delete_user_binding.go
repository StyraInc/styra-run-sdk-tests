package delete_user_binding

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

func deleteUserBinding() test.Test {
	apiResponse := imap{}

	mockResponse := imap{}

	settings := &test.Settings{
		Name: "delete-user-binding",
		Api: &test.Api{
			Request: &test.Request{
				Path:    fmt.Sprintf(pathSdk, user),
				Method:  http.MethodDelete,
				Cookies: test.AuthzCookie(tenant, subject),
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
					test.CheckRequestMethod(http.MethodDelete),
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
			deleteUserBinding(),
		}
	}
}
