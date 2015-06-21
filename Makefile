BUILDOPT := -ldflags '-s -w'
SOURCES  := $(wildcard */*.go)

gom:
	go get github.com/mattn/gom
	gom install

run:
	gom run main.go ${ARGS}

fmt:
	@$(foreach FILE, $(SOURCES), gom exec goimports -w $(FILE);)

build: fmt $(SOURCES)
	@$(foreach FILE, $(SOURCES), echo $(FILE); gom build $(BUILDOPT) -o bin/`basename $(FILE) .go` $(FILE);)

clean:
	rm -f bin/*
