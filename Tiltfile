# version_settings() enforces a minimum Tilt version
# https://docs.tilt.dev/api.html#api.version_settings
version_settings(constraint='>=0.22.1')

# load() can be used to split your Tiltfile logic across multiple files
# the special ext:// prefix loads the corresponding extension from
# https://github.com/tilt-dev/tilt-extensions instead of a local file
load('ext://restart_process', 'docker_build_with_restart')
load('ext://uibutton', 'cmd_button')

# k8s_yaml automatically creates resources in Tilt for the entities
# and will inject any images referenced in the Tiltfile when deploying
# https://docs.tilt.dev/api.html#api.k8s_yaml
k8s_yaml([
    'filters/glitch/k8s.yaml',
    'filters/color/k8s.yaml',
    'filters/bounding-box/k8s.yaml',
    'storage/k8s.yaml',
    'muxer/k8s.yaml',
    'object-detector/k8s.yaml',
    'frontend/k8s.yaml',
])

# k8s_resource allows customization where necessary such as adding port forwards
# https://docs.tilt.dev/api.html#api.k8s_resource
k8s_resource("frontend", port_forwards="3000", labels=["frontend"])
k8s_resource("storage", port_forwards="8080", labels=["infra"])
k8s_resource("max-object-detector", labels=["infra"], new_name="object-detector")
k8s_resource("glitch", labels=["backend"])
k8s_resource("color", labels=["backend"])
k8s_resource("bounding-box", labels=["backend"])
k8s_resource("muxer", labels=["backend"])

# cmd_button extension adds custom buttons to a resource to execute tasks on demand
# https://github.com/tilt-dev/tilt-extensions/tree/master/uibutton
cmd_button(
    name='flush-storage',
    resource='storage',
    argv=['curl', '-s', 'http://localhost:8080/flush'],
    text='Flush DB',
    icon_name='delete'
)


# frontend is a next.js app which has built-in support for hot reload
# live_update only syncs changed files to the correct place for it to pick up
# https://docs.tilt.dev/api.html#api.docker_build
# https://docs.tilt.dev/live_update_reference.html
docker_build(
    "frontend",
    context="./frontend",
    live_update=[
        sync('./frontend', '/usr/src/app')
    ]
)

# the various go services share a base image to avoid re-downloading the same
# dependencies numerous times - `only` is used to prevent unnecessary rebuilds
# https://docs.tilt.dev/api.html#api.docker_build
docker_build(
    "pixeltilt-base",
    context=".",
    dockerfile="base.dockerfile",
    only=['go.mod', 'go.sum']
)

# docker_build_with_restart automatically restarts the process defined by
# `entrypoint` argument after completing the live_update (which syncs .go
# source files and recompiles inside the container)
# https://github.com/tilt-dev/tilt-extensions/tree/master/restart_process
# https://docs.tilt.dev/live_update_reference.html
docker_build_with_restart(
    "glitch",
    context=".",
    dockerfile="filters/glitch/Dockerfile",
    only=['filters/glitch', 'render/api'],
    entrypoint='/usr/local/bin/glitch',
    live_update=[
        sync('filters/glitch', '/app/glitch'),
        sync('render/api', '/app/render/api'),
        run('go build -mod=vendor -o /usr/local/bin/glitch ./glitch')
    ]
)

# for the remainder of the services, plain docker_build is used - these
# services are changed less frequently, so live_update is less important
# any of them can be adapted to use live_update by using "glitch" as an
# example above!
docker_build(
    "muxer",
    context=".",
    dockerfile="muxer/Dockerfile",
    only=['muxer', 'render/api', 'storage/api', 'storage/client']
)

docker_build(
    "color",
    context=".",
    dockerfile="filters/color/Dockerfile",
    only=['filters/color', 'render/api']
)

docker_build(
    "bounding-box",
    context=".",
    dockerfile="filters/bounding-box/Dockerfile",
    only=['filters/bounding-box', 'render/api']
)

docker_build(
    "storage",
    context=".",
    dockerfile="storage/Dockerfile",
    only=['storage'],
    entrypoint='/usr/local/bin/storage'
)
