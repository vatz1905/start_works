package archiver_service

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

	// example.com/ have to be replace by relrvant repository path
	"example.com/archiver_lib/rarop" // example.com/ have to be replace by relrvant repository path
	"example.com/archiver_lib/tarop" // example.com/ have to be replace by relrvant repository path
	"example.com/archiver_lib/zipop" // example.com/ have to be replace by relrvant repository path
)

type Pload struct {
	Cmd              string `json:"cmd"`
	Path             string `json:"path"`
	Url              string `json:"url"`
	Archive_type     string `json:"archive_type"`
	Fs_path          string `json:"fs_path"`
	Output_folder_id string `json:"output_folder_id"`
	Archive_file_ids string `json:"archive_file_ids"`
}
type Task_state struct {
	curr_event_type string // 4 options: start, pause, resume, finish
	progress        string // downloading progress (option)
	state           string // 4 options: started, paused, resumed, finished
	//	archive_metadata []zipop.FileMetadata //[]interface{} //temp
	zip_meta []zipop.FileMetadata //[]interface{} //temp
	tar_meta []tarop.FileMetadata //[]interface{} //temp
	Payload  Pload
}

var log_test_file = "data_default"
var Wg sync.WaitGroup
var tasks_state = make(map[string]Task_state) // key is task_id
var curr_task_state Task_state
var curr_task_id string
var curr_event_type string
var curr_task_zip_meta []zipop.FileMetadata //temp
var curr_task_tar_meta []tarop.FileMetadata //temp

func archive_handler(files_str, dest_archive_path string) error {

	fmt.Println("archive_handler") // for debug
	var files []string
	// get archive name from src string
	files = strings.Split(files_str, ",")
	fmt.Printf("%q\n", files) //for debug

	// check the type of archive (.zip , .tar , tar.gz , .tgz) and call archive function
	if strings.HasSuffix(dest_archive_path, ".zip") {
		err := zipop.Archive_zip(files, dest_archive_path)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.HasSuffix(dest_archive_path, ".tar") {
		err := tarop.Archive_tar(files, dest_archive_path)
		if err != nil {
			log.Fatal(err)
		}
	} else if (strings.HasSuffix(dest_archive_path, ".tar.gz")) || (strings.HasSuffix(dest_archive_path, ".tgz")) {
		err := tarop.Archive_targz(files, dest_archive_path)
		if err != nil {
			log.Fatal(err)
		}

	} else if strings.HasSuffix(dest_archive_path, ".rar") {
		//Rar(srcRarArchive, files string)
		rarop.Rar(dest_archive_path, files_str)
	}
	fmt.Println("The new archive stored in: ", dest_archive_path)

	return nil

}

//
func unarchive_handler(src_archive string, url_flag bool, dest_dir_path string) error {

	fmt.Println("unarchive_handler") // for debug
	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
	if strings.HasSuffix(src_archive, ".zip") {

		_, err := zipop.Unzip(src_archive, url_flag, dest_dir_path)
		if err != nil {
			log.Fatal(err)
		}

	} else if (strings.HasSuffix(src_archive, ".tar")) || (strings.HasSuffix(src_archive, ".tar.gz")) || (strings.HasSuffix(src_archive, ".tgz")) {

		tarGzArchiveName, err := tarop.Untar(src_archive, url_flag, dest_dir_path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Archive name: %s\n", tarGzArchiveName)

	} else if strings.HasSuffix(src_archive, ".rar") {

		rarop.Unrar(src_archive, url_flag, dest_dir_path)

	}
	fmt.Printf("unarchive %s stored in: %s\n", src_archive, dest_dir_path)

	return nil
}

//
func extract_handler(src_archive string, url_flag bool, target_ids_files string, target_str_files string, dest_dir_path string) error {

	fmt.Println("extract_handler") // for debug

	var file_ids_str []string
	var files_names []string
	var file_ids []int
	files_names = make([]string, 0)
	file_ids = make([]int, 0)

	// get archive name from src string
	file_ids_str = strings.Split(target_ids_files, ",")
	fmt.Printf("file_ids_str= %s\n", file_ids_str) //for debug
	//convert file ids to file paths
	for _, str_id := range file_ids_str {
		file_id, err := strconv.Atoi(str_id)
		if err != nil {
			log.Fatal(err)
		}
		file_ids = append(file_ids, file_id)
		fmt.Printf("file_ids= %d\n", file_ids) //for debug
	}

	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
	if strings.HasSuffix(src_archive, ".zip") {
		//
		for _, id := range file_ids {
			curr_file_name := tasks_state[curr_task_id].zip_meta[id].FileName
			files_names = append(files_names, curr_file_name)
			fmt.Printf("file_names in loop= %s\n", files_names) //for debug
		}

		err := zipop.Extract_zip(src_archive, url_flag, files_names, dest_dir_path)
		if err != nil {
			log.Fatal(err)
		}
	} else if (strings.HasSuffix(src_archive, ".tar")) || (strings.HasSuffix(src_archive, ".tar.gz")) || (strings.HasSuffix(src_archive, ".tgz")) {

		//
		for _, id := range file_ids {
			curr_file_name := tasks_state[curr_task_id].tar_meta[id].FileName
			files_names = append(files_names, curr_file_name)
			fmt.Printf("file_names in loop= %s\n", files_names) //for debug
		}

		_, err := tarop.Extract_tar(src_archive, url_flag, files_names, dest_dir_path)
		if err != nil {
			log.Fatal(err)
		}

	} else if strings.HasSuffix(src_archive, ".rar") {
		// change  flie names string to be seperated by space instead of ','
		file_str_names := strings.ReplaceAll(target_str_files, ",", " ")
		fmt.Printf("file_names= %s\n", file_str_names) //for debug

		rarop.Extract_rar(src_archive, url_flag, file_str_names, dest_dir_path)

	}
	fmt.Printf("extracted files stored in: %s\n", dest_dir_path)

	return nil

}

//
func meta_handler(src_archive string, url_flag bool) error {

	fmt.Println("meta_handler") // for debug
	//check the type of archive (.zip , .tar , tar.gz , .tgz) and call meta function
	if strings.HasSuffix(src_archive, ".zip") {

		curr_task_zip_meta, err := zipop.Get_zip_meta(src_archive, url_flag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("curr_task_zip_meta= %s\n", curr_task_zip_meta) // for debug

		arr_len := len(curr_task_zip_meta)
		fmt.Printf("arr_len= %d\n", arr_len) // for debug
		//update state
		curr_task_state.zip_meta = make([]zipop.FileMetadata, arr_len)
		curr_task_state.zip_meta = curr_task_zip_meta
		//		fmt.Printf("zip_metadata[10]= %s\n", curr_task_state.zip_meta[10]) //for debug
		tasks_state[curr_task_id] = curr_task_state
		//		fmt.Printf("zip_metadata[10]= %s\n", curr_task_state.zip_meta[10]) //for debug

	} else if (strings.HasSuffix(src_archive, ".tar")) || (strings.HasSuffix(src_archive, ".tar.gz")) || (strings.HasSuffix(src_archive, ".tgz")) {

		curr_task_tar_meta, err := tarop.Get_tar_meta(src_archive, url_flag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("curr_task_targz_metadata= %s\n", curr_task_tar_meta) // temp  for debug
		arr_len := len(curr_task_tar_meta)
		fmt.Printf("arr_len= %d\n", arr_len) //for debug
		//update state
		curr_task_state.tar_meta = make([]tarop.FileMetadata, arr_len)
		curr_task_state.tar_meta = curr_task_tar_meta
		fmt.Printf("tar_meta[10]= %s\n", curr_task_state.tar_meta[10]) //for debug
		tasks_state[curr_task_id] = curr_task_state
		fmt.Printf("tar_meta[10]= %s\n", curr_task_state.tar_meta[10]) // for debug

	} else if strings.HasSuffix(src_archive, ".rar") {

		rarop.Get_rar_meta(src_archive, url_flag)

		// to do - yet not implemented  update state by meta (only print to console)

	}

	return nil
}
func commit_test() {

	// wait 2 seconds and return true with no errors
}

func commit_seedrfs() {

}

func commit_rclone() {

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
	fmt.Println("in start q") //for debug
	defer Wg.Done()

	var line_tokens []string
	var curr_pload Pload

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
				fmt.Printf("in reading: line_tokens= %s\n", line_tokens)
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
	fmt.Println("in start") //for debug
	//	Wg.Add(1) //temp
	//	go start_aria2_server()
	//	time.Sleep(2 * time.Second)
	Wg.Add(1)
	go start_queue()

	//	defer Wg.Done()
	//	Wg.Wait()
	//	go..
	//	go..

}

// parse line from logfile to commands and payload
func parse_line(line string) (line_tokens []string, pload Pload) {
	fmt.Println("in parse") //for debug
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
	var path_or_url string
	var url_flag bool
	url_flag = false //temp
	curr_task_id = line_tokens[0]
	//	curr_event_id = line_tokens[1]
	curr_event_type = line_tokens[2]
	cmd := curr_pload.Cmd
	url := curr_pload.Url
	//	fs_path := curr_pload.Fs_path
	output_folder_id := curr_pload.Output_folder_id
	// to archive cmd - file paths seperated with ',' . to extract cmd - file ids seperated with ',' .
	// to extract rar - file names seperated with ","
	files := curr_pload.Archive_file_ids
	files_to_rar := files
	path := curr_pload.Path
	if url != "" {
		path_or_url = url
		url_flag = true
	} else {
		path_or_url = path
	}

	fmt.Printf("path= %s\n", path)         //for debug
	fmt.Printf("url= %s\n", url)           //for debug
	fmt.Printf("url_flag= %s\n", url_flag) //for debug

	// update state before specific handler
	curr_task_state.curr_event_type = curr_event_type //curr_event_type=="start"
	curr_task_state.progress = ""
	curr_task_state.Payload.Cmd = curr_pload.Cmd
	curr_task_state.Payload.Url = curr_pload.Url
	curr_task_state.Payload.Fs_path = curr_pload.Fs_path
	curr_task_state.Payload.Output_folder_id = curr_pload.Output_folder_id

	tasks_state[curr_task_id] = curr_task_state

	fmt.Printf("tasks_state[curr_task_id]= %s\n", tasks_state[curr_task_id]) //for debug

	if cmd == "archive" && curr_event_type == "start" {

		fmt.Println("in start queue before handler") //for debug

		//call archive handler
		err := archive_handler(files, path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("in start queue after handler") //for debug

		//update state after archive
		curr_task_state.state = "started"
		tasks_state[curr_task_id] = curr_task_state
		fmt.Printf("tasks_state[curr_task_id]= %s\n", tasks_state[curr_task_id]) //for debug

	}
	if cmd == "unarchive" && curr_event_type == "start" {

		//call unarchive handler with local path or url acording to url_flag
		err := unarchive_handler(path_or_url, url_flag, output_folder_id)
		if err != nil {
			log.Fatal(err)
		}

		//to do - update state after unarchive
		curr_task_state.state = "started"
		tasks_state[curr_task_id] = curr_task_state

	}

	if cmd == "extract" && curr_event_type == "start" {

		//call archive handler
		err := extract_handler(path_or_url, url_flag, files, files_to_rar, output_folder_id)
		if err != nil {
			log.Fatal(err)
		}
		// update state after extract

		curr_task_state.state = "started"
		tasks_state[curr_task_id] = curr_task_state
	}
	if cmd == "meta" && curr_event_type == "start" {

		//call archive handler with local path or url acording to url_flag
		err := meta_handler(path_or_url, url_flag)
		if err != nil {
			log.Fatal(err)
		}
		// update state after meta
		curr_task_state.state = "started"
		tasks_state[curr_task_id] = curr_task_state

	}

	if curr_event_type == "finish" {
		//delete curr_task_state
		delete(tasks_state, curr_task_id)
	}

}
