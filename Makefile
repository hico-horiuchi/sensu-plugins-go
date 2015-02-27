BUILDOPT := -ldflags '-s -w'
SOURCES  := $(wildcard *.go)

build: $(SOURCES)
	@$(foreach FILE, $(SOURCES), echo $(FILE); go build $(BUILDOPT) -o bin/`basename $(FILE) .go` $(FILE);)

clean:
	rm -f bin/*
