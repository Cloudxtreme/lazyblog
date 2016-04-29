release:
	mkdir mac linux windows

	# Mac 64
	env GOOS=darwin GOARCH=amd64 go build -v -o lazyblog
	mv lazyblog mac/lazyblog

	# Linux 64
	# env GOOS=linux GOARCH=amd64 go build -v -o lazyblog

	# Linux 32
	# env GOOS=linux GOARCH=386 go build -v -o lazyblog

	# Linux ARM
	env GOOS=linux GOARCH=arm go build -v -o lazyblog
	mv lazyblog linux/lazyblog

	# Windows 64
	env GOOS=windows GOARCH=amd64 go build -v -o lazyblog

	# Windows 32
	env GOOS=windows GOARCH=386 go build -v -o lazyblog
	mv lazyblog windows/lazyblog


clean:
	rm -rf mac linux

.PHONY: clean
