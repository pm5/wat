package main

import "github.com/domluna/watcher"
import "log"
import "os"
import "os/exec"
import "strings"

var cmds []*exec.Cmd
var files []string

func readArgs() {
	i := 1
	if len(os.Args) < 3 {
		return
	}
	for i < len(os.Args)-1 && os.Args[i] != "--" {
		terms := strings.Split(os.Args[i], " ")
		cmd := exec.Command(terms[0], terms[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmds = append(cmds, cmd)
		i++
	}
	if os.Args[i] == "--" {
		i++
	}
	for i < len(os.Args) {
		files = append(files, os.Args[i])
		i++
	}
}

func runCmds() {
	var err error
	for i, cmd := range cmds {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			cmds[i] = exec.Command(cmd.Path, cmd.Args[1:]...)
			cmd = cmds[i]
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func watchFiles() error {
	done := make(chan bool)
	for _, f := range files {
		w, err := watcher.New(f)
		defer w.Close()
		if err != nil {
			log.Fatal(err)
		}
		w.Watch()

		go func() {
			for {
				select {
				case <-w.Events:
					runCmds()
				}
			}
		}()
	}
	<-done
	return nil
}

func main() {
	readArgs()
	if len(cmds) == 0 || len(files) == 0 {
		os.Exit(-1)
	}
	runCmds()
	watchFiles()
}
