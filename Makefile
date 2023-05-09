sprites:
	@go env GOOS=windows
	@go run .\pkg\packer\cmd\main.go --input ./images --stats

sprites-linux:
	@go env GOOS=linux
	@go run .\pkg\packer\cmd\main.go --input ./images --stats

sprites-darwin:
	@go env GOOS=darwin
	@go run .\pkg\packer\cmd\main.go --input ./images --stats
