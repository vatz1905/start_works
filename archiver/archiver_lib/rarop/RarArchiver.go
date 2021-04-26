package rarop

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	srcRarDir   string
	dstUnrarDir string
	srcSplitRar string
)

// archive archive.rar from files string seperated by space, example: "file1 file2 file3"
func Rar(srcRarArchive, files string) {

	cmd := exec.Command("rar", "a", srcRarArchive, files)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

// unrar local archive.rar to dstUnrarDir
func Unrar(srcRarArchive string, url_flag bool, dstUnrarDir string) {
	if !url_flag {
		cmd := exec.Command("unrar", "e", srcRarArchive, dstUnrarDir)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else if url_flag {
		// download rar archive by its url, to pwd.
		fileURL, err := url.Parse(srcRarArchive)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		// Create blank file
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		// http get request
		resp, err := http.Get(srcRarArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("unrar", "e", fileName, dstUnrarDir)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// extract files from local rar archive to dstUnrarDir
// filesToExtract- files string seperated by space, example: "file1 file2 file3"
func Extract_rar(srcRarArchive string, url_flag bool, filesToExtract string, dstUnrarDir string) {
	if !url_flag {
		cmd := exec.Command("unrar", "e", srcRarArchive, filesToExtract, dstUnrarDir)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("local path only")
	}
}

// display local archive.rar content
func Get_rar_meta(srcRarArchive string, url_flag bool) {
	if !url_flag {
		cmd := exec.Command("unrar", "l", srcRarArchive)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("local path only")
	}
}

//unrar local split archive.rar by its first part: file.part1.rar
func UnrarSplit(srcSplitRar string) {

	//for split rar file use:  unrar x -e file.part1.rar
	// all other parts in the same dir as part1 and unrar prog will find them automatically

	cmd := exec.Command("unrar", "x", "-e", srcSplitRar)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}
