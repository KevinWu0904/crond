all: crond-api

crond-api:
	mkdir -p ./bin/
	go build -o ./bin/ ./cmd/crond-api/

clean:
	rm -rf ./bin/