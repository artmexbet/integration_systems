package main

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SOAP Envelope structures
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Username string   `xml:"username,omitempty"`
	Password string   `xml:"password,omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Content interface{}
}

type SOAPFault struct {
	XMLName     xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
}

// Service request/response structures
type UploadFileRequest struct {
	XMLName      xml.Name `xml:"http://tempuri.org/ UploadFile"`
	FileName     string   `xml:"fileName"`
	FileData     string   `xml:"fileData"` // Base64 encoded
	CallbackURL  string   `xml:"callbackURL"`
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
	XMLName  xml.Name `xml:"http://tempuri.org/ GetLastFileInfoResponse"`
	FileName string   `xml:"fileName"`
	FileSize int64    `xml:"fileSize"`
	UploadTime string `xml:"uploadTime"`
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
	XMLName xml.Name `xml:"http://tempuri.org/ UploadNotification"`
	Success bool     `xml:"success"`
	Message string   `xml:"message"`
	FileName string  `xml:"fileName"`
}

// FileInfo represents stored file metadata
type FileInfo struct {
	FileName   string
	FileSize   int64
	UploadTime time.Time
	Username   string
}

// FileStorage manages uploaded files
type FileStorage struct {
	mu           sync.RWMutex
	files        map[string]FileInfo // username -> last file
	allFiles     []FileInfo
	storageDir   string
	maxStorage   int64
	currentUsage int64
}

var (
	storage    *FileStorage
	startTime  time.Time
	validUsers = map[string]string{
		"user1": "pass1",
		"user2": "pass2",
		"admin": "admin",
	}
)

func init() {
	startTime = time.Now()
	storageDir := "./uploaded_files"
	os.MkdirAll(storageDir, 0755)
	
	storage = &FileStorage{
		files:      make(map[string]FileInfo),
		allFiles:   make([]FileInfo, 0),
		storageDir: storageDir,
		maxStorage: 100 * 1024 * 1024, // 100 MB max storage
	}
}

func main() {
	http.HandleFunc("/soap", handleSOAP)
	
	log.Println("SOAP Server starting on :8080")
	log.Println("WSDL available at: http://localhost:8080/soap?wsdl")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleWSDL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	w.Write([]byte(getWSDL()))
}

func handleSOAP(w http.ResponseWriter, r *http.Request) {
	// Check for WSDL request first
	if _, exists := r.URL.Query()["wsdl"]; exists {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(getWSDL()))
		return
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var envelope SOAPEnvelope
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		sendSOAPFault(w, "Client", "Failed to read request body")
		return
	}
	
	// Parse envelope
	if err := xml.Unmarshal(bodyData, &envelope); err != nil {
		sendSOAPFault(w, "Client", "Invalid SOAP request: "+err.Error())
		return
	}

	// Authenticate
	username := ""
	if envelope.Header != nil {
		if !authenticate(envelope.Header.Username, envelope.Header.Password) {
			sendSOAPFault(w, "Server", "Authentication failed")
			return
		}
		username = envelope.Header.Username
	} else {
		sendSOAPFault(w, "Server", "Authentication required")
		return
	}

	// Process request based on the full body content
	var response interface{}
	
	bodyStr := string(bodyData)
	if strings.Contains(bodyStr, "UploadFile") && !strings.Contains(bodyStr, "UploadFileResponse") {
		var envWithReq struct {
			XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
			Body    struct {
				UploadFile UploadFileRequest `xml:"UploadFile"`
			} `xml:"Body"`
		}
		xml.Unmarshal(bodyData, &envWithReq)
		response = handleUploadFile(envWithReq.Body.UploadFile, username)
	} else if strings.Contains(bodyStr, "GetLastFileInfo") {
		response = handleGetLastFileInfo(username)
	} else if strings.Contains(bodyStr, "GetFileListCSV") {
		response = handleGetFileListCSV()
	} else if strings.Contains(bodyStr, "GetUptime") {
		response = handleGetUptime()
	} else {
		sendSOAPFault(w, "Client", "Unknown operation")
		return
	}

	sendSOAPResponse(w, response)
}

func authenticate(username, password string) bool {
	expectedPass, exists := validUsers[username]
	return exists && expectedPass == password
}

func handleUploadFile(req UploadFileRequest, username string) UploadFileResponse {
	// Decode file data
	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		return UploadFileResponse{Success: false, Message: "Invalid file encoding"}
	}

	// Validate file
	if err := validateFile(req.FileName, fileData); err != nil {
		resp := UploadFileResponse{Success: false, Message: err.Error()}
		// Send async notification
		go sendNotification(req.CallbackURL, req.FileName, false, err.Error())
		return resp
	}

	// Save file
	filePath := filepath.Join(storage.storageDir, fmt.Sprintf("%s_%s", username, req.FileName))
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		resp := UploadFileResponse{Success: false, Message: "Failed to save file"}
		go sendNotification(req.CallbackURL, req.FileName, false, "Failed to save file")
		return resp
	}

	// Update storage
	fileInfo := FileInfo{
		FileName:   req.FileName,
		FileSize:   int64(len(fileData)),
		UploadTime: time.Now(),
		Username:   username,
	}
	
	storage.mu.Lock()
	storage.files[username] = fileInfo
	storage.allFiles = append(storage.allFiles, fileInfo)
	storage.currentUsage += fileInfo.FileSize
	storage.mu.Unlock()

	// Send async notification
	go sendNotification(req.CallbackURL, req.FileName, true, "File uploaded successfully")

	return UploadFileResponse{Success: true, Message: "File uploaded successfully"}
}

func validateFile(fileName string, fileData []byte) error {
	fileSize := int64(len(fileData))
	
	// Check if file is empty
	if fileSize == 0 {
		return fmt.Errorf("file is empty")
	}
	
	// Check size limit (3 MB)
	if fileSize > 3*1024*1024 {
		return fmt.Errorf("file size exceeds 3 MB limit")
	}
	
	// Check if filename contains 'Ж' or 'ж'
	if strings.ContainsAny(fileName, "Жж") {
		return fmt.Errorf("filename contains forbidden character 'Ж'")
	}
	
	// Check if file contains only valid JSON
	var js interface{}
	if err := json.Unmarshal(fileData, &js); err == nil {
		// File is valid JSON, reject it
		trimmed := strings.TrimSpace(string(fileData))
		if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
			(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
			return fmt.Errorf("file contains only JSON data")
		}
	}
	
	// Check storage space
	storage.mu.RLock()
	available := storage.maxStorage - storage.currentUsage
	storage.mu.RUnlock()
	
	if fileSize > available {
		return fmt.Errorf("insufficient storage space")
	}
	
	return nil
}

func sendNotification(callbackURL, fileName string, success bool, message string) {
	if callbackURL == "" {
		return
	}

	notification := UploadNotification{
		Success:  success,
		Message:  message,
		FileName: fileName,
	}

	envelope := SOAPEnvelope{
		Body: SOAPBody{Content: notification},
	}

	data, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", callbackURL, strings.NewReader(xml.Header+string(data)))
	if err != nil {
		log.Printf("Failed to create notification request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Notification sent to %s: %v", callbackURL, success)
}

func handleGetLastFileInfo(username string) GetLastFileInfoResponse {
	storage.mu.RLock()
	fileInfo, exists := storage.files[username]
	storage.mu.RUnlock()

	if !exists {
		return GetLastFileInfoResponse{
			FileName: "",
			FileSize: 0,
			UploadTime: "",
		}
	}

	return GetLastFileInfoResponse{
		FileName:   fileInfo.FileName,
		FileSize:   fileInfo.FileSize,
		UploadTime: fileInfo.UploadTime.Format(time.RFC3339),
	}
}

func handleGetFileListCSV() GetFileListCSVResponse {
	storage.mu.RLock()
	files := make([]FileInfo, len(storage.allFiles))
	copy(files, storage.allFiles)
	storage.mu.RUnlock()

	// Create CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)
	
	// Write header
	writer.Write([]string{"Username", "FileName", "FileSize", "UploadTime"})
	
	// Write data
	for _, f := range files {
		writer.Write([]string{
			f.Username,
			f.FileName,
			fmt.Sprintf("%d", f.FileSize),
			f.UploadTime.Format(time.RFC3339),
		})
	}
	writer.Flush()

	// Encode as base64
	csvData := base64.StdEncoding.EncodeToString([]byte(buf.String()))

	return GetFileListCSVResponse{
		CSVData: csvData,
	}
}

func handleGetUptime() GetUptimeResponse {
	uptime := time.Since(startTime)
	return GetUptimeResponse{
		Uptime: uptime.String(),
	}
}

func sendSOAPResponse(w http.ResponseWriter, response interface{}) {
	envelope := SOAPEnvelope{
		Body: SOAPBody{Content: response},
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(xml.Header))
	if err := xml.NewEncoder(w).Encode(envelope); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func sendSOAPFault(w http.ResponseWriter, code, message string) {
	fault := SOAPFault{
		FaultCode:   code,
		FaultString: message,
	}

	envelope := SOAPEnvelope{
		Body: SOAPBody{Content: fault},
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(envelope)
}

func getWSDL() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<definitions name="FileService"
   targetNamespace="http://tempuri.org/"
   xmlns="http://schemas.xmlsoap.org/wsdl/"
   xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
   xmlns:tns="http://tempuri.org/"
   xmlns:xsd="http://www.w3.org/2001/XMLSchema">

   <types>
      <xsd:schema targetNamespace="http://tempuri.org/">
         <xsd:element name="UploadFile">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="fileName" type="xsd:string"/>
                  <xsd:element name="fileData" type="xsd:string"/>
                  <xsd:element name="callbackURL" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
         <xsd:element name="UploadFileResponse">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="success" type="xsd:boolean"/>
                  <xsd:element name="message" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
         
         <xsd:element name="GetLastFileInfo">
            <xsd:complexType>
               <xsd:sequence/>
            </xsd:complexType>
         </xsd:element>
         <xsd:element name="GetLastFileInfoResponse">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="fileName" type="xsd:string"/>
                  <xsd:element name="fileSize" type="xsd:long"/>
                  <xsd:element name="uploadTime" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
         
         <xsd:element name="GetFileListCSV">
            <xsd:complexType>
               <xsd:sequence/>
            </xsd:complexType>
         </xsd:element>
         <xsd:element name="GetFileListCSVResponse">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="csvData" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
         
         <xsd:element name="GetUptime">
            <xsd:complexType>
               <xsd:sequence/>
            </xsd:complexType>
         </xsd:element>
         <xsd:element name="GetUptimeResponse">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="uptime" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
         
         <xsd:element name="UploadNotification">
            <xsd:complexType>
               <xsd:sequence>
                  <xsd:element name="success" type="xsd:boolean"/>
                  <xsd:element name="message" type="xsd:string"/>
                  <xsd:element name="fileName" type="xsd:string"/>
               </xsd:sequence>
            </xsd:complexType>
         </xsd:element>
      </xsd:schema>
   </types>

   <message name="UploadFileRequest">
      <part name="parameters" element="tns:UploadFile"/>
   </message>
   <message name="UploadFileResponse">
      <part name="parameters" element="tns:UploadFileResponse"/>
   </message>
   
   <message name="GetLastFileInfoRequest">
      <part name="parameters" element="tns:GetLastFileInfo"/>
   </message>
   <message name="GetLastFileInfoResponse">
      <part name="parameters" element="tns:GetLastFileInfoResponse"/>
   </message>
   
   <message name="GetFileListCSVRequest">
      <part name="parameters" element="tns:GetFileListCSV"/>
   </message>
   <message name="GetFileListCSVResponse">
      <part name="parameters" element="tns:GetFileListCSVResponse"/>
   </message>
   
   <message name="GetUptimeRequest">
      <part name="parameters" element="tns:GetUptime"/>
   </message>
   <message name="GetUptimeResponse">
      <part name="parameters" element="tns:GetUptimeResponse"/>
   </message>

   <portType name="FileServicePortType">
      <operation name="UploadFile">
         <input message="tns:UploadFileRequest"/>
         <output message="tns:UploadFileResponse"/>
      </operation>
      <operation name="GetLastFileInfo">
         <input message="tns:GetLastFileInfoRequest"/>
         <output message="tns:GetLastFileInfoResponse"/>
      </operation>
      <operation name="GetFileListCSV">
         <input message="tns:GetFileListCSVRequest"/>
         <output message="tns:GetFileListCSVResponse"/>
      </operation>
      <operation name="GetUptime">
         <input message="tns:GetUptimeRequest"/>
         <output message="tns:GetUptimeResponse"/>
      </operation>
   </portType>

   <binding name="FileServiceBinding" type="tns:FileServicePortType">
      <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http"/>
      <operation name="UploadFile">
         <soap:operation soapAction="http://tempuri.org/UploadFile"/>
         <input>
            <soap:body use="literal"/>
         </input>
         <output>
            <soap:body use="literal"/>
         </output>
      </operation>
      <operation name="GetLastFileInfo">
         <soap:operation soapAction="http://tempuri.org/GetLastFileInfo"/>
         <input>
            <soap:body use="literal"/>
         </input>
         <output>
            <soap:body use="literal"/>
         </output>
      </operation>
      <operation name="GetFileListCSV">
         <soap:operation soapAction="http://tempuri.org/GetFileListCSV"/>
         <input>
            <soap:body use="literal"/>
         </input>
         <output>
            <soap:body use="literal"/>
         </output>
      </operation>
      <operation name="GetUptime">
         <soap:operation soapAction="http://tempuri.org/GetUptime"/>
         <input>
            <soap:body use="literal"/>
         </input>
         <output>
            <soap:body use="literal"/>
         </output>
      </operation>
   </binding>

   <service name="FileService">
      <documentation>File Upload and Management Service</documentation>
      <port name="FileServicePort" binding="tns:FileServiceBinding">
         <soap:address location="http://localhost:8080/soap"/>
      </port>
   </service>
</definitions>`
}
