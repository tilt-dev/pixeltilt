FROM pixeltilt-base

COPY ../render/api render/api/
COPY ../storage/api storage/api/
COPY ../storage/client storage/client/
COPY muxer/ muxer/

RUN --mount=type=cache,target=/cache/gomod \
    --mount=type=cache,target=/cache/gobuild,sharing=locked \
    find . -name "*.go" && \
    go mod vendor && \
    go build -mod=vendor -o /usr/local/bin/muxer ./muxer

CMD ["/usr/local/bin/muxer"]
