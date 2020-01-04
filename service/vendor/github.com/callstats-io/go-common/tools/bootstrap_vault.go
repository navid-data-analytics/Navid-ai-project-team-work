package main

import (
	"fmt"

	"github.com/callstats-io/go-common/vaultbootstrap"
)

func bootstrapVaultForDev(credsFilename string, bootstrapMongo, bootstrapPostgres, bootstrapTLSCert, bootstrapAWS, bootstrapGeneric, bootstrapTransit bool) {
	client := vaultbootstrap.NewBootstrapClient().
		UnmountAll().
		MountAppRoleAuth().
		WriteCredentialsFile(credsFilename)

	if bootstrapMongo {
		fmt.Println("Boostrapping mongo..")
		client.MountMongo()
	}

	if bootstrapPostgres {
		fmt.Println("Boostrapping postgres..")
		client.MountPostgres()
	}

	if bootstrapTLSCert {
		fmt.Println("Boostrapping tls..")
		client.MountTLSCert()
	}

	if bootstrapAWS {
		fmt.Println("Boostrapping aws..")
		client.MountAWS()
	}

	if bootstrapGeneric {
		fmt.Println("Boostrapping generic..")
		client.MountGeneric()
	}

	if bootstrapTransit {
		fmt.Println("Boostrapping transit..")
		client.MountTransit()
	}

	fmt.Printf("Boostrapped successfully. Credentials written to %s\n", credsFilename)
}
