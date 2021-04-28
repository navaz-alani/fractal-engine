GC=go

ENG=fractal-engine
ENG_DIR=cmd/fractal-engine

JULIA_ANIM=julia-animation
ANIM_DIR=cmd/julia-anim

SRCS=$(wildcard *.go)

${JULIA_ANIM}: ${SRCS} $(wildcard ${ANIM_DIR}/*.go)
	${GC} build -o ${JULIA_ANIM} ./${ANIM_DIR}

${ENG}: ${SRCS} $(wildcard ${ENG_DIR}/*.go)
	${GC} build -o ${EXEC} ./${ENG_DIR}

.PHONY: clean

clean:
	rm -rf ${ENG}
	rm -rf ${JULIA_ANIM}
