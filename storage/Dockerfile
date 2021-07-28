FROM pixeltilt-base

COPY storage/ storage/

RUN --mount=type=cache,target=/cache/gomod \
    --mount=type=cache,target=/cache/gobuild,sharing=locked \
    go mod vendor && \
    go build -mod=vendor -o /usr/local/bin/storage ./storage

CMD ["/usr/local/bin/storage"]
