package config

import (
	"errors"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	"github.com/tilt-dev/wmclient/pkg/dirs"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func NewConfig() (*rest.Config, error) {
	home := homedir.HomeDir()
	if home == "" {
		return nil, errors.New("no homedir could be found")
	}

	dir, err := dirs.GetTiltDevDir()
	if err != nil {
		return nil, err
	}
	p := filepath.Join(dir, "config")

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: p},
		&clientcmd.ConfigOverrides{CurrentContext: "tilt-default"}).ClientConfig()

	if err != nil {
		return nil, err
	}

	scheme := v1alpha1.NewScheme()
	cfg.ContentConfig.GroupVersion = &schema.GroupVersion{
		Group: v1alpha1.GroupName,
		Version: v1alpha1.Version,
	}
	cfg.APIPath = "/apis"
	cfg.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
	cfg.UserAgent = rest.DefaultKubernetesUserAgent()

	return cfg, nil
}
