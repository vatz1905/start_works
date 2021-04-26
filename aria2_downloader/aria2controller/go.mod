module example.com/aria2controller

go 1.13

replace example.com/aria2Downloader => ../aria2Downloader

replace example.com/service => ../service

require (
	example.com/aria2Downloader v0.0.0-00010101000000-000000000000
	example.com/service v0.0.0-00010101000000-000000000000
	github.com/urfave/cli v1.22.5
)
