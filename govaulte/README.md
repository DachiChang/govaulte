# govaulte

govaulte aims to use approle to acces secret which in hashicorp Vault and produce runtime environment for executing program to get secret what it wants.

## Build

- make build (for amd64)
- go build . (for your os version)

## Test your vaulte credential is worked

1. Export your credential to local environments. (see also http://github.com/dachichang/govaulte/vaulte-credential-generator)
  ```
  export $(xargs < vaulte_credential)
  ```
2. check govaulte get your secret in the right way.
  ```
  govaulte /bin/bash -l -c "export"
  ```
## Running your code with govaulte

- example
    - /usr/sbin/govaulte /path/to/your/app/bin/main -argv1=data1 -argv2=data2 command1
    - /usr/sbin/govaulte /bin/bash -l -c "cd /path/to/your/workdir/; ./bin/main"

## Reference

- Use [vaulte-credential-generator](https://github..com/dachichang/govaulte/vaulte-credential-generator) to easy generate your vaulte credential.
