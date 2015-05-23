BUILDOPT := -ldflags '-s -w'
SOURCES  := $(wildcard */*.go)

gom:
	go get github.com/mattn/gom
	gom install

run:
	gom run main.go ${ARGS}

fmt:
	gom exec goimports -w check/*.go metrics/*.go handler/*.go

build: fmt $(SOURCES)
	@$(foreach FILE, $(SOURCES), echo $(FILE); gom build $(BUILDOPT) -o bin/`basename $(FILE) .go` $(FILE);)

clean:
	rm -f bin/*
