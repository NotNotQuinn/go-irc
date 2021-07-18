# Check OS
ifeq ($(OS),Windows_NT)
	EXE=.exe
else
	EXE=
endif

BOT_OUTFILE=bin/wanductbot$(EXE)
POPULATOR_OUTFILE=bin/populator$(EXE)

bot:
	go build -o $(BOT_OUTFILE) .

populator:
	go build -o $(POPULATOR_OUTFILE) ./data/populator

clean:
	rm -rf bin/

run: bot
	$(BOT_OUTFILE)

lint:
# To install this lint tool
# See this for your operating system
# https://golangci-lint.run/usage/install/
	golangci-lint run ./...

push: lint-test
# Hahahah, git in my makefile!
	git push

test:
# Only build for now, may publish on docker hub and use an image.
	docker compose build
	docker compose up --abort-on-container-exit

reset-test:
	docker compose down
	docker volume rm twitch-bot_mariadb_test_data

lint-test: lint test
all: bot populator
