package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type CLI struct {
	count string
	ip    string
	pingL int
}

// "setCMD" function sets parameters according to the operating system on which it is run.
// and gets the target remote server ip address from the user.
func setCMD(cli CLI) CLI {
	switch runtime.GOOS {
	case "windows":
		cli.count = "-n"
		cli.pingL = 2
	default:
		cli.count = "-c"
		cli.pingL = 1
	}
	fmt.Printf("IP address: ")
	fmt.Scan(&cli.ip)

	return cli
}

// "check_loop" sends ICMP request to target server at 30 second intervals.
// If it cannot reach the target server, it sends an ICMP request to 216.58.214.142 "www.google.com" to check whether the client is connected to the internet.
// In case the client disconnects and the remote server does not respond to ICMP requests, it saves the date and response in the "log.txt" file.
func check_loop(cli CLI) {
	file, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("File not created:", err)
		panic("exit")
	}
	fmt.Printf("Start time: D %d:M %d - %d:%d:%d\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	for {
		out, _ := exec.Command("ping", cli.count, "1", cli.ip).Output()

		if strings.Contains(strings.ToLower(string(out)), "unreachable") {
			if out, _ := exec.Command("ping", cli.count, "1", "216.58.214.142").Output(); strings.Contains(strings.ToLower(string(out)), "unreachable") {
				fmt.Printf("D %d:M %d - %d:%d:%d - Host network seems down\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
				fmt.Fprintf(file, "D %d:M %d - %d:%d:%d - Host network seems down\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			} else {
				fmt.Printf("D %d:M %d - %d:%d:%d - Remote server seems down\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
				fmt.Fprintf(file, "D %d:M %d - %d:%d:%d - Remote server seems down\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			}
		} else if strings.Contains(strings.ToLower(string(out)), "timed out") {
			fmt.Printf("D %d:M %d - %d:%d:%d - Remote server timed out\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

			fmt.Fprintf(file, "D %d:M %d - %d:%d:%d - Remote server timed out\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		} else {
			lines := strings.Split(string(out), "\n")
			line := strings.TrimSpace(lines[cli.pingL])
			fmt.Printf("D %d:M %d - %d:%d:%d - %s\n", time.Now().Day(), time.Now().Month(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), line)
		}

		time.Sleep(30 * time.Second)
	}
}

func main() {
	var cli CLI
	cli = CLI{}

	cli = setCMD(cli)
	out, _ := exec.Command("ping", cli.count, "1", cli.ip).Output()

	if strings.Contains(strings.ToLower(string(out)), "unreachable") {
		if out, _ := exec.Command("ping", cli.count, "1", "216.58.214.142").Output(); strings.Contains(strings.ToLower(string(out)), "unreachable") {
			panic("You are not connected to the internet.")
		}
	} else {
		check_loop(cli)
	}
}
