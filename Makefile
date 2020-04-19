.PHONY: genast
genast:
	go build script/genast.go
	./genast internal/

.PHONY: glox
glox:
	go build cmd/glox.go
