package list_user_bindings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	pathSdk  = "/user_bindings"
	pathMock = "/data/rbac/user_bindings/%s/%s"
	tenant   = "acmecorp"
	subject  = "alice"
	page     = 2
)

type imap map[string]interface{}
type ilist []interface{}

func listUserBindings() test.Test {
	apiResponse := imap{
		"result": ilist{
			imap{
				"id": "emily",
				"roles": ilist{
					"ADMIN",
					"VIEWER",
				},
			},
			imap{
				"id": "harold",
				"roles": ilist{
					"VIEWER",
				},
			},
			imap{
				"id":    "vivian",
				"roles": ilist{},
			},
		},
		"page": imap{
			"index": 2,
			"total": 2,
		},
	}

	settings := &test.Settings{
		Name: "list-user-bindings",
		Api: &test.Api{
			Request: &test.Request{
				Path:    pathSdk,
				Method:  http.MethodGet,
				Cookies: test.AuthzCookie(tenant, subject),
				Queries: map[string]string{
					"page": fmt.Sprintf("%d", page),
				},
			},
			Checks: []test.CheckResponse{
				test.CheckResponseCode(200),
				test.CheckResponseBody(apiResponse),
			},
		},
		Mocks: map[string]*test.Mock{
			test.AuthzPath: test.AuthzMock(tenant, subject, true),
			fmt.Sprintf(pathMock, tenant, "emily"): {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodGet),
				},
				Response: bindingResponse([]string{"ADMIN", "VIEWER"}),
			},
			fmt.Sprintf(pathMock, tenant, "harold"): {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodGet),
				},
				Response: bindingResponse([]string{"VIEWER"}),
			},
			fmt.Sprintf(pathMock, tenant, "vivian"): {
				Checks: []test.CheckRequest{
					test.CheckRequestMethod(http.MethodGet),
				},
				Response: test.DefaultResponse(404, imap{}),
			},
		},
	}

	return test.New(settings)
}

func bindingResponse(roles []string) test.EmitResponse {
	return func(w http.ResponseWriter, r *test.MockRequest) error {
		body := imap{
			"result": roles,
		}
		if bytes, err := json.Marshal(body); err != nil {
			test.InternalServerError(w)

			return err
		} else {
			w.WriteHeader(200)
			w.Header().Set(test.ContentTypeHeader, test.ApplicationJson)

			if _, err := w.Write(bytes); err != nil {
				test.InternalServerError(w)

				return err
			}
		}

		return nil
	}
}

func New() test.Factory {
	return func() []test.Test {
		return []test.Test{
			listUserBindings(),
		}
	}
}
