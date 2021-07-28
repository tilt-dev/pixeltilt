FROM pixeltilt-base

COPY ../../render/api render/api/
COPY filters/glitch/ glitch/

RUN --mount=type=cache,target=/cache/gomod \
    --mount=type=cache,target=/cache/gobuild,sharing=locked \
    go mod vendor && \
    go build -mod=vendor -o /usr/local/bin/glitch ./glitch

CMD ["/usr/local/bin/glitch"]
