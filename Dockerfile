FROM scratch
ARG ARCH=amd64
ADD ./build/${ARCH} /app
ENTRYPOINT ["/app"]
