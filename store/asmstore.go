package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

const (
	// DefaultASMKeyID is the default alias for the KMS key used to encrypt/decrypt secrets
	DefaultASMKeyID = "alias/secrets_manager_key"
)

// ensure ASMStore confirms to Store interface
var _ Store = &ASMStore{}

// ASMStore implements the Store interface for storing secrets in AWS Secrets Manager
type ASMStore struct {
	svc secretsmanageriface.SecretsManagerAPI
}

// KMSKey returns the key, prepending alias/ if necessary
func (s *ASMStore) KMSKey() string {
	fromEnv, ok := os.LookupEnv("CHAMBER_KMS_KEY_ALIAS")
	if !ok {
		return DefaultASMKeyID
	}
	if !strings.HasPrefix(fromEnv, "alias/") {
		return fmt.Sprintf("alias/%s", fromEnv)
	}

	return fromEnv
}

// NewASMStore creates a new ASMStore
func NewASMStore(numRetries int) *ASMStore {
	asmSession, region := getSession(numRetries)

	svc := secretsmanager.New(asmSession, &aws.Config{
		MaxRetries: aws.Int(numRetries),
		Region:     region,
	})

	return &ASMStore{
		svc: svc,
	}
}

// Write either creates a new secret or a new version of an existing secret
func (s *ASMStore) Write(id SecretId, value string) error {
	fmt.Println("Write stub method")

	return nil
}

// Read returns a secret
func (s *ASMStore) Read(id SecretId, version int) (Secret, error) {
	if version == -1 {
		return s.readLatest(id)
	}

	return s.readVersion(id, version)
}

func (s *ASMStore) readLatest(id SecretId) (Secret, error) {
	describeSecretInput := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(s.idToName(id)),
	}
	desc, err := s.svc.DescribeSecret(describeSecretInput)
	if err != nil {
		return Secret{}, ErrSecretNotFound
	}

	fmt.Printf("%v", desc.VersionIdsToStages)

	getSecretValueInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.idToName(id)),
	}

	resp, err := s.svc.GetSecretValue(getSecretValueInput)
	if err != nil {
		return Secret{}, ErrSecretNotFound
	}

	secretMeta := SecretMetadata{}
	/*
			Created:   *resp.CreatedDate,
			CreatedBy: aws.String("N/A"),
			Version:   *desc,
			Key:       *p.Name,
		}
	*/
	return Secret{
		Value: resp.Name,
		Meta:  secretMeta,
	}, nil
}

func (s *ASMStore) readVersion(id SecretId, version int) (Secret, error) {
	describeSecretInput := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(s.idToName(id)),
	}

	resp, err := s.svc.DescribeSecret(describeSecretInput)
	fmt.Printf("%v", resp)
	if err != nil {
		return Secret{}, ErrSecretNotFound
	}
	/*
		for _, history := range resp.VersionIdsToStages {
			return Secret{}, nil
			thisVersion := 0
			if history.Description != nil {
				thisVersion, _ = strconv.Atoi(*history.Description)
			}
			if thisVersion == version {
				return Secret{
					Value: history.Value,
					Meta: SecretMetadata{
						Created:   *history.LastModifiedDate,
						CreatedBy: *history.LastModifiedUser,
						Version:   thisVersion,
						Key:       *history.Name,
					},
				}, nil
			}
		}
	*/
	return Secret{}, nil
}

func (s *ASMStore) idToName(id SecretId) string {
	return fmt.Sprintf("%s/%s", id.Service, id.Key)
}

func (s *ASMStore) validateName(name string) bool {
	return validPathKeyFormat.MatchString(name)
}

// List returns the key and metadata of all secrets
func (s *ASMStore) List(service string, includeValues bool) ([]Secret, error) {
	fmt.Println("List stub method")

	return []Secret{}, nil
}

// ListRaw returns raw list data
func (s *ASMStore) ListRaw(service string) ([]RawSecret, error) {
	fmt.Println("ListRaw stub method")

	return []RawSecret{}, nil
}

// History returns the history of a secret
func (s *ASMStore) History(id SecretId) ([]ChangeEvent, error) {
	fmt.Println("History stub method")

	return []ChangeEvent{}, nil
}

// Delete removes a secret
func (s *ASMStore) Delete(id SecretId) error {
	fmt.Println("Delete stub method")

	return nil
}

// Rotate rotates the secret data
func (s *ASMStore) Rotate(id SecretId) error {
	fmt.Println("Rotate stub method")

	return nil
}
