package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
			fs, err := fileSecurity.NewClient(apiKey, d)
			if err != nil {
				returnError = err
				return
			}
			_, err = fs.ScanBuffer([]byte{}, "detection", nil)
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
	flagAPIKey   = flag.String("api-key", "", "API key")
	flagRegion   = flag.String("region", "", "Region")
	flagFilename = flag.String("filename", "", "Filename")
	flagFolder   = flag.String("folder", "", "Folder")
)

func ScanFile(fsec *fileSecurity.AmaasClient, filePath string) (err error) {
	var data []byte
	if filePath == "eicar.com" {
		data = []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return
		}
	}
	fileName := filepath.Base(filePath)
	response, err := fsec.ScanBuffer(data, fileName, nil)
	if err != nil {
		return
	}
	var result fileSecurity.ScanResult2Client
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		return
	}
	fmt.Printf("Filename: %s\n", result.FileName)
	fmt.Printf("Scan result: %v\n", result.ScanResult)
	for _, found := range result.FoundMalwares {
		fmt.Printf("Malware: %s (%s)\n", found.MalwareName, found.Engine)
	}
	return
}

func ScanFolder(fsec *fileSecurity.AmaasClient, folderPath string) (err error) {
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return ScanFile(fsec, path)
	})
	return
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("FScan Starting...")
	flag.Parse()
	fileSecurity.SetLoggingLevel(fileSecurity.LogLevelDebug)
	apiKey := *flagAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("TMFS_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("API key not specified")
	}
	region := *flagRegion
	if region == "" {
		log.Println("Detecting region...")
		var err error
		region, err = DetectRegion(apiKey)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Region: " + region)
	}
	fsec, err := fileSecurity.NewClient(apiKey, region)
	if err != nil {
		log.Fatal(err)
	}

	defer fsec.Destroy()
	fsec.SetPMLEnable()
	if *flagFilename != "" {
		err := ScanFile(fsec, *flagFilename)
		if err != nil {
			log.Fatal(err)
		}
	} else if *flagFolder != "" {
		err := ScanFolder(fsec, *flagFolder)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Filename or folder must be specified")
	}
}
