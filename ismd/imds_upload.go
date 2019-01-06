package main

import (
	//"bufio"
	//"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"path/filepath"
)

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func handleErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				fmt.Println("Received 409. Container already exists")
				return
			}
		}
		log.Fatal(err)
	}
}
func listFiles(rootPath, searchPath string, modifiedSpan int) ([]string, []string ){
	fullpathlist := []string{}
	filelist := []string{}
	old := time.Now().Add(time.Duration(-modifiedSpan) * time.Hour)
	fis, err := ioutil.ReadDir(searchPath)
	
    if err != nil {
        panic(err)
    }

    for _, fi := range fis {

		if ( old.Before(fi.ModTime()) ){break}
		
		fullPath := filepath.Join(searchPath, fi.Name())
		fmt.Println(fullPath)
        if fi.IsDir() {
			subf,sub := listFiles(rootPath, fullPath, modifiedSpan)
			fullpathlist = append(fullpathlist, subf...)
			filelist = append(filelist, sub...)

		} else {
			fmt.Printf("Time:%v\n", fi.ModTime())
			fullpathlist = append(fullpathlist, fullPath)
			rel, err := filepath.Rel(rootPath, fullPath)

            if err != nil {
                panic(err)
			}
			filelist = append(filelist, rel)
		}

	}
	return fullpathlist, filelist
}

func main() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	// From the Azure portal, get your storage account name and key and set environment variables.
	//accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	accountName, accountKey := "xxx", "xxx=="
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a container name
	containerName := "ismd-logs"

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create the container
	fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background() // This example uses a never-expiring context
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	handleErrors(err)

	// Create a file to test the upload and download.
	//fmt.Printf("Creating a dummy file to test the upload and download\n")
	/*data := []byte("hello world this is a blob\n")
	fileName := randomString()
	err = ioutil.WriteFile(fileName, data, 0700)
	handleErrors(err)*/

	// Here's how to upload a blob.
	//fileName := randomString()
	root := "/home/azureuser/go/src/ismd/"
	logs := root + "logs"
	archive := root + "archive"
	fullpathlist, relpathlist := listFiles(logs, logs, 3)

	for fi := range fullpathlist {

		blobURL := containerURL.NewBlockBlobURL(relpathlist[fi])
		file, err := os.Open(fullpathlist[fi])
		handleErrors(err)
		
		fmt.Printf("Uploading the file with blob name: %s\n", relpathlist[fi])
		_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16})
		handleErrors(err)

		if err := os.Rename(fullpathlist[fi], archive+"/"+relpathlist[fi]+".old"); err != nil {
			fmt.Println(err)
		}

		file.Close()
	}
}
