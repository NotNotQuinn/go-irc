build:
	go build -o bin/bot.exe .
clean:
	rm -rf bin/
run:
	make clean
	make build
	bin/bot.exe