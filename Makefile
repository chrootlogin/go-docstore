DEP=dep
GOLANG=go

all: deps go_app

deps:
	$(DEP) ensure

go_app:
	$(GOLANG) build

clean_testfiles:
	find . -type f -iname data.db -prune -exec rm -f '{}' '+'

test: clean_testfiles
	$(GOLANG) test ./...
	find . -type f -iname data.db -prune -exec rm -f '{}' '+'

clean: clean_testfiles
	rm -rf vendor