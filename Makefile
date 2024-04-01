SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

.PHONY: all 
all: out/teamscli.x86_64.bin out/teamscli.x86_64.exe

out/teamscli.x86_64.bin: $(SOURCES)
	go build -trimpath  -ldflags "-s -w"  -o $@ .
	strip -s $@

out/teamscli.x86_64.exe: $(SOURCES)
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o $@ .
	strip -s $@
