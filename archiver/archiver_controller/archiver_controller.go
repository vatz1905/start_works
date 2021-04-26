package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	// include 3 packages for archives operations: zipop , tarop , rarop
	"example.com/archiver_lib/rarop"

	"example.com/archiver_lib/tarop"
	"example.com/archiver_lib/zipop"

	// include other packages
	"example.com/archiver_service"

	"github.com/urfave/cli"
)

var src = "data_default" // log file updating with new events
var src_path string
var dest_path string
var src_files []string
var src_files_rar string
var file_path string
var url_str string
var url_flag bool

// archiver cli , implements 2 functions:
// 1. In developing stage - testing service functions and library functions.
// 2. In production stage - starting archiver_service.
//
// archiver cli commands: archive , unarchive, extract , meta.
//
// archiver cli commands syntaxt:
//
//// archiver archive	-src_files  -dest=dest_path
//
//		src_files- filepaths strings seperated by ','
//
//// archiver unarchive -src=src_path/url -url_flag=true/false	-dest=dest_path
//
//// archiver extract -src=files_to_extract -url_flag=true/false -file_path=src_path/url	-dest=dest_path
//
//		files_to_extract - filepaths strings seperated by ','
//
//// archiver meta  -src=src_path/url -url_flag=true/false	-dest=dest_path
//
// example: from the command line
// archiver	archive	-src=file1,file2	-dest=file3.zip
// or as a code line in the program:
// os.Args = []string{"archiver", "archive", "-src=file1,file2"	"-dest=file3.zip"}

// For starting archiver_service by cli in development stage - add next line
// os.Args = []string{"archiver", "start_service"}

func main() {

	url_flag = false

	app := cli.NewApp()
	app.Name = "archiver"
	app.Commands = []cli.Command{
		{
			Name:  "start_service",
			Usage: "Starts the archiver service and task runner",
			Action: func(c *cli.Context) error {
				fmt.Println("in cli") //for debug
				archiver_service.Start()
				return nil
			},
		},
		{ //command: archive
			Name:    "archive",
			Aliases: []string{"a"},
			Usage:   "create archive of src files into 'dest_path' directory",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archive name from src string
				src_files = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", src_files)
				// get destDirPath from dest string
				dest_path = c.String("dest")

				// check the type of archive (.zip , .tar , tar.gz , .tgz) and call archive function
				if strings.HasSuffix(dest_path, ".zip") {
					err := zipop.Archive_zip(src_files, dest_path)
					if err != nil {
						log.Fatal(err)
					}
				} else if strings.HasSuffix(dest_path, ".tar") {
					err := tarop.Archive_tar(src_files, dest_path)
					if err != nil {
						log.Fatal(err)
					}
				} else if (strings.HasSuffix(dest_path, ".tar.gz")) || (strings.HasSuffix(dest_path, ".tgz")) {
					err := tarop.Archive_targz(src_files, dest_path)
					if err != nil {
						log.Fatal(err)
					}

				}
				fmt.Println("The new archive stored in: ", dest_path)
				return nil
			},
		},
		{ //command: unarchive
			Name:    "unarchive",
			Aliases: []string{"u"},
			Usage:   "unarchive src archive file to dest directory",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "url"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archiveFile name from src string
				src_path = c.String("src")
				// get url_flag from url string
				url_str = c.String("url")
				if url_str == "true" {
					url_flag = true
				}
				// get destDirPath from dest string
				dest_path = c.String("dest")
				fmt.Println("src:", src_path)

				//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(src_path, ".zip") {
					_, err := zipop.Unzip(src_path, url_flag, dest_path)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(src_path, ".tar")) || (strings.HasSuffix(src_path, ".tar.gz")) || (strings.HasSuffix(src_path, ".tgz")) {

					tarArchiveName, err := tarop.Untar(src_path, url_flag, dest_path)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Archive name: %s\n", tarArchiveName)

				} else if strings.HasSuffix(src_path, ".rar") {

					rarop.Unrar(src_path, url_flag, dest_path)

				}

				fmt.Printf("unarchive %s stored in: %s\n", src_path, dest_path)
				return nil
			},
		},

		{ //command: extract
			Name:    "extract",
			Aliases: []string{"t"},
			Usage:   "extract single file from archive",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "url"},
				&cli.StringFlag{Name: "file_path"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archive file name from src string
				src_path = c.String("src")
				fmt.Println("src_path:", src_path)
				// get archive name from src string
				src_files = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", src_files)

				// get url_flag from url string
				url_str = c.String("url")
				if url_str == "true" {
					url_flag = true
				}
				// get filepath to extract from archive
				file_path = c.String("file_path")
				fmt.Println("file_path:", file_path)
				// get destDirPath from dest string
				dest_path = c.String("dest")
				fmt.Println("dest:", dest_path)

				////check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(src_path, ".zip") {
					err := zipop.Extract_zip(src_path, url_flag, src_files, dest_path)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(src_path, ".tar")) || (strings.HasSuffix(src_path, ".tar.gz")) || (strings.HasSuffix(src_path, ".tgz")) {

					_, err := tarop.Extract_tar(src_path, url_flag, src_files, dest_path)
					if err != nil {
						log.Fatal(err)
					}

				} else if strings.HasSuffix(src_path, ".rar") {
					// files names seperated with space instead of ","
					src_files_rar = strings.ReplaceAll(src_path, ",", " ")
					rarop.Extract_rar(src_path, url_flag, src_files_rar, dest_path)

				}
				fmt.Printf("file %s stored in: %s\n", file_path, dest_path)
				return nil
			},
		},

		{ //command: meta
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "get metadata of src archive file to dest directory",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "url"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archiveFile name from src string
				src_path = c.String("src")
				// get url_flag from url string
				url_str = c.String("url")
				if url_str == "true" {
					url_flag = true
				}
				// get destDirPath from dest string
				dest_path = c.String("dest")

				//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(src_path, ".zip") {
					_, err := zipop.Get_zip_meta(src_path, url_flag)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(src_path, ".tar")) || (strings.HasSuffix(src_path, ".tar.gz")) || (strings.HasSuffix(src_path, ".tgz")) {

					_, err := tarop.Get_tar_meta(src_path, url_flag)
					if err != nil {
						log.Fatal(err)
					}

				} else if strings.HasSuffix(src_path, ".rar") {

					rarop.Get_rar_meta(src_path, url_flag)

				}
				fmt.Printf("unarchive %s stored in: %s\n", src_path, dest_path)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	archiver_service.Wg.Wait() // wait for all goroutines to finish before ending main
}
