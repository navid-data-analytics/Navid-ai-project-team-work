package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	bootstrapVaultFlagSet = flag.NewFlagSet("bootstrap-vault", flag.ExitOnError)
	credsPath             = bootstrapVaultFlagSet.String("with-creds-path", "/output/vault_dev_creds.json", "Path where to save the vault bootstrap credentials")
	bootstrapMongo        = bootstrapVaultFlagSet.Bool("with-mongo", false, "Include to bootstrap dev mongo db")
	bootstrapPostgres     = bootstrapVaultFlagSet.Bool("with-postgres", false, "Include to bootstrap dev postgres db")
	bootstrapTLSCert      = bootstrapVaultFlagSet.Bool("with-tls", false, "Include to bootstrap dev tls certs")
	bootstrapAWS          = bootstrapVaultFlagSet.Bool("with-aws", false, "Include to bootstrap dev aws")
	bootstrapGeneric      = bootstrapVaultFlagSet.Bool("with-generic", false, "Include to bootstrap dev generic")
	bootstrapTransit      = bootstrapVaultFlagSet.Bool("with-transit", false, "Include to bootstrap dev transit")
)

func main() {
	switch os.Args[1] {
	case "bootstrap-vault":
		bootstrapVaultFlagSet.Parse(os.Args[2:])

		if bootstrapVaultFlagSet.Parsed() {
			fmt.Println("Bootstrapping vault..")
			bootstrapVaultForDev(*credsPath, *bootstrapMongo, *bootstrapPostgres, *bootstrapTLSCert, *bootstrapAWS, *bootstrapGeneric, *bootstrapTransit)
			return
		}
	}
	bootstrapVaultFlagSet.PrintDefaults()
}
