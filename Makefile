OUTFILE=bin/Wanductbot.exe
TEST_OUTFILE=bin/Wanductbot.test.exe

build:
	go build -o $(OUTFILE) .
clean:
	rm -rf bin/
run: build
	$(OUTFILE)

lint:
#   To install this lint tool
#   See this for your operating system
#   https://golangci-lint.run/usage/install/
	golangci-lint run ./...

push: lint
#   Hahahah, git in my makefile!
	git push

test:
#	Needed because working directory is depended upon
	go test ./...
