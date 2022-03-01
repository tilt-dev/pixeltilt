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


krsync_path = os.path.join(os.getcwd(), 'krsync.sh')

MIN_LOCAL_RSYNC_VERSION = '3.0.0'
verify_rsync_path = os.path.join(os.getcwd(), 'verify_rsync.sh')

# if no local rsync/insufficient version, bail
print('-- syncback extension checking for local rsync --')
local([verify_rsync_path, MIN_LOCAL_RSYNC_VERSION])

DEFAULT_EXCLUDES = ['.git', '.gitignore', '.dockerignore', 'Dockerfile', '.tiltignore', 'Tiltfile', 'tilt_modules']

def syncback_command_args(k8s_object, src_dir, ignore=None, delete=False, paths=None, target_dir='.', container='', namespace='', verbose=False):
    """
    Generate a list of command arguments to run that will sync the specified files from the given k8s object to the local filesystem.
    :param k8s_object (str): a Kubernetes object identifier (e.g. deploy/my-deploy, job/my-job, or a pod ID) that Tilt
           can use to select a pod. As per the behavior of `kubectl exec`, we will act on the first pod of the specified
           object, using the first container by default.
    :param src_dir (str): directory IN THE KUBERNETES CONTAINER to sync from. Any paths specified, if relative,
           should be relative to this dir.
    :param ignore (List[str], optional): files to ignore when syncing (relative to src_dir).
    :param delete (bool, optional): run rsync with the --delete flag, i.e. delete files locally if not present in
           the container. By default, False. THIS OPTION RISKS WIPING OUT FILES that exist locally but not in the
           container. Tilt will protect some files automatically, but we recommend syncing specific paths (via `paths`
           and/or using the `ignore` parameter to explicitly protect other files that exist locally but not on the container.
    :param paths (List[str], optional): paths IN THE KUBERNETES CONTAINER to sync, relative to src_dir. May be files or dirs.
           Note that these must not begin with `./`. If this arg is not passed, sync all of src_dir.
    :param target_dir (str, optional): directory ON THE LOCAL FS to sync to. Defaults to '.'
    :param container (str, optiona): name of the container to sync from (by default, the first container)
    :param namespace (str, optiona): namespace of the desired k8s_object, if not `default`.
    :param verbose (bool, optional): if true, print additional rsync information.
    """
    # Verify inputs
    if not src_dir.endswith('/'):
        fail('src_dir must be a directory and have a trailing slash (because of rsync syntax rules)')

    if paths:
        for p in paths:
            if p.startswith('./'):
                fail('Found illegal path "{}": paths may not begin with ./ (because of rsync syntax rules)'.format(p))
            if p.startswith('/'):
                fail('Found illegal path "{}": paths may not begin with / and must be relative to src_dir (because of rsync syntax rules)'.format(p))

    to_include = []
    to_exclude = ignore
    if not ignore:
        to_exclude = []

    if paths:
        # TODO: if you're rsync-savvy you might want to do the wildcarding manually--
        #   give an option to turn off automatic +'***'
        to_include = ['--include={}***'.format(p) for p in paths]
    else:
        # Sync the entire src_dir. Danger, Will Robinson! Exclude some stuff
        # that may exist locally but not in your container so it
        # doesn't get wiped out locally on your first sync
        to_exclude = DEFAULT_EXCLUDES + to_exclude

    to_exclude = ['--exclude={}***'.format(ex) for ex in to_exclude]

    # set remote name to a dummy name
    remote_name = 'syncback'

    # instead of wrestling with passing optional args to krsync.sh that do not
    # then get passed to rsync, just bundle container and namespace flags with
    # k8s object specifier (quoted 1st argument)
    if container:
        k8s_object = '{obj} -c {container}'.format(obj=k8s_object, container=container)

    if namespace:
        k8s_object = '{obj} -n {namespace}'.format(obj=k8s_object, namespace=namespace)

    flags = '-aOv'
    if verbose:
        flags = '-aOvvi'

    argv = [
        krsync_path,
        k8s_object,
        flags,
        '--progress',
        '--stats',
    ]
    if delete:
        argv.append('--delete')
    argv.append('-T=/tmp/rsync.tilt')
    argv.extend(to_include)
    argv.extend(to_exclude)
    argv.append(remote_name + ':' + src_dir)
    argv.append(target_dir)
    return argv


def syncback_command(k8s_object, src_dir, ignore=None, delete=False, paths=None, target_dir='.', container='', namespace='', verbose=False):
    """
    Generate a properly-quoted shell command string that will sync the specified files from the given k8s object to the local filesystem.
    """
    argv = syncback_command_args(k8s_object, src_dir, ignore, delete, paths, target_dir, container, namespace, verbose)
    return ' '.join([shlex.quote(arg) for arg in argv])


def syncback(name, k8s_object, src_dir, ignore=None, delete=False, paths=None, target_dir='.', container='', namespace='', verbose=False, labels=[], resource_deps=[]):
    """
    Create a local resource that will (via rsync) sync the specified files
    from the specified k8s object to the local filesystem.

    :param name (str): name of the created local resource.
    :param k8s_object (str): a Kubernetes object identifier (e.g. deploy/my-deploy, job/my-job, or a pod ID) that Tilt
           can use to select a pod. As per the behavior of `kubectl exec`, we will act on the first pod of the specified
           object, using the first container by default.
    :param src_dir (str): directory IN THE KUBERNETES CONTAINER to sync from. Any paths specified, if relative,
           should be relative to this dir.
    :param ignore (List[str], optional): files to ignore when syncing (relative to src_dir).
    :param delete (bool, optional): run rsync with the --delete flag, i.e. delete files locally if not present in
           the container. By default, False. THIS OPTION RISKS WIPING OUT FILES that exist locally but not in the
           container. Tilt will protect some files automatically, but we recommend syncing specific paths (via `paths`
           and/or using the `ignore` parameter to explicitly protect other files that exist locally but not on the container.
    :param paths (List[str], optional): paths IN THE KUBERNETES CONTAINER to sync, relative to src_dir. May be files or dirs.
           Note that these must not begin with `./`. If this arg is not passed, sync all of src_dir.
    :param target_dir (str, optional): directory ON THE LOCAL FS to sync to. Defaults to '.'
    :param container (str, optiona): name of the container to sync from (by default, the first container)
    :param namespace (str, optiona): namespace of the desired k8s_object, if not `default`.
    :param verbose (bool, optional): if true, print additional rsync information.
    :param labels (Union[str, List[str]], optional): Used to group resources in the Web UI.
    :param resource_deps (Union[str, List[str]], optional): Used to declare dependencies on other resources.
    """
    # Ensure extra kwargs are only passed to local_resource if specified to
    # provide backward compatibility with older version of tilt.

    command = syncback_command(k8s_object, src_dir,
                               ignore=ignore, delete=delete, paths=paths,
                               target_dir=target_dir,
                               container=container, namespace=namespace,
                               verbose=verbose)
    extra_args = {}

    if labels:
        extra_args["labels"] = labels

    if resource_deps:
        extra_args["resource_deps"] = resource_deps

    local_resource(name, command, trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False, **extra_args)


syncback('syncback-js', 'deploy/frontend',
        '/usr/src/app/', ignore=['node_modules/'],
        target_dir='./frontend-copy',
)