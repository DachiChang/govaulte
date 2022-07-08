# Vaulte Credential Generator

This project aims to easy inject secret to vault and auto generate credential file for govaulte

It will be easy to host your secret in Vault and easy to use approle to grab the secret/env from Vault.

## Install

- go install .

## What will it does

1. parse your env file to k/v map
2. create vault kv2 secret
3. create a policy for the secret
4. create a approle and bind with the policy
5. get approle id and generate secret id
6. help you dump the information which vaulte.pl needed in the file.

## How to use

- prepare your secret file (or env with key/value) EX: secret.data
  ```
  MQ_HOST=rabbitmq.test-dev.com
  MQ_VHOST=
  MQ_PORT=5672
  MQ_USER=user1
  MQ_Password=12345
  MQ_EXCHANGE=test_exchange
  ```
- prepare
    - vault server address EX: https://vault.test-dev.com/
    - vault token (which you can login vault WebUI by LDAP authentication and copy **token** from profile) EX: xxxxxx
        - ![](https://i.imgur.com/TyB2Prr.png)
    - your secret name (Maybe it will follow up by the probject name) EX: test-secretdata
- check usage of this program
    - vaulte-credential-generator -h
- run program
    - cat secret.data | vaulte-credential-generator -a https://vault.test-dev.com -t xxxx -n test-secretdata
- you will have
    - environment.test-secretdata file in the folder.
    - environment.test-secretdata file is used by govaulte
