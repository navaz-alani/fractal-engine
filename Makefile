GC=go
EXEC=fractal-engine
CMD_DIR=cmd
SRCS=$(wildcard *.go cmd/*.go)

${EXEC}: ${SRCS}
	go build -o ${EXEC} ./${CMD_DIR}

.PHONY: clean

clean:
	rm -rf ${EXEC}
