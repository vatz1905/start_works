module example.com/archiver_controller

go 1.13

replace example.com/archiver_lib/zipop => ../archiver_lib/zipop

replace example.com/archiver_lib/tarop => ../archiver_lib/tarop

replace example.com/archiver_lib/rarop => ../archiver_lib/rarop

require (
	example.com/archiver_lib/rarop v0.0.0-00010101000000-000000000000
	example.com/archiver_lib/tarop v0.0.0-00010101000000-000000000000
	example.com/archiver_lib/zipop v0.0.0-00010101000000-000000000000
	example.com/archiver_service v0.0.0-00010101000000-000000000000
	github.com/urfave/cli v1.22.5
)

replace example.com/archiver_service => ../archiver_service
