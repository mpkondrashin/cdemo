package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	apiKey := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJjaWQiOiJiZTg0OTU0NS1lNjc0LTQwZjAtOTlkYy1mYjU2NWYzMjQ3NjAiLCJjcGlkIjoic3ZwIiwicHBpZCI6ImN1cyIsIml0IjoxNzQxMTcxNjc4LCJldCI6MTc3MjcwNzY3NywiaWQiOiJhNzdlYTNiZS1jMGFkLTRhNTUtOGRiNS04MDQzNTYzNTVkNDMiLCJ0b2tlblVzZSI6ImN1c3RvbWVyIn0.lfz2exUL1rhUY85y9JxaqqCQ9y3LRwIoxixkrB35gnotFFOU8Q4k1Qa3-_zSi87U2hRO7U9Vv99BkIEov6iAsEUWRBZU70iW-QbhrN8ziea4OEtZXwIhJUcZIayOnkpw1Bt3caT5yhUkXtlY4_0r-gQ1yqPv7L-nTk-LC8M2-R6V6O-nAVxTg_AEsgQiiY7hV6d-MoSbQ99AwQN47Gft1WZ_49YbpAGBVyRB4o31e8KbWRjdVonfoSOmJcORI9dsIk0wsUbb6KuTU4e4YaJ0HgbjqQS73GSXl-jtNn6bF0gPEFWNbL0CMHsgZ9Xcb-Mh09qAMrPa03LLehxLycSZvbdQbRV4NPkrZ5COuWh1EcdBRyJwSs1to2pb1owMmPwKbfeWJP81YD79LkVI0qEDRExctrJcDUWaIEk8Jr0RTMnILUmlLy0ZzX-BmWd7weoMirfz5NTQOr-qf9Am2Qdb4WAlbVY3sIEjU2wH8LnhmgALt95ZmQMCMMaVmupCys8nHXvxiZ9-BSaCrX9Wu8tMlBlGs61ZWcRrsCqM5y8BdpwJ4LyZ4Y4e9UA-VQCQR3onZyLM9MkwNwHUhItfe_LEPlwvqVLrBAOoAu1tHdaoQLktDqCo3TBTntitvvFumuoXYfXP3mUeI3QSxaaJ5vG-mPDwX2JXM55dIGfZkQpM5Gs"
	region := "eu-central-1"

	fileSecurity.SetLoggingLevel(fileSecurity.LogLevelDebug)

	fsec, err := fileSecurity.NewClient(apiKey, region)
	if err != nil {
		panic(err)
	}

	defer fsec.Destroy()
	fsec.SetPMLEnable()
	filePath := os.Args[1]
	var data []byte
	if filePath == "eicar.com" {
		data = []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
	}
	fileName := filepath.Base(filePath)
	response, err := fsec.ScanBuffer(data, fileName, nil)
	if err != nil {
		panic(err)
	}
	var result fileSecurity.ScanResult2Client
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Filename: %s\n", result.FileName)
	fmt.Printf("Scan result: %v\n", result.ScanResult)
	for _, found := range result.FoundMalwares {
		fmt.Printf("Malware: %s (%s)\n", found.MalwareName, found.Engine)
	}
}
