package secrets

import (
	"errors"
	"os"
	"strings"
)

type SecretManager struct {
	SecretDirectory string
	Secrets         map[string]interface{}
	TypeConversion  bool
}

func NewSecretManager(secretDirectory string, typeConversion bool) (*SecretManager, error) {
	if secretDirectory == "" {
		secretDirectory = "/run/secrets"
	}

	sm := &SecretManager{
		SecretDirectory: secretDirectory,
		Secrets:         make(map[string]interface{}),
		TypeConversion:  typeConversion,
	}

	err := sm.Load()
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func (sm *SecretManager) Load() error {
	if !isDirectory(sm.SecretDirectory) {
		return errors.New("secret directory not exist or not reachable")
	}

	files, err := os.ReadDir(sm.SecretDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		secretName := file.Name()
		secretValue, err := os.ReadFile(sm.SecretDirectory + "/" + secretName)
		if err != nil {
			return err
		}
		sm.Secrets[secretName] = strings.Trim(string(secretValue), "\t\n\v\f\r")
	}

	return nil
}

func (sm *SecretManager) GetSecret(secretName string) interface{} {
	return sm.Secrets[secretName]
}

func (sm *SecretManager) SetSecret(secretName string, secretValue interface{}) {
	sm.Secrets[secretName] = secretValue
}

func (sm *SecretManager) ExtendSecretFromFile(secretName string, secretFilePath string) error {
	secretValue, err := os.ReadFile(secretFilePath)
	if err != nil {
		return err
	}
	sm.Secrets[secretName] = string(secretValue)

	return nil
}

func (sm *SecretManager) ExtendSecretsFromFiles(secretFiles map[string]string) error {
	for k, v := range secretFiles {
		err := sm.ExtendSecretFromFile(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func isDirectory(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
