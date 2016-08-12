package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/chzyer/readline"
)

func checkGo() error {
	cmd := exec.Command("go")
	if err := cmd.Start(); err != nil {
		return err
	}
	// it's exit code should return 2
	return nil
}

func runPprof(args []string) error {
	newArgs := make([]string, 2+len(args))
	copy(newArgs[2:], args)
	newArgs[0] = "tool"
	newArgs[1] = "pprof"
	cmd := exec.Command("go", newArgs...)

	pr, pw := io.Pipe()
	cmd.Stderr = pw
	cmd.Stdout = pw
	pprofIn, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_ = pprofIn

	if err := cmd.Start(); err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		errChan <- cmd.Wait()
	}()
	go process(pr, pprofIn)

	err = <-errChan
	pr.CloseWithError(err)
	return err
}

func process(pr io.ReadCloser, pprofIn io.Writer) {
	r := bufio.NewReader(pr)
	defer pr.Close()

	rl := newrl()
	defer rl.Close()
	buffer := bytes.NewBuffer(nil)

	for {
		line, err := lineOrEOF(r, buffer)
		if err != nil {
			break
		}
		if line != "(pprof) " {
			os.Stdout.Write([]byte(line))
			continue
		}

	reread:
		ret := rl.Line()
		if ret.CanBreak() {
			pprofIn.Write([]byte("quit\n"))
			break
		} else if ret.CanContinue() {
			goto reread
		}

		pprofIn.Write([]byte(ret.Line))
		pprofIn.Write([]byte("\n"))
	}

}

func lineOrEOF(r *bufio.Reader, buf *bytes.Buffer) (string, error) {
	buf.Reset()
	for {
		c, err := r.ReadByte()
		if err != nil {
			break
		}
		buf.WriteByte(c)
		if c == '\n' {
			break
		}
		if bytes.Equal(buf.Bytes(), []byte("(pprof) ")) {
			break
		}
	}
	return buf.String(), nil
}

func newrl() *readline.Instance {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "(gopprof) ",
	})
	if err != nil {
		panic(err)
	}
	return rl
}

func main() {
	if err := checkGo(); err != nil {
		println(err.Error())
		os.Exit(2)
	}

	if err := runPprof(os.Args[1:]); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
