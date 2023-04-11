sprites:
	@go env GOOS=windows
	@go run .\pkg\packer\cmd\main.go --input ./images --stats
