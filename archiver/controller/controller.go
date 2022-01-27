package main

// archiver cli , implements 2 functions:
//
// 1. In production stage - starting service.
// 2. In developing stage - testing service functions and library functions.

// in production stage - to run the program from terminal:
// cd ./controller
// go run controller.go "startService"
//
// Tasks are delivered by log file ( example file data_default)
////////////////////////////////////////////////////////////////
// archiver cli commands: create , extract, extractSome , meta.
//
// archiver cli commands syntaxt:
//
//// create	-srcFiles=FilesList  -dest=destPath
//
//		FilesList= filepaths strings seperated by ','
//
//// extract -src=srcPath/url -url=true/false	-dest=destPath
//
//// extractSome -src=filesToExtract -url=true/false -filePath=srcPath/url	-dest=destPath
//
//		filesToExtract = filepaths strings seperated by ','
//
//// meta  -src=srcPath/url -url=true/false	-dest=destPath
//
// example: from the command line
// create	-src=file1,file2	-dest=file3.zip
// or as a code line in the program:
// os.Args = []string{"archiver", "crate", "-src=file1,file2"	"-dest=file3.zip"}

// For starting service by cli in development stage - add next line
// os.Args = []string{"archiver", "startService"}

import (
	"fmt"
	"log"
	"os"
	"strings"

	// include 3 packages for archives operations: zipop , tarop , rarop
	"example.com/archiver/archiverLib/rarop"

	"example.com/archiver/archiverLib/tarop"
	"example.com/archiver/archiverLib/zipop"

	// include other packages
	"example.com/archiver/service"

	"github.com/urfave/cli"
)

var src = "data_default" // log file updating with new events
var srcPath string
var destPath string
var srcFiles []string
var srcFilesRar string
var filePath string
var urlStr string
var isUrl bool

func main() {

	isUrl = false

	app := cli.NewApp()
	app.Name = "archiver"
	app.Commands = []cli.Command{

		//command: start_service
		{
			Name:  "startService",
			Usage: "Starts the archiver service and task runner",
			Action: func(c *cli.Context) error {
				fmt.Println("in cli") //for debug
				service.Start()
				return nil
			},
		},

		//command: create
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create archive of src files into 'dest_path' directory",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archive name from src string
				srcFiles = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", srcFiles)
				// get destDirPath from dest string
				destPath = c.String("dest")

				// check the type of archive (.zip , .tar , tar.gz , .tgz) and call archive function
				if strings.HasSuffix(destPath, ".zip") {
					err := zipop.Create(srcFiles, destPath)
					if err != nil {
						log.Fatal(err)
					}
				} else if strings.HasSuffix(destPath, ".tar") {
					err := tarop.CreateTar(srcFiles, destPath)
					if err != nil {
						log.Fatal(err)
					}
				} else if (strings.HasSuffix(destPath, ".tar.gz")) || (strings.HasSuffix(destPath, ".tgz")) {
					err := tarop.CreateTargz(srcFiles, destPath)
					if err != nil {
						log.Fatal(err)
					}

				}
				fmt.Println("The new archive stored in: ", destPath)
				return nil
			},
		},

		//command: extract
		{
			Name:    "extract",
			Aliases: []string{"e"},
			Usage:   "extract src archive files to dest directory",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "url"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archiveFile name from src string
				srcPath = c.String("src")
				// get url_flag from url string
				urlStr = c.String("url")
				if urlStr == "true" {
					isUrl = true
				}
				// get destDirPath from dest string
				destPath = c.String("dest")
				fmt.Println("src:", srcPath)

				//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(srcPath, ".zip") {
					_, err := zipop.Extract(srcPath, isUrl, destPath)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(srcPath, ".tar")) || (strings.HasSuffix(srcPath, ".tar.gz")) || (strings.HasSuffix(srcPath, ".tgz")) {

					tarArchiveName, err := tarop.Extract(srcPath, isUrl, destPath)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Archive name: %s\n", tarArchiveName)

				} else if strings.HasSuffix(srcPath, ".rar") {

					rarop.Extract(srcPath, isUrl, destPath)

				}

				fmt.Printf("unarchive %s stored in: %s\n", srcPath, destPath)
				return nil
			},
		},

		//command: extractSome
		{
			Name:    "extractSome",
			Aliases: []string{"s"},
			Usage:   "extract single file from archive",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				&cli.StringFlag{Name: "url"},
				&cli.StringFlag{Name: "filePath"},
				&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				// get archive file name from src string
				srcPath = c.String("src")
				fmt.Println("srcPath:", srcPath)
				// get archive name from src string
				srcFiles = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", srcFiles)

				// get url_flag from url string
				urlStr = c.String("url")
				if urlStr == "true" {
					isUrl = true
				}
				// get filepath to extract from archive
				filePath = c.String("filePath")
				fmt.Println("filePath:", filePath)
				// get destDirPath from dest string
				destPath = c.String("dest")
				fmt.Println("dest:", destPath)

				////check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(srcPath, ".zip") {
					err := zipop.ExtractSome(srcPath, isUrl, srcFiles, destPath)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(srcPath, ".tar")) || (strings.HasSuffix(srcPath, ".tar.gz")) || (strings.HasSuffix(srcPath, ".tgz")) {

					_, err := tarop.ExtractSome(srcPath, isUrl, srcFiles, destPath)
					if err != nil {
						log.Fatal(err)
					}

				} else if strings.HasSuffix(srcPath, ".rar") {
					// files names seperated with space instead of ","
					srcFilesRar = strings.ReplaceAll(srcPath, ",", " ")
					rarop.ExtractSome(srcPath, isUrl, srcFilesRar, destPath)

				}
				fmt.Printf("file %s stored in: %s\n", filePath, destPath)
				return nil
			},
		},

		//command: meta
		{
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
				srcPath = c.String("src")
				// get url_flag from url string
				urlStr = c.String("url")
				if urlStr == "true" {
					isUrl = true
				}
				// get destDirPath from dest string
				destPath = c.String("dest")

				//check the type of archive (.zip , .tar , tar.gz , .tgz) and call unarchive function
				if strings.HasSuffix(srcPath, ".zip") {
					_, err := zipop.Meta(srcPath, isUrl)
					if err != nil {
						log.Fatal(err)
					}

				} else if (strings.HasSuffix(srcPath, ".tar")) || (strings.HasSuffix(srcPath, ".tar.gz")) || (strings.HasSuffix(srcPath, ".tgz")) {

					_, err := tarop.Meta(srcPath, isUrl)
					if err != nil {
						log.Fatal(err)
					}

				} else if strings.HasSuffix(srcPath, ".rar") {

					rarop.Meta(srcPath, isUrl)

				}
				fmt.Printf("unarchive %s stored in: %s\n", srcPath, destPath)
				return nil
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	// wait for all goroutines to finish before ending main
	service.Wg.Wait()
}
