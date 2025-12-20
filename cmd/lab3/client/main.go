package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SOAPEnvelope structures
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	XMLNS   string   `xml:"xmlns,attr"`
	Header  *SOAPHeader
	Body    *SOAPBody
}

type SOAPHeader struct {
	XMLName  xml.Name `xml:"Header"`
	Username string   `xml:"username,omitempty"`
	Password string   `xml:"password,omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"Body"`
	Content string   `xml:",innerxml"`
}

type SOAPFault struct {
	XMLName     xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
}

// Service request/response structures
type UploadFileRequest struct {
	XMLName     xml.Name `xml:"http://tempuri.org/ UploadFile"`
	FileName    string   `xml:"fileName"`
	FileData    string   `xml:"fileData"` // Base64 encoded
	CallbackURL string   `xml:"callbackURL"`
}

type UploadFileResponse struct {
	XMLName xml.Name `xml:"http://tempuri.org/ UploadFileResponse"`
	Success bool     `xml:"success"`
	Message string   `xml:"message"`
}

type GetLastFileInfoRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ GetLastFileInfo"`
}

type GetLastFileInfoResponse struct {
	XMLName    xml.Name `xml:"http://tempuri.org/ GetLastFileInfoResponse"`
	FileName   string   `xml:"fileName"`
	FileSize   int64    `xml:"fileSize"`
	UploadTime string   `xml:"uploadTime"`
}

type GetFileListCSVRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ GetFileListCSV"`
}

type GetFileListCSVResponse struct {
	XMLName xml.Name `xml:"http://tempuri.org/ GetFileListCSVResponse"`
	CSVData string   `xml:"csvData"` // Base64 encoded CSV
}

type GetUptimeRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ GetUptime"`
}

type GetUptimeResponse struct {
	XMLName xml.Name `xml:"http://tempuri.org/ GetUptimeResponse"`
	Uptime  string   `xml:"uptime"`
}

type UploadNotification struct {
	XMLName  xml.Name `xml:"http://tempuri.org/ UploadNotification"`
	Success  bool     `xml:"success"`
	Message  string   `xml:"message"`
	FileName string   `xml:"fileName"`
}

// Client represents SOAP client
type Client struct {
	serverURL      string
	username       string
	password       string
	webhookPort    int
	notificationCh chan UploadNotification
	httpClient     *http.Client
}

func main() {
	client := &Client{
		serverURL:      "http://localhost:8080/soap",
		webhookPort:    9090,
		notificationCh: make(chan UploadNotification, 10),
		httpClient:     &http.Client{Timeout: 30 * time.Second},
	}

	// Start webhook server for notifications
	go client.startWebhookServer()

	// Wait a bit for server to start
	time.Sleep(time.Second)

	fmt.Println("=== SOAP File Upload Client ===")
	fmt.Println()

	// Check server availability
	if !client.checkServerAvailability() {
		log.Fatal("Server is not available")
	}
	fmt.Println("✓ Server is available")
	fmt.Println()

	// Authenticate
	if !client.authenticate() {
		log.Fatal("Failed to authenticate")
	}

	// Main menu loop
	for {
		fmt.Println("\n--- Main Menu ---")
		fmt.Println("1. Upload file")
		fmt.Println("2. Check last uploaded file info")
		fmt.Println("3. Get file list (CSV)")
		fmt.Println("4. Get server uptime")
		fmt.Println("5. Exit")
		fmt.Print("Choose option: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			client.uploadFileFlow()
		case "2":
			client.getLastFileInfo()
		case "3":
			client.getFileList()
		case "4":
			client.getUptime()
		case "5":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func (c *Client) checkServerAvailability() bool {
	// Try to get WSDL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.serverURL+"?wsdl", nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (c *Client) authenticate() bool {
	fmt.Print("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	c.username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	c.password = strings.TrimSpace(password)

	// Test authentication with GetUptime request
	request := GetUptimeRequest{}
	var response GetUptimeResponse

	err := c.sendSOAPRequest(request, &response)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		return false
	}

	fmt.Println("✓ Authentication successful")
	return true
}

func (c *Client) uploadFileFlow() {
	fmt.Print("Enter file path to upload: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	fileName := filepath.Base(filePath)
	fileDataEncoded := base64.StdEncoding.EncodeToString(fileData)

	// Prepare callback URL
	callbackURL := fmt.Sprintf("http://localhost:%d/webhook", c.webhookPort)

	request := UploadFileRequest{
		FileName:    fileName,
		FileData:    fileDataEncoded,
		CallbackURL: callbackURL,
	}

	var response UploadFileResponse
	err = c.sendSOAPRequest(request, &response)
	if err != nil {
		fmt.Printf("Upload failed: %v\n", err)
		return
	}

	fmt.Printf("Upload initiated: %s - %s\n", fileName, response.Message)

	if response.Success {
		// Wait for notification with timeout
		fmt.Println("Waiting for upload confirmation...")
		select {
		case notification := <-c.notificationCh:
			fmt.Printf("\n✓ Notification received: %s - %s\n", notification.FileName, notification.Message)
		case <-time.After(15 * time.Second):
			fmt.Println("\n⚠ Notification timeout, checking upload status proactively...")
			c.getLastFileInfo()
		}
	}
}

func (c *Client) getLastFileInfo() {
	request := GetLastFileInfoRequest{}
	var response GetLastFileInfoResponse

	err := c.sendSOAPRequest(request, &response)
	if err != nil {
		fmt.Printf("Failed to get file info: %v\n", err)
		return
	}

	if response.FileName == "" {
		fmt.Println("No files uploaded yet by this user")
	} else {
		fmt.Printf("\nLast uploaded file:\n")
		fmt.Printf("  Name: %s\n", response.FileName)
		fmt.Printf("  Size: %d bytes\n", response.FileSize)
		fmt.Printf("  Upload Time: %s\n", response.UploadTime)
	}
}

func (c *Client) getFileList() {
	request := GetFileListCSVRequest{}
	var response GetFileListCSVResponse

	err := c.sendSOAPRequest(request, &response)
	if err != nil {
		fmt.Printf("Failed to get file list: %v\n", err)
		return
	}

	// Decode CSV
	csvData, err := base64.StdEncoding.DecodeString(response.CSVData)
	if err != nil {
		fmt.Printf("Failed to decode CSV: %v\n", err)
		return
	}

	fmt.Println("\nFile List (CSV):")
	fmt.Println("---")
	fmt.Println(string(csvData))
	fmt.Println("---")
}

func (c *Client) getUptime() {
	request := GetUptimeRequest{}
	var response GetUptimeResponse

	err := c.sendSOAPRequest(request, &response)
	if err != nil {
		fmt.Printf("Failed to get uptime: %v\n", err)
		return
	}

	fmt.Printf("\nServer uptime: %s\n", response.Uptime)
}

func (c *Client) sendSOAPRequest(request interface{}, response interface{}) error {
	// Create SOAP envelope
	envelope := SOAPEnvelope{
		XMLNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Header: &SOAPHeader{
			Username: c.username,
			Password: c.password,
		},
		Body: &SOAPBody{},
	}

	// Marshal request content
	requestBytes, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	envelope.Body.Content = string(requestBytes)

	// Marshal request
	data, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	// Send request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", c.serverURL, strings.NewReader(xml.Header+string(data)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Debug: print raw response
	fmt.Printf("DEBUG: Raw response:\n%s\n", string(respData))

	// Parse response envelope
	var respEnvelope SOAPEnvelope
	if err := xml.Unmarshal(respData, &respEnvelope); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Debug: check Body
	fmt.Printf("DEBUG: Body is nil? %v\n", respEnvelope.Body == nil)
	if respEnvelope.Body != nil {
		fmt.Printf("DEBUG: Body.Content type: %T\n", respEnvelope.Body.Content)
		fmt.Printf("DEBUG: Body.Content value: %v\n", respEnvelope.Body.Content)
	}

	if respEnvelope.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	// Get body content as bytes
	bodyBytes := []byte(respEnvelope.Body.Content)

	// Check for fault
	if strings.Contains(string(bodyBytes), "Fault") {
		var fault SOAPFault
		if err := xml.Unmarshal(bodyBytes, &fault); err == nil {
			return fmt.Errorf("SOAP fault: %s - %s", fault.FaultCode, fault.FaultString)
		}
	}

	// Unmarshal response
	if err := xml.Unmarshal(bodyBytes, response); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

func (c *Client) startWebhookServer() {
	http.HandleFunc("/webhook", c.handleWebhook)

	addr := fmt.Sprintf(":%d", c.webhookPort)
	log.Printf("Webhook server listening on %s\n", addr)

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Webhook server error: %v", err)
	}
}

func (c *Client) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		log.Printf("Failed to decode webhook: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Parse notification
	bodyBytes := []byte(envelope.Body.Content)
	log.Printf("Webhook received body: %s", string(bodyBytes))
	var notification UploadNotification
	if err := xml.Unmarshal(bodyBytes, &notification); err != nil {
		log.Printf("Failed to unmarshal notification: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Send notification to channel
	select {
	case c.notificationCh <- notification:
		log.Printf("Received notification: %s - %v", notification.FileName, notification.Success)
	default:
		log.Println("Notification channel full")
	}

	w.WriteHeader(http.StatusOK)
}
