package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"example.com/aria2Downloader" // example.com/ have to be replace by relrvant repository path
)

type Pload struct {
	Cmd              string `json:"cmd"`
	Url              string `json:"url"`
	Fs_path          string `json:"fs_path"`
	Output_folder_id string `json:"output_folder_id"`
}
type Task_state struct {
	curr_event_type string // 4 options: start, pause, resume, finish
	progress        string // downloading progress (option)
	state           string // 4 options: started, paused, resumed, finished
	gid             string // after first addDl, aria2 retuen gid (the download identifire)
	Payload         Pload
}

var log_test_file = "data_default"
var Wg sync.WaitGroup
var tasks_state = make(map[string]Task_state) // key of map is task_id
var curr_task_state Task_state
var curr_task_id string
var curr_event_type string

// start download handler and return download gid
func start_handler(urisList []string) (ret_gid string) {
	fmt.Println("in start handler") // for debug
	ret_gid = aria2Downloader.Download(urisList)

	return ret_gid
}

// pause download by its gid
func pause_handler(gid string) (paused_gid string) {
	fmt.Println("in pause handler") // for debug
	paused_gid = aria2Downloader.Pause_dl(gid)
	return paused_gid
}

// unpause download by its gid
func resume_handler(gid string) (resumed_gid string) {
	fmt.Println("in resume handler") // for debug
	resumed_gid = aria2Downloader.Unpause_dl(gid)
	return resumed_gid
}

// remove download by its gid
func remove_handler(gid string) (removed_gid string) {
	fmt.Println("in remove handler") // for debug
	removed_gid = aria2Downloader.Remove_dl(gid)
	return removed_gid
}

func commit_test() {

	// wait 2 seconds and return true with no errors
}

func commit_seedrfs() {

}

func commit_rclone() {

}

//to do -
func start_aria2_server() {
	fmt.Println("in start aria2 server function")
	defer Wg.Done()
	cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-port=6800")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func start_http() {

	// net.http
}

// start_queue():
// read log file for test
// got adduris start (add uri to download) event ->
// 1. add to tasks table
// 2. add payload to payloads dictionary
// 3. call adduris start (add uri to download) command
func start_queue() {

	defer Wg.Done()

	//var line string
	var line_tokens []string
	var curr_pload Pload
	//	var out io.Writer

	// read logfile
	f, err := os.Open(log_test_file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	oldSize := info.Size()
	for {

		for line, prefix, err := r.ReadLine(); err != io.EOF; line, prefix, err = r.ReadLine() {
			if prefix {
				//	fmt.Fprint(out, string(line)) //for debug

				// call parse
				line_tokens, curr_pload = parse_line(string(line))
				// call general handler
				general_handler(line_tokens, curr_pload)

			} else {
				//	fmt.Fprintln(out, string(line)) //for debug

				// call parse
				line_tokens, curr_pload = parse_line(string(line))
				fmt.Printf("in reading: line_tokens= %s", line_tokens)
				// call general handler
				general_handler(line_tokens, curr_pload)
			}
		}

		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(err)
		}

		for {
			time.Sleep(time.Second)
			newinfo, err := f.Stat()
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize {
					f.Seek(0, 0)
				} else {
					f.Seek(pos, io.SeekStart)
				}
				r = bufio.NewReader(f)
				oldSize = newSize
				break
			}
		}
	}
}

//
func start_commit() {

	// calls and waits for external HTTP service with HTTP client
}

func Start() {
	Wg.Add(1) //temp
	go start_aria2_server()
	time.Sleep(2 * time.Second)
	Wg.Add(1)
	go start_queue()

	//	defer Wg.Done()
	//	Wg.Wait()
	//	go..
	//	go..

}

// parse line from logfile to commands and payload
func parse_line(line string) (line_tokens []string, pload Pload) {

	var (
		str_payload string
		curr_pload  Pload
	)

	//
	line_str := strings.TrimSpace(string(line))
	line_tokens = strings.Split(line_str, "<->")
	str_payload = line_tokens[3]

	fmt.Printf("str_payload= %s\n", str_payload) // for debug

	var jsonBlob = []byte(str_payload)
	err := json.Unmarshal(jsonBlob, &curr_pload)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("curr pload= %s\n", curr_pload) // for debug

	return line_tokens, curr_pload
}

func general_handler(line_tokens []string, curr_pload Pload) {
	var uris_list []string
	var ret_gid string
	uris_list = make([]string, 1, 1)
	curr_task_id = line_tokens[0]
	//curr_event_id = line_tokens[1]
	curr_event_type = line_tokens[2]
	cmd := curr_pload.Cmd
	url := curr_pload.Url
	//	fs_path := curr_pload.Fs_path
	//	output_folder_id := curr_pload.Output_folder_id

	if cmd == "download" && curr_event_type == "start" {
		fmt.Println("in start queue before handler")

		//	curr_event := events[eve_id]
		uris_list[0] = url
		ret_gid = start_handler(uris_list)
		fmt.Println("in start queue after handler")
		//update ret_gid state
		curr_task_state.gid = ret_gid

		// update other Task_state fields
		curr_task_state.curr_event_type = curr_event_type //curr_event_type=="start"
		curr_task_state.progress = ""
		curr_task_state.state = "started"
		curr_task_state.Payload.Cmd = curr_pload.Cmd
		curr_task_state.Payload.Url = curr_pload.Url
		curr_task_state.Payload.Fs_path = curr_pload.Fs_path
		curr_task_state.Payload.Output_folder_id = curr_pload.Output_folder_id

		tasks_state[curr_task_id] = curr_task_state

		fmt.Printf("tasks_state[curr_task_id]= %s\n", tasks_state[curr_task_id]) //for debug
		fmt.Printf("curr_task_state.gid= %s\n", curr_task_state.gid)             //for debug

	}
	if cmd == "download" && curr_event_type == "pause" {
		//update curr_task_state
		curr_task_state.curr_event_type = curr_event_type //curr_event_type=="pause"
		// read gid of the task (download) for doing pause
		curr_task_state = tasks_state[curr_task_id]
		gid := curr_task_state.gid
		//call pause handler
		paused_gid := pause_handler(gid)
		// update state after pause
		if paused_gid == gid {
			curr_task_state.state = "paused"
		}
	}

	if cmd == "download" && curr_event_type == "resume" {
		//update curr_task_state
		curr_task_state.curr_event_type = curr_event_type //curr_event_type=="resume"
		// read gid of the task (download) for doing pause
		curr_task_state = tasks_state[curr_task_id]
		gid := curr_task_state.gid
		//call resume handler
		resumed_gid := resume_handler(gid)
		// update state after pause
		if resumed_gid == gid {
			curr_task_state.state = "resumed"
		}

	}
	if cmd == "download" && curr_event_type == "finish" {
		//update curr_task_state
		curr_task_state.curr_event_type = curr_event_type // curr_event_type=="finish"
		// read gid of the task (download) for doing pause
		curr_task_state = tasks_state[curr_task_id]
		gid := curr_task_state.gid
		//call resume handler
		removed_gid := remove_handler(gid)
		// update state after pause
		if removed_gid == gid {
			curr_task_state.state = "removed"
			//delete the download from map
			delete(tasks_state, curr_task_id)
		}

	}

}
