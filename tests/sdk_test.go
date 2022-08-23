package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/styrainc/styra-run-sdk-tests/tests/rbac/get_roles"
	"github.com/styrainc/styra-run-sdk-tests/tests/server"
	"github.com/styrainc/styra-run-sdk-tests/tests/test"
)

const (
	port = 4000
	url  = "http://localhost:3000"
)

var (
	factories = []test.Factory{
		get_roles.New(),
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
				t.Error(err)
			}
		}
	}

	if err := server.Shutdown(); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
