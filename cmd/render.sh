#!/bin/bash

# This script demonstrates example useage of the fractal-engine binary.

EXEC=fractal-engine
[ ! -f ${EXEC} ] && echo "executable not found" && exit 1;

# this renders an image
: <<'END'
./${EXEC} --mode=img                         \
          --img-width=1024 --img-height=1024 \
          --plot-width=3   --plot-height=3   \
          --palette=bw                       \
          --julia-exp=4                      \
          --julia-escape-rad=2               \
          --init-iteratex=0                  \
          --init-iteratey=0                  \
          --of-name=renders/render.png
END

# this renders a gif animation
./${EXEC} --mode=gif                         \
          --img-width=2048 --img-height=1024 \
          --plot-width=4   --plot-height=2   \
          --plot-cx=0   --plot-cy=1          \
          --palette=uf                       \
          --julia-exp=2                      \
          --julia-escape-rad=2               \
          --num-frames=300                   \
          --iterations=2500                  \
          --progress                         \
          --of-name=renders/render-uf.gif
