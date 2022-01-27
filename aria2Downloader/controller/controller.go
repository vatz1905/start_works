package main

// aria2Downloader cli , implements 2 functions:
//
// 1. In production stage - starting service.
// 2. In developing stage - testing service functions and library functions.

// in production stage - to run the program from terminal:
// cd ./controller
// go run controller.go "startService"
//
// Tasks are delivered by log file ( example file: data_default)
//
//
//  aria2Downloader cli commands: start , pause, resume , status, remove.
//
// commands syntax:
// downloader start-src=uri1,uri2,..
// downloader pause		-src=gid
// downloader resume	-src=gid
// downloader status	-src=gid
// downloader remove	-src=gid

// Using cli in development stage, add command by os.Arg , for example:
// os.Args = []string{"downloader", "start", "-src=https://github.com/resource1,https://github.com/resource2"}
// os.Args = []string{"downloader", "remove", "-src=downloadGid"}

import (
	"fmt"
	"log"
	"os"
	"strings"

	//example.com have to be replace by relrvant repository path
	"example.com/aria2Downloader/downloader"
	"example.com/aria2Downloader/service"
	"github.com/urfave/cli"
)

var srcUris []string
var srcGid string

func main() {

	app := cli.NewApp()
	app.Name = "downloader"
	app.Commands = []cli.Command{
		{
			Name:  "startService",
			Usage: "Starts the downloader service and task runner",
			Action: func(c *cli.Context) error {
				service.Start()
				return nil
			},
		},

		//command: download
		// command format: downloader download -src=urisList
		// srcUris seperated by ','
		{
			Name:    "start",
			Aliases: []string{"sd"},
			Usage:   "add download to 'downloads queue',start download and and return gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},
			Action: func(c *cli.Context) error {

				var retGid string

				// get uris list from src string
				srcUris = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", srcUris)

				// call Download function
				retGid = downloader.Start(srcUris)
				fmt.Println("ret_gid ", retGid)

				return nil
			},
		},

		//command: remove
		//command format: downloader remove -src=gid
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "remove download from queue by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},
			Action: func(c *cli.Context) error {

				var retGid string

				// get 'download gid' from src string
				srcGid = c.String("src")
				fmt.Println("src:", srcGid)

				//call removeDl function
				retGid = downloader.Remove(srcGid)
				fmt.Printf("ret_gid= %s\n", retGid)

				return nil
			},
		},

		//command: pause
		//command format: downloader pause   -src=gid
		{
			Name:    "pause",
			Aliases: []string{"p"},
			Usage:   "pause download by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},
			Action: func(c *cli.Context) error {

				var retGid string

				// get 'download gid' from src string
				srcGid = c.String("src")
				fmt.Println("srcGid:", srcGid)

				// call pauseDl function
				retGid = downloader.Pause(srcGid)

				fmt.Printf("retgid= %s\n", retGid)
				return nil
			},
		},

		//command: resume
		//command format: downloader unpause	-src=gid
		{
			Name:    "resume",
			Aliases: []string{"rs"},
			Usage:   "resume download by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},

			Action: func(c *cli.Context) error {

				var retGid string

				// get 'download gid' from src string
				srcGid = c.String("src")
				fmt.Println("srcGid= \n", srcGid)

				// call unpauseDl function
				retGid = downloader.Resume(srcGid)
				fmt.Println("retgid= \n", retGid)

				return nil
			},
		},

		//command status
		//command format: downloader getStatus	-src=gid
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "get download status by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},

			Action: func(c *cli.Context) error {

				// get 'download gid' from src string
				srcGid = c.String("src")
				fmt.Println("srcGid= \n", srcGid)

				// get download status by its gid
				downloader.DownloadStatus(srcGid) // aria2Downloader.getStatus(srcGid)

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
