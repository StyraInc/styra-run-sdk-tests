package tests

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
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
	url  = "http://localhost:3000"
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

var sdkProcess *exec.Cmd

func startSdk() bool {
	sdkDir := os.Getenv("SDK_DIR")
	if sdkDir == "" {
		fmt.Println("SDK_DIR environment variable not set")
		return false
	}

	sdkCommand := strings.Split(os.Getenv("SDK_CMD"), " ")
	if len(sdkCommand) < 1 {
		fmt.Println("SDK_CMD environment variable not set")
		return false
	}

	sdkExecutable := sdkCommand[0]
	sdkArguments := sdkCommand[1:]
	fmt.Printf("Starting SDK at %s\n", sdkDir)

	sdkProcess = exec.Command(sdkExecutable, sdkArguments...)
	sdkProcess.Dir = sdkDir
	if err := sdkProcess.Start(); err != nil {
		fmt.Printf("Failed to start SDK process: %v\n", err)
		return false
	}

	var running = false

	fmt.Print("Waiting for SDK")
	for i := 0; i < 20; i++ {
		fmt.Print(".")
		response, err := http.Get("http://localhost:3000/ready")
		if err != nil && !errors.Is(err, syscall.ECONNREFUSED) {
			break
		}
		if response != nil && response.StatusCode == 200 {
			running = true
			break
		}
		time.Sleep(time.Second)
	}

	fmt.Printf(" started: %v\n", running)

	return running
}

func stopSdk() {
	fmt.Println("Stopping SDK")
	if err := sdkProcess.Process.Kill(); err != nil {
		fmt.Printf("Failed to stop SDK process: %v\n", err)
	}
}

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

	if !startSdk() {
		t.Error("SDK never started")
	}

	defer stopSdk()

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
