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
}

func setCMD(cli CLI) CLI {
	switch runtime.GOOS {
	case "windows":
		cli.count = "-n"
	default:
		cli.count = "-c"
	}
	fmt.Printf("IP address: ")
	fmt.Scan(&cli.ip)

	return cli
}

func check_loop(cli CLI) {
	file, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("File not created:", err)
		panic("exit")
	}
	fmt.Printf("Start time: %d:%d:%d\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	for {
		out, _ := exec.Command("ping", cli.count, "1", cli.ip).Output()

		if strings.Contains(strings.ToLower(string(out)), "unreachable") {
			if out, _ := exec.Command("ping", cli.count, "1", "216.58.214.142").Output(); strings.Contains(strings.ToLower(string(out)), "unreachable") {
				fmt.Printf("%d:%d:%d - Host network seems down\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
				fmt.Fprintf(file, "%d:%d:%d - Host network seems down\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			}else{
				fmt.Printf("%d:%d:%d - Remote server seems down\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
				fmt.Fprintf(file, "%d:%d:%d - Remote server seems down\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
			}
		} else {
			lines := strings.Split(string(out), "\n")
			line := strings.TrimSpace(lines[1])
			fmt.Printf("%d:%d:%d - %s\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second(), line)
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
