release:
	mkdir mac linux windows

	# Mac 64
	env GOOS=darwin GOARCH=amd64 go build -v -o lazyblog
	mv lazyblog mac/lazyblog

	# Linux 64
	env GOOS=linux GOARCH=amd64 go build -v -o lazyblog
	mv lazyblog linux/lazyblog64

	# Linux 32
	env GOOS=linux GOARCH=386 go build -v -o lazyblog
	mv lazyblog linux/lazyblog32

	# Linux ARM
	env GOOS=linux GOARCH=arm go build -v -o lazyblog
	mv lazyblog linux/lazyblogarm

	# Windows 64
	env GOOS=windows GOARCH=amd64 go build -v -o lazyblog
	mv lazyblog windows/lazyblog64

	# Windows 32
	env GOOS=windows GOARCH=386 go build -v -o lazyblog
	mv lazyblog windows/lazyblog32


clean:
	rm -rf mac linux windows

.PHONY: clean
