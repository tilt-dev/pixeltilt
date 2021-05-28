module github.com/windmilleng/pixeltilt

go 1.16

require (
	github.com/fogleman/gg v1.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sug0/go-glitch v0.0.0-20190629024109-a11fbafffa96
	github.com/tilt-dev/tilt-api-client-go v0.0.2
	github.com/tilt-dev/tilt-apiserver v0.3.1
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	k8s.io/apimachinery v0.20.2
)

replace (
	github.com/pkg/browser v0.0.0-00010101000000-000000000000 => github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4

	k8s.io/apimachinery => github.com/tilt-dev/apimachinery v0.20.2-tilt-20210505
)
