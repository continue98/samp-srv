package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {

	if len(os.Args) < 3 && os.Args[1] != "list" {
		fmt.Println("Usage: samp-srv {project_name} stop/start/restart/uptime/list/logs")
		return
	}

	var cmd *exec.Cmd
	var out bytes.Buffer
	var stderr bytes.Buffer

	var project_name string
	var commands string

	if os.Args[1] == "list" {
		commands = os.Args[1]
	} else {
		project_name = os.Args[1]
		commands = os.Args[2]
	}

	s := fmt.Sprintf(`(?m)^\s+(\d+)[.](%s)\s+\(\d+\/\d+\/\d+\s+\d+\:\d+\:\d+\s+\w+\)`, project_name)

	var re1 = regexp.MustCompile(s)

	switch commands {
	case "stop":
		cmd = exec.Command("sh", "-c", "screen -ls")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		match := re1.FindAllStringSubmatch(out.String(), -1)
		if len(match) == 0 {
			fmt.Println("Server not started!")
			return
		}
		for _, val := range match {
			if val[2] == project_name {
				cmd = exec.Command("sh", "-c", "screen -XS "+val[1]+" quit")
				cmd.Run()
			}
		}
		fmt.Println("Server is stop...")
	case "start":
		cmd = exec.Command("sh", "-c", "screen -ls")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		match := re1.FindAllStringSubmatch(out.String(), -1)
		if len(match) != 0 {
			fmt.Println("Server is already running!")
			return
		}
		cmd = exec.Command("sh", "-c", "cd /home/samp_servers/servers/"+project_name+" && screen -L -Logfile /home/samp_servers/servers/"+project_name+" -dmS "+project_name+" ./samp03svr")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		fmt.Println("Server is started...")
	case "restart":
		cmd = exec.Command("sh", "-c", "screen -ls")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		match := re1.FindAllStringSubmatch(out.String(), -1)
		if len(match) == 0 {
			fmt.Println("Server not started!")
			return
		}
		for _, val := range match {
			if val[2] == project_name {
				cmd = exec.Command("sh", "-c", "screen -XS "+val[1]+" quit")
				cmd.Run()
			}
		}
		cmd = exec.Command("sh", "-c", "cd /home/samp_servers/servers/"+project_name+" && screen -L -Logfile /home/samp_servers/servers/"+project_name+" -dmS "+project_name+" ./samp03svr")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()
		fmt.Println("Server is restarted...")
	case "uptime":
		cmd = exec.Command("sh", "-c", "pgrep samp03svr")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()

		var re = regexp.MustCompile(`(?m)^(\d+)$`)
		match := re.FindAllStringSubmatch(out.String(), -1)

		out.Reset()
		stderr.Reset()

		for i := range match {
			cmd = exec.Command("sh", "-c", "pwdx "+match[i][1])
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cmd.Run()
		}

		var re1 = regexp.MustCompile(`(?m)[^\/]*\/(\w+)$`)
		match1 := re1.FindAllStringSubmatch(out.String(), -1)

		var count_server int64
		for i := range match1 {
			out.Reset()
			stderr.Reset()
			if match1[i][1] == project_name {
				cmd = exec.Command("sh", "-c", `ps -o etime= -p "`+match[i][1]+`"`)
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				cmd.Run()

				var re = regexp.MustCompile(`(?m)\s+`)
				uptime_str := re.ReplaceAllString(out.String(), ``)
				fmt.Println("Project: ", project_name, " | PID: ", match[i][1], " | uptime: ", uptime_str)
				count_server++
			}
		}
		if count_server == 0 {
			fmt.Println("Project " + project_name + " not finding process!")
		}
	case "list":
		cmd = exec.Command("sh", "-c", "pgrep samp03svr")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()

		var re = regexp.MustCompile(`(?m)^(\d+)$`)
		match := re.FindAllStringSubmatch(out.String(), -1)

		out.Reset()
		stderr.Reset()

		for i := range match {
			cmd = exec.Command("sh", "-c", "pwdx "+match[i][1])
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cmd.Run()
		}

		var re1 = regexp.MustCompile(`(?m)[^\/]*\/(\w+)$`)
		match1 := re1.FindAllStringSubmatch(out.String(), -1)

		var count_server int64
		for i := range match1 {
			out.Reset()
			stderr.Reset()
			cmd = exec.Command("sh", "-c", `ps -o etime= -p "`+match[i][1]+`"`)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cmd.Run()

			var re = regexp.MustCompile(`(?m)\s+`)
			uptime_str := re.ReplaceAllString(out.String(), ``)
			fmt.Println("Project: ", match1[count_server][1], " | PID: ", match[i][1], " | uptime: ", uptime_str)
			count_server++
		}
	case "logs":
		cmd = exec.Command("sh", "-c", "pgrep samp03svr")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		cmd.Run()

		var re = regexp.MustCompile(`(?m)^(\d+)$`)
		match := re.FindAllStringSubmatch(out.String(), -1)

		out.Reset()
		stderr.Reset()

		for i := range match {
			cmd = exec.Command("sh", "-c", "pwdx "+match[i][1])
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cmd.Run()
		}

		var re1 = regexp.MustCompile(`(?m)\/([A-z0-9-_+]+\/)*([A-z0-9]+)$`)
		match1 := re1.FindAllStringSubmatch(out.String(), -1)

		var count_server int64
		for i := range match1 {
			out.Reset()
			stderr.Reset()
			if match1[i][2] == project_name {

				cmd = exec.Command("sh", "-c", "cat "+match1[i][0]+"/server_log.txt")
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				cmd.Run()
				fmt.Println(out.String())
				count_server++
			}
		}
		if count_server == 0 {
			fmt.Println("Project " + project_name + " not finding process!")
		}
	}
}
