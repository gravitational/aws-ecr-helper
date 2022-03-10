package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"k8s.io/client-go/rest"
)

const (
	// REGISTRY_URI is environment name for variable
	REGISTRY_URI = "REGISTRY_URI"

	// NAMESPACE is environment name for variable
	NAMESPACE = "K8S_NAMESPACE"

	// SECRET_NAME is environment name for variable
	SECRET_NAME = "PULL_SECRET_NAME"

	// SECRET_NAME_DEFAULT is default value for SECRET_NAME
	SECRET_NAME_DEFAULT = "ecr-pull-secret"
)

type Config struct {
	// SecretName is the name of the secret created by this helper
	// This value defaults to `ecr-pull-secret` but can be overwridden by
	// the environment variable: PULL_SECRET_NAME
	SecretName string

	// RegistryURI is location of AWS ECR Registry.
	RegistryURI string

	// Namespace is kubernetes namespace aws-ecr-helper is running in.
	// This is typically set within the pod spec.
	Namespace string

	// AWSConfig is used for generating SDK Clients
	AWSCfg aws.Config

	// KubeConfig is kubernetes config for kubernetes clientset
	KubeCfg *rest.Config
}

// ConfigFromEnv populates Config with corresponding environment variables
func ConfigFromEnv() (*Config, error) {

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	k8sCfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	registryURI := os.Getenv(REGISTRY_URI)
	if registryURI == "" {
		return nil, fmt.Errorf("%s is unset or contains an empty value", REGISTRY_URI)
	}

	namespace := os.Getenv(NAMESPACE)
	if namespace == "" {
		return nil, fmt.Errorf("%s is unset or contains an empty value", NAMESPACE)
	}

	secretName := os.Getenv(SECRET_NAME)
	if secretName == "" {
		secretName = SECRET_NAME_DEFAULT
	}

	return &Config{
			RegistryURI: registryURI,
			SecretName:  secretName,
			Namespace:   namespace,
			AWSCfg:      awsCfg,
			KubeCfg:     k8sCfg,
		},
		nil
}
