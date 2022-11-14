package list_user_bindings

import (
	"fmt"
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/user_bindings_all"
	pathMock = "/data/rbac/user_bindings/%s"
	tenant   = "acmecorp"
	subject  = "alice"
)

type imap map[string]interface{}
type ilist []interface{}

func listUserBindings() test.Test {
	apiResponse := imap{
		"result": ilist{
			imap{
				"id": "alice",
				"roles": ilist{
					"ADMIN",
				},
			},
			imap{
				"id": "bob",
				"roles": ilist{
					"VIEWER",
				},
			},
			imap{
				"id": "bryan",
				"roles": ilist{
					"VIEWER",
				},
			},
			imap{
				"id": "emily",
				"roles": ilist{
					"VIEWER",
				},
			},
		},
	}

	mockResponse := imap{
		"result": imap{
			"alice": ilist{
				"ADMIN",
			},
			"bob": ilist{
				"VIEWER",
			},
			"bryan": ilist{
				"VIEWER",
			},
			"emily": ilist{
				"VIEWER",
			},
		},
	}

	settings := &test.Settings{
		Name: "list-user-bindings",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodGet,
				Cookies: test.AuthzCookie(tenant, subject),
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			test.AuthzPath: test.AuthzMock(tenant, subject, true),
			fmt.Sprintf(pathMock, tenant): {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodGet),
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
			listUserBindings(),
		}
	}
}
