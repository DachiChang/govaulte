package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

var (
	vaultADDR     string
	vaultKVSecret string
	vaultRoleID   string
	vaultSecretID string
)

func main() {
	args := os.Args[1:]
	vaultADDR = os.Getenv("VAULT_ADDR")
	vaultKVSecret = os.Getenv("VAULT_KV_SECRET")
	vaultRoleID = os.Getenv("VAULT_ROLE_ID")
	vaultSecretID = os.Getenv("VAULT_SECRET_ID")

	if vaultADDR == "" || vaultKVSecret == "" || vaultRoleID == "" || vaultSecretID == "" {
		fmt.Println("Missing VAULT_ADDR or VAULT_KV_SECRET or VAULT_ROLE_ID or VAULT_SECRET_ID in environment.")
		os.Exit(1)
	}
	if len(args) == 0 {
		fmt.Println("Not provide the execute program.")
		os.Exit(1)
	}

	config := vault.DefaultConfig()
	config.Address = vaultADDR

	client, err := vault.NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	approleSecret := &auth.SecretID{
		FromString: vaultSecretID,
	}
	approleAuth, err := auth.NewAppRoleAuth(
		vaultRoleID,
		approleSecret,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	authInfo, err := client.Auth().Login(context.Background(), approleAuth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if authInfo == nil {
		fmt.Println("appRole login failed, can't get secret token")
		os.Exit(1)
	}
	// already login, start to retrive secret
	secret, err := client.Logical().Read("kv/data/" + vaultKVSecret)
	if err != nil {
		fmt.Println(err)
	}
	secretData := (secret.Data["data"]).(map[string]interface{})
	for k, v := range secretData {
		os.Setenv(k, v.(string))
	}

	// syscall exec program from arguments
	env := os.Environ()
	execErr := syscall.Exec(args[0], args, env)
	if execErr != nil {
		panic(execErr)
	}
}
