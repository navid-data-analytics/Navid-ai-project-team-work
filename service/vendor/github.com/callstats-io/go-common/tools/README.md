### Development environment usage

1. Add the following entry to another project's `docker-compose.yaml`:
    ```yaml
    vaultbootstrapdev:
      build: .
      dockerfile: ./vendor/github.com/callstats-io/go-common/tools/Dockerfile
      volumes:
        - ./creds/:/output
        - ./vendor/:/go/src/github.com/callstats-io/go-common/tools/vendor
      environment:
        - ENV=dev
        - SERVICE_NAME=
        - VAULT_TOKEN=
        - VAULT_ADDR=http://vault:8200/
        - VAULT_MONGO_CLUSTER_NAME=devc1
        - VAULT_MONGO_ROOT_URL=mongodb://vault:vault@mongo:27017/admin?ssl=false
        - VAULT_CERT_NAME=x509
        - TLS_CERT_FILE=/go/src/github.com/callstats-io/go-common/tools/devcerts/cert.pem
        - TLS_CERT_KEY_FILE=/go/src/github.com/callstats-io/go-common/tools/devcerts/key.pem
      links:
        - mongo
        - vault
      command: '/go/bin/tools -bootstrap-vault-dev'
    ```
1. Match the
    - `SERVICE_NAME` with another project's name
    - `VAULT_TOKEN` with the `VAULT_DEV_ROOT_TOKEN_ID` you configure for vault image
    - auth credentials in `VAULT_MONGO_ROOT_URL` with what you have as `INIT_USERNAME` and `INIT_PASSWORD` for mongo
1. Run `docker-compose run vaultbootstrapdev`
