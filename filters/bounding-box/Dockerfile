FROM pixeltilt-base

COPY ../../render/api render/api/
COPY filters/bounding-box/ rectangler/

RUN --mount=type=cache,target=/cache/gomod \
    --mount=type=cache,target=/cache/gobuild,sharing=locked \
    go mod vendor && \
    go build -mod=vendor -o /usr/local/bin/rectangler ./rectangler

CMD ["/usr/local/bin/rectangler"]
