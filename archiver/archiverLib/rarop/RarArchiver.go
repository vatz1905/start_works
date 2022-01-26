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

// Create archive.rar from files string seperated by space, example: "file1 file2 file3"
func Create(srcArchive, files string) {

	cmd := exec.Command("rar", "a", srcArchive, files)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

// Extract local archive.rar to destDir
func Extract(srcArchive string, isUrl bool, destDir string) {
	if !isUrl {
		cmd := exec.Command("unrar", "e", srcArchive, destDir)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else if isUrl {
		// download rar archive by its url, to pwd.
		fileURL, err := url.Parse(srcArchive)
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
		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			log.Fatal(err)
		}
		cmd := exec.Command("unrar", "e", fileName, destDir)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ExtractSome files from local rar archive to destDir
// filesToExtract- files string seperated by space, example: "file1 file2 file3"
func ExtractSome(srcArchive string, isUrl bool, filesToExtract string, destDir string) {
	if !isUrl {
		cmd := exec.Command("unrar", "e", srcArchive, filesToExtract, destDir)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("local path only")
	}
}

// display local archive.rar content
func Meta(srcArchive string, isUrl bool) {
	if !isUrl {
		cmd := exec.Command("unrar", "l", srcArchive)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("local path only")
	}
}

//extract local split archive.rar by its first part: file.part1.rar
func ExtractSplit(srcSplitRar string) {

	//for split rar file use:  unrar x -e file.part1.rar
	// all other parts in the same dir as part1 and unrar prog will find them automatically

	cmd := exec.Command("unrar", "x", "-e", srcSplitRar)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

}
