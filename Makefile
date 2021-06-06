build:
	go build -o bin/bot.exe .
clean:
	rm -rf bin/
run: build
	bin/bot.exe

lint:
# To install this lint tool
# See this for your operating system
# https://golangci-lint.run/usage/install/
	golangci-lint run ./...

push: lint
# Hahahah, git in my makefile!
	git push
