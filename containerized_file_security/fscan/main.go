// Copyright (c) 2026 Michael Kondrashin (mkondrashin@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	fileSecurity "github.com/trendmicro/tm-v1-fs-golang-sdk"
)

var (
	address = flag.String("address", "", "File Security address")
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
	log.Printf("Scan result: %v\n", result.ScanResult)
	for _, found := range result.FoundMalwares {
		log.Printf("Malware: %s (%s)\n", found.MalwareName, found.Engine)
	}
}

func main() {
	flag.Parse()
	fileSecurity.SetLoggingLevel(fileSecurity.LogLevelDebug)
	client, err := fileSecurity.NewClientInternal("", *address, true, "")
	if err != nil {
		log.Fatalf("NewClientInternal: %v", err)
	}
	defer client.Destroy()
	log.Printf("Scanning %d files", len(flag.Args()))
	for _, filePath := range flag.Args() {
		log.Printf("Scanning file: %s", filePath)
		ScanFile(client, filePath)
	}
	log.Println("Scan complete")
	log.Printf("Finished scanning %d files", len(flag.Args()))
}
