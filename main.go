package main

import "github.com/go-fsnotify/fsnotify"
import "log"
import "os"
import "os/exec"
import "strings"

var cmds []*exec.Cmd
var files []string

func readArgs() {
	i := 1
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
		}
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func watchFiles() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			log.Println(".")
			select {
			case _ = <-watcher.Events:
				runCmds()
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	for _, f := range files {
		err = watcher.Add(f)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
	return nil
}

func main() {
	readArgs()
	watchFiles()
}
