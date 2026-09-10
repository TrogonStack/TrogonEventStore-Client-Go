package samples

import (
	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func UserCertificates() {
	// region client-with-user-certificates
	settings, err := trogoneventstore.ParseConnectionString("esdb://admin:changeit@{endpoint}?tls=true&userCertFile={pathToCaFile}&userKeyFile={pathToKeyFile}")

	if err != nil {
		panic(err)
	}

	db, err := trogoneventstore.NewClient(settings)
	// endregion client-with-user-certificates

	if err != nil {
		panic(err)
	}

	db.Close()
}
