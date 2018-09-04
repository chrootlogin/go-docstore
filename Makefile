DEP=dep
GOLANG=go

all: deps go_app

deps:
	$(DEP) ensure

go_app:
	$(GOLANG) build github.com/chrootlogin/go-docstore/cmd/go-docstore-server

test:
	$(GOLANG) test ./...