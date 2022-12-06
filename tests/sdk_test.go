package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/styrainc/styra-run-sdk-tests/tests/client/batch_query"
	"github.com/styrainc/styra-run-sdk-tests/tests/client/check"
	"github.com/styrainc/styra-run-sdk-tests/tests/client/delete_data"
	"github.com/styrainc/styra-run-sdk-tests/tests/client/get_data"
	"github.com/styrainc/styra-run-sdk-tests/tests/client/put_data"
	"github.com/styrainc/styra-run-sdk-tests/tests/client/query"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/delete_user_binding"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/get_roles"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/get_user_binding"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/list_user_bindings"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/list_user_bindings_all"
	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/put_user_binding"
	"github.com/styrainc/styra-run-sdk-tests/tests/server"
	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	port = 4000
	url  = "http://127.0.0.1:3000"
)

var (
	factories = []test.Factory{
		get_data.New(),
		put_data.New(),
		delete_data.New(),
		query.New(),
		check.New(),
		batch_query.New(),
		get_roles.New(),
		list_user_bindings_all.New(),
		list_user_bindings.New(),
		get_user_binding.New(),
		put_user_binding.New(),
		delete_user_binding.New(),
	}
)

func TestSdk(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	server := server.NewWebServer(
		&server.Settings{
			Port: port,
		},
	)

	go func() {
		defer wg.Done()

		if err := server.Listen(); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)

	for _, factory := range factories {
		for _, test := range factory() {
			server.SetTest(test)

			for _, err := range test.Run(url) {
				t.Errorf("%s: %v", test.Name(), err)
			}
		}
	}

	if err := server.Shutdown(); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
