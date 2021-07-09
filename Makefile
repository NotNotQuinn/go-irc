ifeq ($(OS),Windows_NT)
	EXE=.exe
else
	EXE=
endif

OUTFILE=bin/wanductbot$(EXE)

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

push: lint-test
#   Hahahah, git in my makefile!
	git push

test:
	go test ./...

lint-test: lint test
