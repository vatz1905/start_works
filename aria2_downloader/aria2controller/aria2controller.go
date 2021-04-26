package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/aria2Downloader" // example.com/ have to be replace by relrvant repository path
	"example.com/service"         //example.com/ have to be replace by relrvant repository path
	"github.com/urfave/cli"
)

var src_uris []string
var src_gid string

func main() {

	// for testing
	//os.Args = []string{"downloader", "addDl", "-src=https://github.com/pyload/pyload/wiki/module.Api-(SourceCode)"}
	os.Args = []string{"downloader", "start_service"}
	//os.Args = []string{"downloader", "removeDl", "-src=gid"}
	//	os.Args = []string{"downloader", "addDl", "-src=https://github.com/pyload/pyload/wiki/module.Api-(SourceCode)",}

	//	os.Args = []string{"downloader", "pauseDl", "-src=gid",}
	//	os.Args = []string{"downloader", "unpauseDl", "-src=gid",}
	//os.Args = []string{"downloader", "getStatus", "-src=gid"}
	// end "for testing"

	// commands syntax
	// downloader addDl  -src=uri1,uri2,..}
	// downloader removeDl -src=gid
	// downloader pauseDl   -src=gid
	// downloader unpauseDl	-src=gid
	// downloader getStatus	-src=gid

	app := cli.NewApp()
	app.Name = "downloader"
	app.Commands = []cli.Command{
		{
			Name:  "start_service",
			Usage: "Starts the downloader service and task runner",
			Action: func(c *cli.Context) error {
				service.Start()
				return nil
			},
		},
		{ //command: download
			// command format: downloader download -src=urisList
			// src- uris seperated by ','
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "add download to 'downloads queue',start download and and return gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
				//	&cli.StringFlag{Name: "dest"},
			},
			Action: func(c *cli.Context) error {
				var ret_gid string
				// get urislist from src string
				src_uris = strings.Split(c.String("src"), ",")
				fmt.Printf("%q\n", src_uris)
				// call addDl function
				ret_gid = aria2Downloader.Download(src_uris)
				fmt.Println("ret_gid ", ret_gid)
				return nil
			},
		},
		{ //command: remove
			//command format: downloader remove -src=gid
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove download from queue by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},
			Action: func(c *cli.Context) error {
				var ret_gid string
				// get 'download gid' from src string
				src_gid = c.String("src")

				fmt.Println("src:", src_gid)
				//call removeDl function

				ret_gid = aria2Downloader.Remove_dl(src_gid)

				fmt.Printf("ret_gid= %s\n", ret_gid)
				return nil
			},
		},

		{ //command: pause
			//command format: downloader pause   -src=gid
			Name:    "pause",
			Aliases: []string{"p"},
			Usage:   "pause download by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},
			Action: func(c *cli.Context) error {
				var ret_gid string
				// get archive file name from src string
				src_gid = c.String("src")
				fmt.Println("srcGid:", src_gid)

				// call pauseDl function
				ret_gid = aria2Downloader.Pause_dl(src_gid)

				fmt.Printf("retgid= %s\n", ret_gid)
				return nil
			},
		},

		{ //command: unpause
			//command format: downloader unpause	-src=gid

			Name:    "unpause",
			Aliases: []string{"u"},
			Usage:   "unpause download by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},

			Action: func(c *cli.Context) error {
				var ret_gid string
				// get rar archive name from src string
				src_gid = c.String("src")
				fmt.Println("srcGid= \n", src_gid)
				// call unpauseDl function
				ret_gid = aria2Downloader.Unpause_dl(src_gid)

				fmt.Println("retgid= \n", ret_gid)
				return nil
			},
		},

		{ //command getStatus
			// command format: downloader getStatus	-src=gid
			Name:    "getStatus",
			Aliases: []string{"s"},
			Usage:   "get download status by its gid",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "src"},
			},

			Action: func(c *cli.Context) error {
				// get download status by its gid, from src string
				src_gid = c.String("src")
				fmt.Println("srcGid= \n", src_gid)

				//fmt.Println("create archive to:", destPath)
				aria2Downloader.GetStatus_dl(src_gid) // aria2Downloader.getStatus(srcGid)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	service.Wg.Wait() // wait for all goroutines to finish before ending main
}
