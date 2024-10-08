package runner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
)

// get input from insignals ...
// parse if own command
// parse if active key combi
// send to running program

func RunProg(inKeys chan int, outData chan byte) (err error) { //add progname and args
	fmt.Println("runner start")
	progName := "bash" // hardcoded for now
	// progArgs := []string{""}

	cmd := exec.Command(progName)
	progOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("ERROR Runner: ", err)
		return err
	}
	progIn, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("ERROR Runner: ", err)
		return err
	}

	//handle input

	//handle output
	fmt.Println("runner start handle output")
	go handleProgOutput(cmd, progOut, outData)

	err = cmd.Start()
	fmt.Println("prog started ", cmd.ProcessState)
	if err != nil {
		fmt.Println("ERROR Runner: ", err)
	}
	// runner loop
	for {
		select {
		case key := <-inKeys:
			fmt.Println("runner got keys")
			progIn.Write([]byte{byte(key)})
		}
	}

}

func handleProgOutput(cmd *exec.Cmd, progOut io.ReadCloser, outData chan byte) {
	buf := bufio.NewReader(progOut)
	outData <- 'A'
	for {
		// line, _, err := buf.ReadLine()
		bits, err := buf.ReadByte()
		if err != nil {
			// fmt.Println("ERROR:-> Runner: buf got EOF")
		} else {
			fmt.Println("got some data")
			outData <- bits
		}

	}

}

// program running here ?
// maiybe put in runner mod but furst draft here

func CheckRunning(cmd *exec.Cmd) bool {
	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() || cmd.Process == nil {
		return false
	}
	return true
}

func RunProg2() {

	me, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println(">>> user is: ", me)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(">>> cwd is: ", cwd)

	shelly := os.Getenv("SHELL")
	fmt.Println(">>> shelly is: ", shelly)

	cmd := exec.Command(shelly)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(">>> error setting out pipe: ", err)
		return
	}

	data := make([]byte, 1024)
	bufData := bufio.NewReader(output)

	go func() {
		for {
			_, err := bufData.Read(data)
			if err != io.EOF {
				//do something with output
				// fmt.Println(">>> got soemthing") //TEST:
				fmt.Print(string(data))
				data = make([]byte, 1024)
			}
		}
	}()

	// os.Stderr.Close()
	cmd.Start()
	retval := cmd.Wait()

	fmt.Println(">>> shell ended: ", retval)
}

// REM: ========================================================================
func RunProg3(ended chan<- bool, outputChan chan<- byte, errChan chan<- byte, inputChan <-chan int) {

	me, err := user.Current()
	if err != nil {
		panic(err)
	}
	log.Println(">>> user is: ", me)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(">>> cwd is: ", cwd)

	shelly := os.Getenv("SHELL")
	log.Println(">>> shelly is: ", shelly)

	progArgs := []string{"-i"}
	cmd := exec.Command(shelly, progArgs...)

	cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stdout
	// cmd.Stdout = os.Stdout

	// os.Stderr.Close()

	progErr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(">>> error setting err pipe: ", err)
		return
	}

	progOutput, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(">>> error setting out pipe: ", err)
		return
	}

	outputData := make([]byte, 128)
	bufProgOutput := bufio.NewReader(progOutput)

	//handle stdout
	go func() {
		for {
			_, err := bufProgOutput.Read(outputData)
			if err != io.EOF {
				for _, b := range outputData {
					if b != 0 {
						outputChan <- b
					}
				}
				outputData = make([]byte, 128)
			}
		}
	}()

	errData := make([]byte, 128)
	bufProgErr := bufio.NewReader(progErr)

	//handle stderr
	go func() {
		for {
			_, err := bufProgErr.Read(errData)
			if err != io.EOF {
				for _, b := range errData {
					if b != 0 {
						errChan <- b
					}
				}
				errData = make([]byte, 128)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println(">>> error starting prog: ", err)
	}

	retval := cmd.Wait()

	ended <- true

	log.Println(">>> shell ended: ", retval)
}
