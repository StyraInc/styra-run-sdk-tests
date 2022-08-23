package get_roles

import (
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	url       = "/roles"
	pathAuthz = "/data/rbac/manage/allow"
	pathRoles = "/data/rbac/roles"
	tenant    = "acmecorp"
	subject   = "alice"
)

type smap map[string]interface{}
type slist []interface{}

func getRoles() test.Test {
	apiResponse := smap{
		"result": slist{
			"ADMIN",
			"VIEWER",
		},
	}

	authzRequest := smap{
		"input": smap{
			"tenant":  "acmecorp",
			"subject": "alice",
		},
	}

	authzResponse := smap{
		"result": true,
	}

	settings := &test.Settings{
		Name: "get-roles",
		Api: &test.Api{
			Request: &test.Request{
				Path:    url,
				Method:  http.MethodGet,
				Cookies: test.AuthzCookie(tenant, subject),
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			pathAuthz: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPost),
					test.CheckRequestContentType(test.ApplicationJson),
					test.CheckRequestBody(authzRequest),
				},
				Code: 200,
				Body: authzResponse,
			},
			pathRoles: {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodPost),
					test.CheckRequestContentType(test.ApplicationJson),
				},
				Code: 200,
				Body: apiResponse,
			},
		},
	}

	return test.New(settings)
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			getRoles(),
		}
	}
}
