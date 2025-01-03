FROM --platform=${TARGETPLATFORM} busybox:latest
COPY scaf /bin/scaf
ENTRYPOINT ["/bin/scaf"]
