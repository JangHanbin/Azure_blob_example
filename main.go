package main

import (
	"azure/auth"
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
	"strings"
)

//
////// this code will be packaged. so, Public function needed
//func CreateContainer(client *armresources.Client) {
//	url := "https://storageblobsgo.blob.core.windows.net/"
//	ctx := context.Background()
//
//	// Create a default Azure credential
//	credential, err := azidentity.NewDefaultAzureCredential(nil)
//	if err != nil {
//		log.Fatal("Invalid credentials with error: " + err.Error())
//	}
//
//	serviceClient, err := azblobb.NewServiceClient(url, credential, nil)
//
//	if err != nil {
//		log.Fatal("Invalid credentials with error: " + err.Error())
//	}
//
//	// Create the container
//	containerName := fmt.Sprintf("quickstart-%s", "d0rk")
//	fmt.Printf("Creating a container named %s\n", containerName)
//
//	containerClient, err2 := serviceClient.NewContainerClient(containerName)
//	_, err = containerClient.Create(ctx, nil)
//
//	if err != nil {
//		println(err2)
//		log.Fatal(err)
//	}
//}

//func GetBlobList() {
//	// List the blobs in the container
//	pager := containerClient.ListBlobsFlat(nil)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
//			fmt.Println(*v.Name)
//		}
//	}
//
//	if err = pager.Err(); err != nil {
//		log.Fatalf("Failure to list blobs: %+v", err)
//	}
//
//	// Download the blob
//	get, err := blobClient.Download(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func DownloadFile() {
//	// Download the blob
//	get, err := blobClient.Download(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	downloadedData := &bytes.Buffer{}
//	reader := get.Body(azblobb.RetryReaderOptions{})
//	_, err = downloadedData.ReadFrom(reader)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = reader.Close()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(downloadedData.String())
//}
//
//func CleanUp() {
//	// Cleaning up the quick start by deleting the blob and container
//	fmt.Printf("Press enter key to delete the blob fils, example container, and exit the application.\n")
//	bufio.NewReader(os.Stdin).ReadBytes('\n')
//	fmt.Printf("Cleaning up.\n")
//
//	// Delete the blob
//	fmt.Printf("Deleting the blob " + blobName + "\n")
//
//	_, err = blobClient.Delete(ctx, nil)
//	if err != nil {
//		log.Fatalf("Failure: %+v", err)
//	}
//
//	// Delete the container
//	fmt.Printf("Deleting the blob " + containerName + "\n")
//	_, err = containerClient.Delete(ctx, nil)
//
//	if err != nil {
//		log.Fatalf("Failure: %+v", err)
//	}
//}

func main() {

	credential, accountName := auth.GetCredentialFromFile("config.json")

	// Create a request pipeline that is used to process HTTP(S) requests and responses. It requires
	// your account credentials. In more advanced scenarios, you can configure telemetry, retry policies,
	// logging, and other options. Also, you can configure multiple request pipelines for different scenarios.
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// The URL typically looks like this:
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline.
	serviceURL := azblob.NewServiceURL(*u, p)

	// Now, you can use the serviceURL to perform various container and blob operations.

	// All HTTP operations allow you to specify a Go context.Context object to control cancellation/timeout.
	ctx := context.Background() // This example uses a never-expiring context.

	// This example shows several common operations just to get you started.

	// Create a URL that references a to-be-created container in your Azure Storage account.
	// This returns a ContainerURL object that wraps the container's URL and a request pipeline (inherited from serviceURL)
	containerURL := serviceURL.NewContainerURL("container-public-access") // Container names require lowercase

	// Create the container on the service (with no metadata and no public access)
	_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessBlob)
	if err != nil {
		println("CREATE ERROR")
		log.Fatal(err)

	}

	// Create a URL that references a to-be-created blob in your Azure Storage account's container.
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := containerURL.NewBlockBlobURL("directory/HelloWorld2.txt") // Blob names can be mixed case

	// Create the blob with string (plain text) content.
	data := "Hello World! play with me!"
	_, err = blobURL.Upload(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Download the blob's contents and verify that it worked correctly
	get, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(azblob.RetryReaderOptions{})
	downloadedData.ReadFrom(reader)
	reader.Close() // The client must close the response body when finished with it
	if data != downloadedData.String() {
		log.Fatal("downloaded data doesn't match uploaded data")
	}

	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatal(err)
		}
		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("Blob name: " + blobInfo.Name + "\n")
		}
	}

	//// Delete the blob we created earlier.
	//_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Delete the container we created earlier.
	//_, err = containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	//if err != nil {
	//	log.Fatal(err)
	//}

}
