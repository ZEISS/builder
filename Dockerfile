# https://goreleaser.com/docker/

FROM gcr.io/distroless/static:nonroot

ARG BINARY=builder
ARG TARGETPLATFORM

WORKDIR /
COPY $TARGETPLATFORM/$BINARY /main

USER 65532:65532

CMD ["/main"]
