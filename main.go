package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func main() {
	var commands string
	var project_name string
	if len(os.Args) < 3 {
		if len(os.Args) == 2 && os.Args[1] == "list" {
			commands = os.Args[1]
		} else {
			fmt.Println("Usage: samp-srv {project_name} stop/start/restart/uptime/list/logs/install")
			return
		}
	} else {
		project_name = os.Args[1]
		commands = os.Args[2]
	}
	var cmd *exec.Cmd
	var out bytes.Buffer
	var stderr bytes.Buffer

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
	case "install":
		// curl files.sa-mp.com/samp037svr_R2-1.tar.gz --output /tmp/samp037svr_R2-1.tar.gz

		fmt.Println("Start download from files.sa-mp.com...")
		cmd = exec.Command("sh", "-c", "curl files.sa-mp.com/samp037svr_R2-1.tar.gz --output /tmp/samp037svr_R2-1.tar.gz")
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		cmd.Run()
		cmd.Wait()

		fmt.Println("Unpack samp037svr_R2-1.tar.gz to ~/servers/" + project_name + "...")
		cmd = exec.Command("sh", "-c", "tar xzvf /tmp/samp037svr_R2-1.tar.gz -C /tmp")
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		cmd.Run()
		cmd.Wait()

		out.Reset()
		stderr.Reset()

		cmd = exec.Command("sh", "-c", "mv -f /tmp/samp03 /home/samp_servers/servers/"+project_name)
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		cmd.Run()
		out.Reset()
		stderr.Reset()

		fmt.Println("Remove samp037svr_R2-1.tar.gz...")

		cmd = exec.Command("sh", "-c", "rm -rf /tmp/samp037svr_R2-1.tar.gz")
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		cmd.Run()

		out.Reset()
		stderr.Reset()

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Input port server (1000 - 9999): ")

		var port int

		for scanner.Scan() {
			if !IsInt(scanner.Text()) {
				fmt.Println("Input port server (1000 - 9999): ")
				continue
			}
			fmt.Sscan(scanner.Text(), &port)

			if port > 9999 || port < 1000 {
				fmt.Println("Input port server (1000 - 9999): ")
			}
			break
		}

		fmt.Println("Generate rcon password...")

		rcon_password := GenerateRandomString(16)

		fmt.Println("rcon password: ", rcon_password)

		path_to_server_cfg := "/home/samp_servers/servers/" + project_name + "/server.cfg"
		file, err := os.Open(path_to_server_cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)

		cfg_content := string(b)
		var re_cfg_port = regexp.MustCompile(`(?m)(port\s+)(\d+)`)

		port_string := fmt.Sprintf("%d", port)
		cfg_content = re_cfg_port.ReplaceAllString(cfg_content, "${1}"+port_string)

		var re_cfg_rcon_password = regexp.MustCompile(`(?m)(rcon_password\s+)(changeme)`)
		cfg_content = re_cfg_rcon_password.ReplaceAllString(cfg_content, "${1}"+rcon_password)
		ioutil.WriteFile(path_to_server_cfg, []byte(cfg_content), 0644)

		cmd = exec.Command("sh", "-c", "samp-srv "+project_name+" start")
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		cmd.Run()
	default:
		{
			fmt.Println("Invalid command!")
			return
		}

	}

}
func IsInt(s string) bool {
	l := len(s)
	if strings.HasPrefix(s, "-") {
		l = l - 1
		s = s[1:]
	}

	reg := fmt.Sprintf("\\d{%d}", l)

	rs, err := regexp.MatchString(reg, s)

	if err != nil {
		return false
	}

	return rs
}

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenerateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
