// SSH_AUTH_SOCK=/tmp/ssh-UqjBf7VqtHwH/agent.18157; export SSH_AUTH_SOCK;
// SSH_AGENT_PID=18158; export SSH_AGENT_PID;

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
)

const PID_ENV = "SSH_AGENT_PID"
const SOCK_ENV = "SSH_AUTH_SOCK"
const EXE = "ssh-agent"

var verbose = flag.Bool("v", false, "Enable verbose output")

func exitErr(format string, args ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintf(os.Stderr, "\n")
	}
	os.Exit(1)
}

func check_pid(pid int) {
	commPath := path.Join("/", "proc", strconv.Itoa(pid), "comm")
	f, err := os.Open(commPath)
	if err != nil {
		exitErr("No such PID: %d", pid)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan()
	if s.Err() != nil {
		exitErr("Error reading from comm!: %v", err)
	}

	comm := s.Text()
	if comm != EXE {
		exitErr("PID %d is %s, expected %s", pid, comm, EXE)
	}
}

func check_sock(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		exitErr("Error checking socket %s: %v", path, err)
	}

	if fi.Mode()&os.ModeSocket == 0 {
		exitErr("%s is not a socket", path)
	}
}

func main() {
	flag.Parse()

	pidStr, env_set := os.LookupEnv(PID_ENV)
	if !env_set {
		exitErr("%s is not set", PID_ENV)
	}
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		exitErr("%s is not a valid PID", pidStr)
	}

	sock, env_set := os.LookupEnv(SOCK_ENV)
	if !env_set {
		exitErr("%s is not set", SOCK_ENV)
	}

	check_pid(pid)
	check_sock(sock)
}
