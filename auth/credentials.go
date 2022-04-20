package auth

import (
	"encoding/json"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"os"
)

type Configuration struct {
	AccountName string
	AccessKey   string
}

func GetCredentialFromFile(path string) (*azblob.SharedKeyCredential, string) { //may function name is too long
	// From the file, get your Storage account's name and account key.
	file, fileErr := os.Open(path)
	if fileErr != nil {
		log.Fatalf("Config file open failure: %+v", fileErr)
		panic(fileErr)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	configuration := Configuration{}

	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Config parsing failure: %+v", err)
		panic(err)
	}
	// Use your Storage account's name and key to create a credential object; this is used to access your account.

	credential, err := azblob.NewSharedKeyCredential(configuration.AccountName, configuration.AccessKey)
	return credential, configuration.AccountName

}
