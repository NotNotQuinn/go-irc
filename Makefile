build:
	go build -o bin/bot.exe .
clean:
	rm -rf bin/
run: build
	bin/bot.exe