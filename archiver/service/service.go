package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.com/archiver/archiverLib/rarop" // example.com/ have to be replace by relrvant repository path
	"example.com/archiver/archiverLib/tarop" // example.com/ have to be replace by relrvant repository path
	"example.com/archiver/archiverLib/zipop" // example.com/ have to be replace by relrvant repository path
)

type Pload struct {
	Cmd            string `json:"cmd"`
	Path           string `json:"path"`
	Url            string `json:"url"`
	ArchiveType    string `json:"archiveType"`
	FsPath         string `json:"fsPath"`
	OutputFolderId string `json:"outputFolderId"`
	ArchiveFileIds string `json:"archiveFileIds"`
}
type Task_state struct {
	currEventType string // 4 options: start, pause, resume, finish
	progress      string // downloading progress (option)
	state         string // 4 options: started, paused, resumed, finished
	zipMeta       []zipop.FileMetadata
	tarMeta       []tarop.FileMetadata
	Payload       Pload
}

var logTestFile = "data_default"
var Wg sync.WaitGroup
var tasksState = make(map[string]Task_state) // key is task_id
var currTaskState Task_state
var currTaskId string
var currEventType string

func Start() {

	Wg.Add(1)
	go startQueue()
}

// startQueue function:
// Read tasks from log file
// Update current Payload (currPload) with task details
// Call general handler to continue handling deferent tasks
func startQueue() {
	fmt.Println("in start q") //for debug
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
				fmt.Printf("in reading: line_tokens= %s\n", lineTokens)
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

// General handling for deferent tasks.
func generalHandler(lineTokens []string, currPload Pload) {
	var pathOrUrl string
	var isUrl bool

	isUrl = false //temp
	currTaskId = lineTokens[0]
	//	curr_event_id = line_tokens[1]
	currEventType = lineTokens[2]
	cmd := currPload.Cmd
	url := currPload.Url
	//	fs_path := curr_pload.Fs_path
	output_folder_id := currPload.OutputFolderId
	// to archive cmd - file paths seperated with ',' . to extract cmd - file ids seperated with ',' .
	// to extract rar - file names seperated with ","
	files := currPload.ArchiveFileIds
	files_to_rar := files
	path := currPload.Path
	if url != "" {
		pathOrUrl = url
		isUrl = true
	} else {
		pathOrUrl = path
	}

	// update state before specific handler
	currTaskState.currEventType = currEventType //curr_event_type=="start"
	currTaskState.progress = ""
	currTaskState.Payload.Cmd = currPload.Cmd
	currTaskState.Payload.Url = currPload.Url
	currTaskState.Payload.FsPath = currPload.FsPath
	currTaskState.Payload.OutputFolderId = currPload.OutputFolderId

	tasksState[currTaskId] = currTaskState

	if cmd == "archive" && currEventType == "start" {

		fmt.Println("in start queue before handler") //for debug

		//call archive handler
		err := createArchiveHandler(files, path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("in start queue after handler") //for debug

		//update state after archive
		currTaskState.state = "started"
		tasksState[currTaskId] = currTaskState

	}

	if cmd == "unarchive" && currEventType == "start" {

		//call unarchive handler with local path or url acording to url_flag
		err := extractHandler(pathOrUrl, isUrl, output_folder_id)
		if err != nil {
			log.Fatal(err)
		}

		//to do - update state after unarchive
		currTaskState.state = "started"
		tasksState[currTaskId] = currTaskState

	}

	if cmd == "extract" && currEventType == "start" {

		//call archive handler
		err := extractSomeHandler(pathOrUrl, isUrl, files, files_to_rar, output_folder_id)
		if err != nil {
			log.Fatal(err)
		}
		// update state after extract

		currTaskState.state = "started"
		tasksState[currTaskId] = currTaskState
	}

	if cmd == "meta" && currEventType == "start" {

		//call archive handler with local path or url acording to url_flag
		err := metaHandler(pathOrUrl, isUrl)
		if err != nil {
			log.Fatal(err)
		}
		// update state after meta
		currTaskState.state = "started"
		tasksState[currTaskId] = currTaskState

	}

	if currEventType == "finish" {
		//delete curr_task_state
		delete(tasksState, currTaskId)
	}

}

// Specific handlers for archive functions : Create , Extract, ExtractSome , Meta

// Handler for Create archive functions
func createArchiveHandler(filesStr, destArchivePath string) error {

	var files []string
	// get archive name from src string
	files = strings.Split(filesStr, ",")

	// check the type of archive (.zip , .tar , tar.gz , .tgz) and call archive function
	if strings.HasSuffix(destArchivePath, ".zip") {
		err := zipop.Create(files, destArchivePath)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasSuffix(destArchivePath, ".tar") {
		err := tarop.CreateTar(files, destArchivePath)
		if err != nil {
			log.Fatal(err)
		}
	} else if (strings.HasSuffix(destArchivePath, ".tar.gz")) || (strings.HasSuffix(destArchivePath, ".tgz")) {
		err := tarop.CreateTargz(files, destArchivePath)
		if err != nil {
			log.Fatal(err)
		}

	} else if strings.HasSuffix(destArchivePath, ".rar") {
		//Rar(srcRarArchive, files string)
		rarop.Create(destArchivePath, filesStr)
	}
	fmt.Println("The new archive stored in: ", destArchivePath)

	return nil

}

// Handler for Extract functions
func extractHandler(srcArchive string, isUrl bool, destDirPath string) error {

	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
	if strings.HasSuffix(srcArchive, ".zip") {

		_, err := zipop.Extract(srcArchive, isUrl, destDirPath)
		if err != nil {
			log.Fatal(err)
		}

	} else if (strings.HasSuffix(srcArchive, ".tar")) || (strings.HasSuffix(srcArchive, ".tar.gz")) || (strings.HasSuffix(srcArchive, ".tgz")) {

		tarGzArchiveName, err := tarop.Extract(srcArchive, isUrl, destDirPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Archive name: %s\n", tarGzArchiveName)

	} else if strings.HasSuffix(srcArchive, ".rar") {

		rarop.Extract(srcArchive, isUrl, destDirPath)

	}
	fmt.Printf("unarchive %s stored in: %s\n", srcArchive, destDirPath)

	return nil
}

// Handler for ExtractSome functions
func extractSomeHandler(srcArchive string, isUrl bool, targetIdsFiles string, targetStrFiles string, destDirPath string) error {

	var fileIdsStr []string
	var filesNames []string
	var fileIds []int

	filesNames = make([]string, 0)
	fileIds = make([]int, 0)

	// get archive name from src string
	fileIdsStr = strings.Split(targetIdsFiles, ",")

	//convert file ids to file paths
	for _, strId := range fileIdsStr {
		fileId, err := strconv.Atoi(strId)
		if err != nil {
			log.Fatal(err)
		}
		fileIds = append(fileIds, fileId)
	}

	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
	if strings.HasSuffix(srcArchive, ".zip") {

		for _, id := range fileIds {
			currFileName := tasksState[currTaskId].zipMeta[id].FileName
			filesNames = append(filesNames, currFileName)
		}

		err := zipop.ExtractSome(srcArchive, isUrl, filesNames, destDirPath)
		if err != nil {
			log.Fatal(err)
		}
	} else if (strings.HasSuffix(srcArchive, ".tar")) || (strings.HasSuffix(srcArchive, ".tar.gz")) || (strings.HasSuffix(srcArchive, ".tgz")) {

		for _, id := range fileIds {
			currFileName := tasksState[currTaskId].tarMeta[id].FileName
			filesNames = append(filesNames, currFileName)
		}

		_, err := tarop.ExtractSome(srcArchive, isUrl, filesNames, destDirPath)
		if err != nil {
			log.Fatal(err)
		}

	} else if strings.HasSuffix(srcArchive, ".rar") {
		// change  flie names string to be seperated by space instead of ','
		fileStrNames := strings.ReplaceAll(targetStrFiles, ",", " ")

		rarop.ExtractSome(srcArchive, isUrl, fileStrNames, destDirPath)

	}
	fmt.Printf("extracted files stored in: %s\n", destDirPath)

	return nil

}

// Handler for Meta functions
func metaHandler(srcArchive string, isUrl bool) error {

	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call meta function
	if strings.HasSuffix(srcArchive, ".zip") {

		currTaskZipMeta, err := zipop.Meta(srcArchive, isUrl)
		if err != nil {
			log.Fatal(err)
		}

		arrLen := len(currTaskZipMeta)

		//update state
		currTaskState.zipMeta = make([]zipop.FileMetadata, arrLen)
		currTaskState.zipMeta = currTaskZipMeta
		tasksState[currTaskId] = currTaskState

	} else if (strings.HasSuffix(srcArchive, ".tar")) || (strings.HasSuffix(srcArchive, ".tar.gz")) || (strings.HasSuffix(srcArchive, ".tgz")) {

		currTaskTarMeta, err := tarop.Meta(srcArchive, isUrl)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("curr_task_targz_metadata= %s\n", currTaskTarMeta) // temp  for debug
		arrLen := len(currTaskTarMeta)

		//update state
		currTaskState.tarMeta = make([]tarop.FileMetadata, arrLen)
		currTaskState.tarMeta = currTaskTarMeta

		tasksState[currTaskId] = currTaskState

	} else if strings.HasSuffix(srcArchive, ".rar") {

		rarop.Meta(srcArchive, isUrl)

		// to do - yet not implemented  update state by meta (only print to console)

	}

	return nil
}

// helper functions

// parse line from logfile to commands and payload
func parseLine(line string) (lineTokens []string, pload Pload) {

	var (
		strPayload string
		currPload  Pload
	)

	lineStr := strings.TrimSpace(string(line))
	lineTokens = strings.Split(lineStr, "<->")
	strPayload = lineTokens[3]

	var jsonBlob = []byte(strPayload)
	err := json.Unmarshal(jsonBlob, &currPload)
	if err != nil {
		fmt.Println("error:", err)
	}

	return lineTokens, currPload
}

// To do - some functions to complete for future use

//
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
