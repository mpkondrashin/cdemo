package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	fileSecurity "github.com/trendmicro/tm-v1-fs-golang-sdk"
)

func DetectRegion(apiKey string) (result string, returnError error) {
	var wg sync.WaitGroup
	for _, domain := range fileSecurity.V1Regions {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			_, err := fileSecurity.NewClient(apiKey, d)
			if err != nil {
				returnError = err
				return
			}
			result = d
		}(domain)
	}
	wg.Wait()
	if result != "" {
		returnError = nil
	}
	return
}

var (
	apiKey    = flag.String("apikey", "", "Vision One API key")
	address   = flag.String("address", "", "File Security address")
	ignoreTLS = flag.Bool("tls", true, "Ignore TLS errros")
	cert      = flag.String("cert", "", "Certificate filepath")
	pml       = flag.Bool("pml", false, "Use Predictive Machine Learning")
)

func ScanFile(client *fileSecurity.AmaasClient, filePath string) {
	log.Printf("Scan File: %s", filePath)
	var data []byte
	var err error
	if filePath == "eicar.com" {
		data = []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Readinf file: %v", err)
		}
	}
	fileName := filepath.Base(filePath)
	response, err := client.ScanBuffer(data, fileName, nil)
	if err != nil {
		log.Fatalf("ScanBuffer: %v", err)
	}
	var result fileSecurity.ScanResult2Client
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Printf("Filename: %s\n", result.FileName)
	log.Printf("Scan result: %v\n", result.ScanResult)
	for _, found := range result.FoundMalwares {
		log.Printf("Malware: %s (%s)\n", found.MalwareName, found.Engine)
	}
}

func main() {
	flag.Parse()
	fileSecurity.SetLoggingLevel(fileSecurity.LogLevelDebug)
	client, err := fileSecurity.NewClientInternal(*apiKey, *address, *ignoreTLS, *cert)
	if err != nil {
		log.Fatalf("NewClientInternal: %v", err)
	}
	defer client.Destroy()
	if *pml {
		client.SetPMLEnable()
		log.Println("PML enabled")
	} else {
		log.Println("PML disabled")
	}
	log.Printf("Scanning %d files", len(flag.Args()))
	for _, filePath := range flag.Args() {
		log.Printf("Scanning file: %s", filePath)
		ScanFile(client, filePath)
	}
	log.Println("Scan complete")
	log.Printf("Finished scanning %d files", len(flag.Args()))
}
