package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

var (
	vaultToken    string
	vaultAddr     string
	kvName        string
	kvMount       string
	approleID     string
	approleSecret string
)

func init() {
	flag.StringVar(&vaultToken, "token", "", "vault token")
	flag.StringVar(&vaultToken, "t", "", "vault token")
	flag.StringVar(&vaultAddr, "address", "http://127.0.0.1:8200", "vault server address ex: https://vault.test-dev.com")
	flag.StringVar(&vaultAddr, "a", "http://127.0.0.1:8200", "vault server address ex: https://vault.test-dev.com")
	flag.StringVar(&kvName, "name", "", "key value entity name")
	flag.StringVar(&kvName, "n", "", "key vaule entity name")
	flag.StringVar(&kvMount, "mount", "kv", "kv mount engine")
	flag.StringVar(&kvMount, "m", "kv", "kv mount engine")
}

func envToMap(stdin *os.File) map[string]interface{} {
	envMaps := make(map[string]interface{})
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		kv := strings.Split(scanner.Text(), "=")
		if len(kv) == 2 {
			envMaps[kv[0]] = kv[1]
		}
	}

	return envMaps
}

func main() {
	flag.Parse() // get token from CLI

	kvMap := envToMap(os.Stdin)
	fmt.Println(strings.Repeat("-", 20), " KV Start", strings.Repeat("-", 20))
	for k, v := range kvMap {
		fmt.Printf("%s=%s\n", k, v)
	}
	fmt.Println(strings.Repeat("-", 20), " KV End ", strings.Repeat("-", 20))

	// create vault client
	config := vault.DefaultConfig()
	config.Address = vaultAddr

	client, err := vault.NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client.SetToken(vaultToken)

	// create or update KV2 with kvName
	_, err = client.KVv2(kvMount).Put(context.Background(), kvName, kvMap)
	if err != nil {
		// TODO because vault server version is too old, so will get error about conflict response
		// only show the error string, but not terminate.
		fmt.Println(err)
	} else {
		fmt.Printf("KV entity '%s' inject done...\n", kvName)
	}

	// create policy for the kv entity
	sys := client.Sys()
	policyName := fmt.Sprintf("policy_%s", kvName)
	policyRule := fmt.Sprintf(`
path "%s/data/%s" {
	capabilities = ["read", "list"]
}
path "%s/data/%s/*" {
	capabilities = ["read", "list"]
}
	`, kvMount, kvName, kvMount, kvName)
	err = sys.PutPolicy(policyName, policyRule)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("policy '%s' set down...\n", policyName)
	}

	// create an approle
	approlePath := fmt.Sprintf("auth/approle/role/%s", kvName)
	_, err = client.Logical().Write(
		approlePath,
		map[string]interface{}{
			"token_policies": policyName,
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("approle '%s' set down...\n", kvName)
	}

	// get approle role-id
	roleIDResult, err := client.Logical().Read(fmt.Sprintf("%s/role-id", approlePath))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	approleID = (roleIDResult.Data["role_id"]).(string)

	// generate approle secret-id
	secretResult, err := client.Logical().Write(fmt.Sprintf("%s/secret-id", approlePath), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	approleSecret = (secretResult.Data["secret_id"]).(string)
	// prepare env for vaulte.pl
	vaulteEnvAddr := fmt.Sprintf("VAULT_ADDR=%s", vaultAddr)
	vaulteEnvKVSecret := fmt.Sprintf("VAULT_KV_SECRET=%s", kvName)
	vaulteEnvRoleID := fmt.Sprintf("VAULT_ROLE_ID=%s", approleID)
	vaulteEnvSecretID := fmt.Sprintf("VAULT_SECRET_ID=%s", approleSecret)
	vaulteInfo := []string{
		vaulteEnvAddr,
		vaulteEnvKVSecret,
		vaulteEnvRoleID,
		vaulteEnvSecretID,
	}

	dumpFileName := fmt.Sprintf("environment.%s", kvName)
	f, err := os.Create(dumpFileName)
	defer f.Close()

	if err != nil {
		fmt.Printf("Can' open file %s, print to stdout\n", dumpFileName)
		fmt.Println(strings.Repeat("-", 20), " here is information for Vaulte.pl ", strings.Repeat("-", 20))
		// just print information to stdout
		for _, v := range vaulteInfo {
			fmt.Println(v)
		}
		return
	}

	for _, v := range vaulteInfo {
		f.WriteString(v + "\n")
	}
	fmt.Printf("'%s' is created...\n", dumpFileName)
}
