PROJECT=json-example
SRC=$(wildcard *.go)


.PHONY: run
run: $(PROJECT)
	@./$(PROJECT) test.json


$(PROJECT): $(SRC)
	go build
