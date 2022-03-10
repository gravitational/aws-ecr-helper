package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// PullSecretDataKey is name of secret data key
	PullSecretDataKey = ".dockerconfigjson"
)

func main() {
	cfg, err := ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ecrClient := ecr.NewFromConfig(cfg.AWSCfg)

	k8sClient, err := kubernetes.NewForConfig(cfg.KubeCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve credentials
	authTokenOutput, err := ecrClient.GetAuthorizationToken(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	authToken := *authTokenOutput.AuthorizationData[0].AuthorizationToken

	dockerConfig := generateDockerConfig(cfg.RegistryURI, authToken)
	fmt.Println(dockerConfig)

	ecrSecret := new(v1.Secret)
	ecrSecret.Name = cfg.SecretName

	ecrSecret.Data[PullSecretDataKey] = []byte(authToken)

	secret, err := k8sClient.CoreV1().Secrets(cfg.Namespace).Create(context.Background(), ecrSecret, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created secret: %s", secret.Name)
}

func generateDockerConfig(registry, auth string) string {
	authMap := make(map[string]DockerAuth)

	authMap[registry] = DockerAuth{
		Auth: auth,
	}

	dockerConfig := DockerConfig{Auths: authMap}

	configBytes, err := json.Marshal(dockerConfig)
	if err != nil {
		panic(err)
	}

	return string(configBytes)
}

type DockerConfig struct {
	Auths map[string]DockerAuth `json:"auths"`
}

type DockerAuth struct {
	Auth string `json:"auth"`
}
