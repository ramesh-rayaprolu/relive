package testtools

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

const buflen = 4096
const defaultBufferLimit = 4096

// These types are used to implement the integration test spec
type (
	// Worker interface is an interface that all nodes in the integration test
	// spec must implement.
	Worker interface {
		Report(int) error   // Prints Name and iterates over subnodes
		Go(*sync.WaitGroup) // Performs function of node
	}
	// Comp is a composite node. It groups multiple nodes together
	Comp struct {
		// Name of this step
		Name       string // Name of the step. Documentation only
		Sequential bool   // The nodes in SubN should be done serially
		SubN       []Worker
		Err        error
	}
	// Curl is used to do an http request a la curl
	Curl struct {
		Name          string // Name of the step. Documentation only
		URL           string
		PrintBodyFlag bool
		Resp          *http.Response
		Err           error
	}

	// Delay - like sleep
	Delay struct {
		Name     string // Name of the step. Documentation only
		Duration time.Duration
		Err      error
	}

	// DelayHealthCheck - delay and perform heartbeat until active
	DelayHealthCheck struct {
		Name         string // Name of the step. Documentation only
		URL          string
		Resp         *http.Response
		PollInterval time.Duration
		NPolls       int
		Err          error
	}

	// Wait - wait on a process to complete.
	Wait struct {
		Name    string // Name of the step. Documentation only
		Token   string
		Timeout time.Duration //
		Err     error
	}

	// Kill - kill one process
	Kill struct {
		Name  string // Name of the step. Documentation only
		Token string
		Err   error
	}

	// GoFunc is used to embed
	GoFunc struct {
		Name string    // Name of the step. Documentation only
		Func GoRoutine // Function to invoke
		Skip bool      // Flag to skip a test case, this will be reported SKIPPED in test report
		Err  error
	}

	// GoRoutine - Type of functions in GoFunc
	GoRoutine func(*GoFunc)
	// Exec is used to start a process (ForkExec)
	// token is used to refer to that process (e.g. in Kill)

	stream struct {
		stream      io.ReadCloser
		firstBuf    []byte // Print the first buffer, even in TailMode
		data        chan []byte
		BufferLimit int
		TailMode    bool
		Report      int // 0 = normal, 1 = on error, 2 = off
		Discards    int
	}

	// Exec - exec command
	Exec struct {
		Name             string // Name of the step. Documentation only
		Token            string
		Directory        string
		Precommand       string
		Command          func() *exec.Cmd
		Args             []string
		Sequential       bool // DEPRECATED, remove all references
		SubN             []Worker
		Cmd              *exec.Cmd
		Pid              int // Pid of process
		Err              error
		stdwg            sync.WaitGroup
		KilledExplicitly bool           // Whether killed explicitly by Kill worker
		VerifyFunc       ExecVerifyFunc // Function to use for execution verification

		Stdout stream
		Stderr stream
	}

	// ExecVerifyFunc verification function to use for Exec runs
	ExecVerifyFunc func(exec *Exec, err error) error
)

var (
	//TailReportOnError - report on err
	TailReportOnError = stream{
		Report:   1,
		TailMode: true,
	}
	//TailReportOn - report on err
	TailReportOn = stream{
		TailMode: true,
	}
	//ReportOn - report on err
	ReportOn = stream{}
)

var (
	// Execs - test
	Execs map[string]*Exec
)

func init() {
	Execs = make(map[string]*Exec)
}

// Go ---------------- Delay Methods
func (w *Delay) Go(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s\n", w.Name)
	time.Sleep(w.Duration)
}

func printName(dep int, name string) {
	name += "                                                  "
	fmt.Printf("%s%s: ", indent(dep), name[:50-dep*8])
}

// ReportNode - report method
func ReportNode(dep int, name string, err error) {
	printName(dep, name)
	if err == nil {
		fmt.Printf("OK\n")
	} else {
		fmt.Printf("FAILED: %s\n", err)
	}
}

// ReportNodeSkipped - report method
func ReportNodeSkipped(dep int, name string) {
	printName(dep, name)
	fmt.Printf("SKIPPED\n")
}

// Report - report method
func (w *Delay) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}

// Go ---------------- DelayHealthCheck Methods
func (w *DelayHealthCheck) Go(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s\n", w.Name)

	// Wait 1 PollInterval to avoid clutter in the logs
	time.Sleep(w.PollInterval)
	for i := 0; i < w.NPolls; i++ {
		w.Resp, w.Err = http.Get(w.URL)
		if w.Err != nil {
			fmt.Printf("HealthCheckDelay-ing\n")
			fmt.Printf("bad read, err = %s\n", w.Err)
			time.Sleep(w.PollInterval)
			continue
		}
		defer w.Resp.Body.Close()
		if w.Resp.StatusCode == 200 {
			buf := make([]byte, 4096)

			defer w.Resp.Body.Close()
			n, err := w.Resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("DelayHealthCheck: error reading body %s\n", err)
			} else {
				fmt.Printf("%s\n", string(buf[:n]))
			}
			break
		}
		time.Sleep(w.PollInterval)
	}
	fmt.Printf("HealthCheckDelay Finished(%s)\n", w.Name)
}

// Report - report method
func (w *DelayHealthCheck) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}

// Go ---------------- GoFunc Methods
func (w *GoFunc) Go(wg *sync.WaitGroup) {
	fmt.Printf("GoFunc  %s\n", w.Name)
	defer wg.Done()

	// The Skip flag can be used to skip a test case from running.
	// For these cases the status will be set to SKIPPED in the report.
	if w.Skip == true {
		fmt.Printf("GoFunc %s skipping", w.Name)
		return
	}

	if w.Func == nil {
		w.Err = errors.New("GoFunc without a func attached")
		return
	}
	w.Func(w)
}

// Report - report method
func (w *GoFunc) Report(dep int) error {
	if w.Skip == true {
		ReportNodeSkipped(dep, w.Name)
	} else {
		ReportNode(dep, w.Name, w.Err)
	}
	return w.Err
}

// Go ---------------- Curl Methods
func (w *Curl) Go(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s\n", w.Name)

	w.Resp, w.Err = http.Get(w.URL)
	if w.Err != nil {
		fmt.Printf("Curl.Go: bad http.Get, err = %s\n", w.Err)
		return
	}
	defer w.Resp.Body.Close()
	if w.PrintBodyFlag {
		// If it's -1, we have an unknown content length. Possibly have
		// chunked Transfer Encoding.
		if w.Resp.ContentLength == -1 {
			//Verify that it's chunked.
			isChunked := false
			for _, xfer := range w.Resp.TransferEncoding {
				if xfer == "chunked" {
					isChunked = true
					break
				}
			}
			if isChunked {
				if resp, err := ioutil.ReadAll(w.Resp.Body); err != nil {
					fmt.Printf("Curl:PrintBodyFlag: error reported %s\n", err)
				} else {
					fmt.Printf("%s\n", string(resp))
				}

			} else {
				fmt.Printf("Unknown Transfer encoding present.\n")
			}
		} else {
			buf := make([]byte, 4096)
			_, err := w.Resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("Curl:PrintBodyFlag: error reported %s\n", err)
			} else {
				l := w.Resp.ContentLength
				if l < 0 {
					fmt.Printf("w.Resp.ContentLength negative response\n")
				} else {
					fmt.Printf("%s\n", string(buf[:l]))
				}
			}
		}
	}
}

// Report - report method
func (w *Curl) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}

// Go ---------------- Exec Methods
func (w *Exec) Go(wg *sync.WaitGroup) {
	defer wg.Done()
	Execs[w.Token] = w

	fmt.Printf("Exec.Go(%s)\n", w.Token)

	// Set some default values (16M each buffer)
	if w.Stdout.BufferLimit == 0 {
		w.Stdout.BufferLimit = (defaultBufferLimit)
	}
	if w.Stderr.BufferLimit == 0 {
		w.Stderr.BufferLimit = (defaultBufferLimit)
	}
	w.Err = os.Chdir(w.Directory)
	if w.Err != nil {
		fmt.Printf("Exec.Go(%s): Chdir failed %s\n", w.Token, w.Err)
		return
	}
	w.Cmd = w.Command()

	w.Stdout.stream, w.Err = w.Cmd.StdoutPipe()
	if w.Err != nil {
		fmt.Printf("Exec.Go(%s): w.Cmd.StdoutPipe() failed%s\n", w.Token, w.Err)
		return
	}
	w.Stderr.stream, w.Err = w.Cmd.StderrPipe()
	if w.Err != nil {
		fmt.Printf("Exec.Go(%s): w.Cmd.StderrPipe() failed%s\n", w.Token, w.Err)
		return
	}
	w.Err = w.Cmd.Start()
	if w.Err != nil {
		fmt.Printf("Exec.Go(%s): w.Cmd.Start() failed %s\n", w.Token, w.Err)
		return
	}

	// Start 2 go routines to collect stdout and stderr
	collector := func(tag string, w *Exec, s *stream) {
		defer w.stdwg.Done()
		s.data = make(chan []byte, s.BufferLimit)
		buffer := make([]byte, buflen)
		for {
			// The read will block until there's data to read OR the pipe is
			// closed by calling .Wait()
			errlen, err := io.ReadAtLeast(s.stream, buffer, buflen)
			if errlen > 0 {
				if s.firstBuf == nil {
					s.firstBuf = buffer[0:errlen]
					buffer = make([]byte, buflen)
					continue
				}
				if s.TailMode {
					if len(s.data) == s.BufferLimit {
						<-s.data // discard
						s.Discards++
					}
					s.data <- buffer[0:errlen] // append
				} else {
					if len(s.data) == s.BufferLimit {
						s.Discards++
						continue
					}
					s.data <- buffer[0:errlen] // append
				}
			}
			if err != nil {
				break
			}
			buffer = make([]byte, buflen)
		}
	}
	w.stdwg.Add(2)
	go collector("stdout", w, &w.Stdout)
	go collector("stderr", w, &w.Stderr)

	if w.Err != nil {
		fmt.Printf("Exec.Go(%s):Start failed %s\n", w.Token, w.Err)
		return
	}
}

// Report - report method
func (w *Exec) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}

// VerifyExitCode is a simple utility function to verify an expected exit
// code returned from a Exec.Cmd.Wait() call and return an error if it does not match.
func VerifyExitCode(expected int, err error) error {
	ex := expected & 0xff
	if ee, ok := err.(*exec.ExitError); ok != false {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok != false {
			ec := ws.ExitStatus()
			if ec != ex {
				return fmt.Errorf("Executable exit code %d does not match expected error code %d",
					ec, ex)
			}
		} else {
			return fmt.Errorf("Unable to verify executable exit code: failed to retrieve WaitStatus")
		}
	} else {
		return fmt.Errorf("Unable to verify executable exit code: failed to retrieve ExitError")
	}

	return nil
}

// Go ---------------- Comp Methods
func (w *Comp) Go(wg *sync.WaitGroup) {
	fmt.Printf("Go  %s\n", w.Name)
	defer wg.Done()

	if len(w.SubN) > 0 {
		var wg2 sync.WaitGroup
		for _, sub := range w.SubN {
			wg2.Add(1)
			if w.Sequential {
				sub.Go(&wg2)
			} else {
				// TODO put back the 'go'
				sub.Go(&wg2)
			}
		}
		wg2.Wait()
	}
	fmt.Printf("End %s\n", w.Name)
}

//Report - report method
func (w *Comp) Report(dep int) error {
	var err error
	ReportNode(dep, w.Name, w.Err)
	for _, sub := range w.SubN {
		e1 := sub.Report(dep + 1)
		if e1 != nil {
			err = e1
		}
	}

	if w.Err != nil {
		return w.Err
	}

	return err
}

func indent(d int) (tabs string) {
	for i := 0; i < d; i++ {
		tabs += "\t"
	}
	return
}

//Start - start
func (w *Comp) Start() {
	var wg sync.WaitGroup

	wg.Add(1)
	w.Go(&wg)
	if w.Err != nil {
		return
	}
	wg.Wait()
}

// Go ----------- Wait Methods
func (w *Wait) Go(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Go %s, token = %s\n", w.Name, w.Token)

	runner := Execs[w.Token]
	if runner == nil {
		fmt.Printf("Wait.Go: Token not found\n")
		return
	}
	var timer *time.Timer
	if w.Timeout > 0 {
		timer = time.AfterFunc(w.Timeout, func() {
			fmt.Printf("Wait timeout reached, killing process\n")
			runner.Cmd.Process.Kill()
		})
	}

	runner.stdwg.Wait()

	err := runner.Cmd.Wait()

	if timer != nil {
		timer.Stop()
	}

	if err != nil {
		if runner.KilledExplicitly {
			fmt.Printf("Wait.Go: Cmd.Wait() returns error %s. Ignore error; killed explicitly", err)
		} else {
			if runner.VerifyFunc != nil {
				w.Err = runner.VerifyFunc(runner, err)
			} else {
				w.Err = err
			}
		}
	}
	strfmt := "\n########################## START      - (%s) -  ####################################\n"
	onefmt := "\n########################## FIRST BUFFER (%s) -  ####################################\n"
	endfmt := "\n########################## END          (%s) -  ####################################\n"
	tailfmt := "\n+++++++++++ TailMode Discarded %d %s buffers +++++++++++++\n"
	headfmt := "\n+++++++++++ HeadMode Discarded %d %s buffers +++++++++++++\n"

	reporter := func(tag string, s *stream) {
		if s.Report == 0 ||
			((s.Report == 1) && (err != nil)) {

			fmt.Printf(strfmt, tag)
			if s.firstBuf != nil {
				fmt.Printf(onefmt, tag)
				fmt.Print(string(s.firstBuf))
			}
			if s.TailMode && s.Discards > 0 {
				fmt.Printf(tailfmt, s.Discards, tag)
			}
			for len(s.data) > 0 {
				buf := <-s.data
				fmt.Print(string(buf))
			}
			if !s.TailMode && s.Discards > 0 {
				fmt.Printf(headfmt, s.Discards, tag)
			}
			fmt.Printf(endfmt, tag)
		}
	}
	reporter("stdout", &runner.Stdout)
	reporter("stderr", &runner.Stderr)
}

//Report - report
func (w *Wait) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}

// Go ----------- Kill Methods
func (w *Kill) Go(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Go %s, token = %s\n", w.Name, w.Token)
	runner := Execs[w.Token]
	if runner == nil {
		fmt.Printf("Kill.Go: Token not found\n")
		return
	}
	w.Err = runner.Cmd.Process.Kill()
	if w.Err != nil {
		fmt.Printf("Kill.Go: syscall.Kill returned %s\n", w.Err)
	} else {
		runner.KilledExplicitly = true
	}
	return // implied
}

//Report - report
func (w *Kill) Report(dep int) error {
	ReportNode(dep, w.Name, w.Err)
	return w.Err
}
