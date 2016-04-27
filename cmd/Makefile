release:
	mkdir mac linux windows
	env GOOS=darwin GOARCH=amd64 go build -v -o lazyblog
	mv lazyblog mac/lazyblog
	env GOOS=linux GOARCH=arm GOARM=7 go build -v -o lazyblog
	mv lazyblog linux/lazyblog
	env GOOS=windows GOARCH=386 go build -v -o lazyblog
	mv lazyblog windows/lazyblog

clean:
	rm -rf mac linux

.PHONY: clean
