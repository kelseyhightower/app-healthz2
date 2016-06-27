package main

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/vault/api"
)

type vaultClient struct {
	client          *api.Client
	dbLeaseID       string
	dbLeaseDuration int
	dbRenewable     bool
}

func newVaultClient(addr, token string) (*vaultClient, error) {
	config := api.Config{Address: addr}
	client, err := api.NewClient(&config)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return &vaultClient{client: client}, nil
}

func (v *vaultClient) getDatabaseCredentials(path string) (string, string, error) {
	var err error
	var secret *api.Secret

	for {
		secret, err = v.client.Logical().Read(path)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	v.dbLeaseID = secret.LeaseID
	v.dbLeaseDuration = secret.LeaseDuration
	v.dbRenewable = secret.Renewable

	username := secret.Data["username"].(string)
	password := secret.Data["password"].(string)
	return username, password, nil
}

func (v *vaultClient) renewDatabaseCredentials() error {
	if !v.dbRenewable {
		return errors.New("credentials not renewable")
	}

	log.Println("Renewing credentials:", v.dbLeaseID)
	_, err := v.client.Sys().Renew(v.dbLeaseID, v.dbLeaseDuration)
	if err != nil {
		log.Println(err)
	}

	// Renew the lease before it expires.
	duration := (v.dbLeaseDuration - 300)

	for {
		time.Sleep(time.Second * time.Duration(duration))
		log.Println("Renewing credentials:", v.dbLeaseID)
		// Should we be reusing the secret?
		_, err := v.client.Sys().Renew(v.dbLeaseID, v.dbLeaseDuration)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
