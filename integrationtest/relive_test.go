package integrationtest

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/msproject/relive/testtools"
)

func TestRelive(t *testing.T) {
	var dbconnect string
	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatal("Error getting wd: ", err.Error())
	}
	var test1 = &testtools.Comp{
		Name:       "Top",
		Sequential: true,
		SubN: []testtools.Worker{
			&testtools.Comp{
				Name:       "Compile Deps",
				Sequential: true,
				SubN: []testtools.Worker{
					&testtools.Exec{
						Name:      "Compile relive",
						Token:     "reliveCompile",
						Directory: rootDir,
						Command: func() *exec.Cmd {
							return exec.Command("go", "build", "-v", "-race", "github.com/msproject/relive")
						},
					},
					&testtools.Wait{
						Name:  "Compile relive Wait",
						Token: "reliveCompile",
					},
				},
			},
			&testtools.Comp{
				Name: "Init",
				SubN: []testtools.Worker{
					&testtools.GoFunc{
						Name: "Startup MySQL",
						Func: func(w *testtools.GoFunc) {
							// Stage our DB
							_, dbconnect, err = testtools.RunMySQL("relive-mysql")
							if err != nil {
								t.Fatalf("Could not start mysql %s", err.Error())
							}
							reliveTestCfg.mysqlAccessAddr = dbconnect

						},
					},
					&testtools.Exec{
						Name:      "Start relive",
						Token:     "relive",
						Directory: rootDir,
						Command: func() *exec.Cmd {
							return exec.Command("./relive",
								"-cert", rootDir+"/../relive_cert.pem",
								"-key", rootDir+"/../relive_key.pem")
						},
					},
					&testtools.DelayHealthCheck{
						Name:         "delay health relive",
						PollInterval: time.Second,
						URL:          reliveTestCfg.reliveServerURL + "/health",
						NPolls:       5,
					},
				},
			},
			&testtools.Comp{
				Name: "Test",
				SubN: []testtools.Worker{
					&testtools.Curl{
						Name:          "health API",
						URL:           reliveTestCfg.reliveServerURL + "/health",
						PrintBodyFlag: false,
					},
					&testtools.Curl{
						Name:          "Version API",
						URL:           reliveTestCfg.reliveServerURL + "/version",
						PrintBodyFlag: true,
					},
					&testtools.GoFunc{
						Name: "Setup Test Environment",
						Func: func(w *testtools.GoFunc) {
							w.Err = reliveSetupTestEnv()
							if w.Err != nil {
								t.Fatalf("Failed to setup the test environment, exit now! err: %v", w.Err)
							}
						},
					},
					&testtools.GoFunc{
						Name: "Test Create Account",
						Func: func(w *testtools.GoFunc) {
							w.Err = testCreateAccount()
						},
					},
				},
			},

			&testtools.Comp{
				Name: "Cleanup",
				SubN: []testtools.Worker{
					&testtools.GoFunc{
						Name: "Cleanup MYSQL Tables",
						Func: func(w *testtools.GoFunc) {
							w.Err = cleanupAllTables(dbconnect)
						},
					},
					&testtools.Kill{
						Name:  "Quit relive",
						Token: "relive",
					},
					&testtools.Wait{
						Name:  "Wait kill relive",
						Token: "relive",
					},
				},
			},
		},
	}

	test1.Start()

	if err := test1.Report(0); err != nil {
		t.Fail()
	}
}
