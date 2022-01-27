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

	// example.com/ have to be replace by relrvant repository path
	"example.com/aria2Downloader/downloader"
)

type Pload struct {
	Cmd            string `json:"cmd"`
	Url            string `json:"url"`
	FsPath         string `json:"fsPath"`
	OutputFolderId string `json:"outputFolderId"`
}
type TaskState struct {
	currEventType string // 5 options: start, pause, resume, status, finish
	progress      string // downloading progress (option)
	state         string // 4 options: started, paused, resumed, finished
	gid           string // after first Download, aria2 retuen gid (the download identifire)
	Payload       Pload
}

var logTestFile = "data_default"
var Wg sync.WaitGroup
var tasksState = make(map[string]TaskState) // key of map is task id
var currTaskState TaskState
var currTaskId string
var currEventType string

func Start() {

	Wg.Add(1)

	go startAria2Server()

	time.Sleep(2 * time.Second)

	Wg.Add(1)

	go startQueue()
}

func startAria2Server() {

	defer Wg.Done()

	cmd := exec.Command("aria2c", "--enable-rpc", "--rpc-listen-port=6800")

	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}
}

// startQueue function:
// Read tasks from log file
// Update current Payload (currPload) with task details
// Call general handler to continue handling deferent tasks
func startQueue() {

	defer Wg.Done()

	var lineTokens []string
	var currPload Pload

	// read logfile
	f, err := os.Open(logTestFile)
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
				lineTokens, currPload = parseLine(string(line))
				// call general handler
				generalHandler(lineTokens, currPload)

			} else {
				//	fmt.Fprintln(out, string(line)) //for debug

				// call parse
				lineTokens, currPload = parseLine(string(line))
				fmt.Printf("in reading: line_tokens= %s", lineTokens)
				// call general handler
				generalHandler(lineTokens, currPload)
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

// General handling for deferent tasks

func generalHandler(lineTokens []string, currPload Pload) {

	var urisList []string
	var retGid string

	urisList = make([]string, 1, 1)
	currTaskId = lineTokens[0]

	currEventType = lineTokens[2]
	cmd := currPload.Cmd
	url := currPload.Url

	if cmd == "download" && currEventType == "start" {

		urisList[0] = url
		retGid = startHandler(urisList)

		//update ret_gid state
		currTaskState.gid = retGid

		// update other Task_state fields
		currTaskState.currEventType = currEventType //currEventType=="start"
		currTaskState.progress = ""
		currTaskState.state = "started"
		currTaskState.Payload.Cmd = currPload.Cmd
		currTaskState.Payload.Url = currPload.Url
		currTaskState.Payload.FsPath = currPload.FsPath
		currTaskState.Payload.OutputFolderId = currPload.OutputFolderId

		tasksState[currTaskId] = currTaskState

	}

	if cmd == "download" && currEventType == "pause" {

		//update curr_task_state
		currTaskState.currEventType = currEventType //currEventType=="pause"

		// read gid of the task (download) for doing pause
		currTaskState = tasksState[currTaskId]
		gid := currTaskState.gid

		//call pause handler
		pausedGid := pauseHandler(gid)

		// update state after pause
		if pausedGid == gid {
			currTaskState.state = "paused"
		}
	}

	if cmd == "download" && currEventType == "resume" {
		//update curr_task_state
		currTaskState.currEventType = currEventType //currEventType=="resume"
		// read gid of the task (download) for doing resume
		currTaskState = tasksState[currTaskId]
		gid := currTaskState.gid
		//call resume handler
		resumedGid := resumeHandler(gid)
		// update state after pause
		if resumedGid == gid {
			currTaskState.state = "resumed"
		}

	}
	/////////////////////////////////////////////////
	//Todo - test status event and statusHandler()///
	/////////////////////////////////////////////////
	if cmd == "download" && currEventType == "status" {
		//update curr_task_state
		currTaskState.currEventType = currEventType //currEventType=="status"
		// read gid of the task (download) for status
		currTaskState = tasksState[currTaskId]
		gid := currTaskState.gid
		//call status handler
		retStatus := statusHandler(gid)
		// update state after status delivered
		if retStatus.Result.Gid == gid {
			currTaskState.state = "statusDone"
		}

	}

	if cmd == "download" && currEventType == "finish" {
		//update curr_task_state
		currTaskState.currEventType = currEventType // currEventType=="finish"
		// read gid of the task (download) for doing pause
		currTaskState = tasksState[currTaskId]
		gid := currTaskState.gid
		//call resume handler
		removedGid := removeHandler(gid)
		// update state after pause
		if removedGid == gid {
			currTaskState.state = "removed"
			//delete the download from map
			delete(tasksState, currTaskId)
		}

	}

}

// Specific handlers for downloader functions : Start , Pause, Resume , Remove

// start download handler and return download gid
func startHandler(urisList []string) (retGid string) {

	retGid = downloader.Start(urisList)

	return retGid
}

// pause download by its gid
func pauseHandler(gid string) (pausedGid string) {

	pausedGid = downloader.Pause(gid)
	return pausedGid
}

// resume download by its gid
func resumeHandler(gid string) (resumedGid string) {

	resumedGid = downloader.Resume(gid)
	return resumedGid
}

// remove download by its gid
func removeHandler(gid string) (removedGid string) {

	removedGid = downloader.Remove(gid)
	return removedGid
}

// remove download by its gid
func statusHandler(gid string) (retStatus downloader.RetStatusMsg) {

	retStatus = downloader.DownloadStatus(gid)
	return retStatus
}

// helper functions

// parse line from logfile to commands and payload
func parseLine(line string) (lineTokens []string, pload Pload) {

	var strPload string
	var currPload Pload

	line_str := strings.TrimSpace(string(line))
	lineTokens = strings.Split(line_str, "<->")
	strPload = lineTokens[3]

	var jsonBlob = []byte(strPload)

	err := json.Unmarshal(jsonBlob, &currPload)

	if err != nil {

		fmt.Println("error:", err)
	}

	return lineTokens, currPload
}

// To do - some functions to complete for future use

func commitTest() {

	// wait 2 seconds and return true with no errors
}

func commitSeedrFs() {

}

func commitRclone() {

}

func startHttp() {

	// net.http
}

//
func startCommit() {

	// calls and waits for external HTTP service with HTTP client
}
