_default: bin/gecko
	@:  # if we have a command this silences "nothing to be done"

bin/gecko: gecko/*.go # help: run the server
	go build -o bin/gecko

clean:
	rm -f bin/gecko