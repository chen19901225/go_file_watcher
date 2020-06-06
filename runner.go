package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/radovskyb/watcher"
	"golang.org/x/sys/unix"
)

func runCommand(command string, directory string, instanceWatcher *watcher.Watcher) {
	fmt.Printf("run command:%s\n", command)
	var cmdStart = []string{"/bin/sh", "-c"}
	var procAttrs = &unix.SysProcAttr{Setpgid: true}
	cs := append(cmdStart, command)
	// 这个...
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = directory
	// 这个东西的作用是什么？
	cmd.SysProcAttr = procAttrs
	if err := cmd.Start(); err != nil {
		select {
		case instanceWatcher.Error <- err:
		default:
		}
		fmt.Printf("Failed to start %s: %s\n", command, err)
		return
	}
}

func handleEvent(event watcher.Event, commandList []interface{},
	directory string,
	instanceWatcher *watcher.Watcher) {
	path := event.Path
	for _, orgPiece := range commandList {
		piece := orgPiece.(map[string]interface{})
		command, ok := piece["command"]
		if ok == false {
			log.Fatal(fmt.Sprintf("no command in %v", piece))
		}
		subDirectory, ok := piece["directory"]
		if ok == false {
			subDirectory = directory
		} else {
			subDirectory = subDirectory.(string)
		}
		pattern, ok := piece["pattern"]
		if ok == false {
			runCommand(command.(string), subDirectory.(string), instanceWatcher)
			continue
		}

		regPattern := regexp.MustCompile(pattern.(string))
		if regPattern.MatchString(path) {
			runCommand(command.(string), subDirectory.(string), instanceWatcher)
			continue
		}

	}

}
