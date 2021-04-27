#!/bin/bash

# This script demonstrates example useage of the fractal-engine binary.

EXEC=fractal-engine
[ ! -f ${EXEC} ] && echo "executable not found" && exit 1;

# this renders an image
: <<'END'
./${EXEC} --mode=img                         \
          --img-width=4096 --img-height=2048 \
          --plot-width=6   --plot-height=3   \
          --palette=bw                       \
          --julia-exp=2                      \
          --julia-escape-rad=2               \
          --init-iteratex=1                  \
          --init-iteratex=1                  \
          --of-name=renders/render.png
END

# this renders a gif animation
./${EXEC} --mode=gif                         \
          --img-width=1024 --img-height=512 \
          --plot-width=4   --plot-height=2   \
          --plot-cx=0   --plot-cy=1          \
          --palette=bw                       \
          --julia-exp=2                      \
          --julia-escape-rad=2               \
          --num-frames=250                   \
          --progress                         \
          --of-name=renders/render-bw.gif
