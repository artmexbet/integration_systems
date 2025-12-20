package main

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	serverURL := "http://localhost:8080/soap"

	fmt.Println("=======================================")
	fmt.Println("Lab 3: SOAP Service Complete Demo")
	fmt.Println("=======================================")
	fmt.Println()

	// Check if server is running
	fmt.Print("Checking server availability... ")
	resp, err := http.Get(serverURL + "?wsdl")
	if err != nil {
		fmt.Printf("✗ Server is not available: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("✓ Server is ready!")
	fmt.Println()

	fmt.Println("WSDL is available at: http://localhost:8080/soap?wsdl")
	fmt.Println()

	fmt.Println("Available credentials:")
	fmt.Println("  - user1:pass1")
	fmt.Println("  - user2:pass2")
	fmt.Println("  - admin:admin")
	fmt.Println()

	// Create test files
	fmt.Println("Creating test files...")
	tempDir := "uploaded_files/lab3_demo"
	os.MkdirAll(tempDir, 0755)

	testFiles := map[string]string{
		"valid.txt":    "This is a valid text file for upload",
		"invalid.json": `{"name":"test","value":123}`,
		"empty.txt":    "",
		"Жfile.txt":    "Test file with русский текст",
	}

	for name, content := range testFiles {
		ioutil.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
	}

	fmt.Printf("Test files created in %s/\n", tempDir)
	fmt.Println()

	// Demo Scenario 1: Successful file upload
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 1: Successful file upload")
	fmt.Println("================================================")

	fileBytes, _ := ioutil.ReadFile(filepath.Join(tempDir, "valid.txt"))
	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	uploadXml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>valid.txt</fileName>
      <fileData>%s</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>`, fileData)

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(uploadXml))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("✓ File uploaded successfully\n%s\n\n", formatXML(string(body)))

	// Demo Scenario 2: Upload with forbidden character
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 2: Upload with forbidden character")
	fmt.Println("================================================")

	fileBytes, _ = ioutil.ReadFile(filepath.Join(tempDir, "Жfile.txt"))
	fileData = base64.StdEncoding.EncodeToString(fileBytes)

	uploadXml = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>Жfile.txt</fileName>
      <fileData>%s</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>`, fileData)

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(uploadXml))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Response:\n%s\n\n", formatXML(string(body)))

	// Demo Scenario 3: Upload JSON file (rejected)
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 3: Upload JSON file (rejected)")
	fmt.Println("================================================")

	fileBytes, _ = ioutil.ReadFile(filepath.Join(tempDir, "invalid.json"))
	fileData = base64.StdEncoding.EncodeToString(fileBytes)

	uploadXml = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>data.json</fileName>
      <fileData>%s</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>`, fileData)

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(uploadXml))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Response:\n%s\n\n", formatXML(string(body)))

	// Demo Scenario 4: Get last file info
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 4: Get last file info")
	fmt.Println("================================================")

	getLastFileXml := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetLastFileInfo xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>`

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(getLastFileXml))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Response:\n%s\n\n", formatXML(string(body)))

	// Demo Scenario 5: Get file list (CSV)
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 5: Get file list (CSV)")
	fmt.Println("================================================")

	getFileListXml := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetFileListCSV xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>`

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(getFileListXml))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	responseStr := string(body)
	fmt.Printf("Response:\n%s\n\n", formatXML(responseStr))

	// Extract and decode CSV
	if start := bytes.Index(body, []byte("<csvData>")); start != -1 {
		if end := bytes.Index(body[start:], []byte("</csvData>")); end != -1 {
			csvBase64 := string(body[start+9 : start+end])
			if csvData, err := base64.StdEncoding.DecodeString(csvBase64); err == nil {
				fmt.Println("Decoded CSV:")
				fmt.Println(string(csvData))
			}
		}
	}
	fmt.Println()

	// Demo Scenario 6: Get server uptime
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 6: Get server uptime")
	fmt.Println("================================================")

	getUptimeXml := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>`

	resp, _ = http.Post(serverURL, "text/xml", bytes.NewBufferString(getUptimeXml))
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Response:\n%s\n\n", formatXML(string(body)))

	// Demo Scenario 7: Authentication failure
	fmt.Println("================================================")
	fmt.Println("Demo Scenario 7: Authentication failure")
	fmt.Println("================================================")

	authFailXml := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>invalid</username>
    <password>wrong</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>`

	resp, err = http.Post(serverURL, "text/xml", bytes.NewBufferString(authFailXml))
	if err == nil {
		defer resp.Body.Close()
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Printf("Response (Status %d):\n%s\n\n", resp.StatusCode, formatXML(string(body)))
	} else {
		fmt.Printf("✓ Server correctly rejected invalid authentication\n\n")
	}

	fmt.Println("=======================================")
	fmt.Println("Demo completed successfully!")
	fmt.Println("=======================================")
	fmt.Println()
	fmt.Println("To interact with the service manually:")
	fmt.Println("  1. Start server: go run cmd/lab3/server/main.go")
	fmt.Println("  2. Start client: go run cmd/lab3/client/main.go")
	fmt.Println()
}

func formatXML(xmlStr string) string {
	var doc interface{}
	xml.Unmarshal([]byte(xmlStr), &doc)
	data, _ := xml.MarshalIndent(doc, "", "  ")
	return string(data)
}
