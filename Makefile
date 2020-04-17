.PHONY: genast
genast:
	go build script/genast.go
	./genast internal/
