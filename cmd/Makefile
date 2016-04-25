release:
	mkdir mac linux
	go build -v -o lazyblog
	mv lazyblog mac/lazyblog
	env GOOS=linux GOARCH=arm GOARM=7 go build -v -o lazyblog
	mv lazyblog linux/lazyblog

clean:
	rm -rf mac linux

.PHONY: clean
