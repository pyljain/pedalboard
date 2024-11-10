# Build and run Pedalboard
.PHONY: run
run: build
	./bin/pedalboard --config ./examples/config.yaml
	
.PHONY: clean
clean:
	rm -rf ./bin && mkdir ./bin

.PHONY: build
build: clean
	go build -o ./bin/pedalboard

