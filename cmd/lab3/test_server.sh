#!/bin/bash

# SOAP Server Test Script
# This script demonstrates all functionality of the Lab 3 SOAP service

set -e

SERVER_URL="http://localhost:8080/soap"
ADMIN_USER="admin"
ADMIN_PASS="admin"

echo "================================"
echo "Lab 3 SOAP Service Test Script"
echo "================================"
echo ""

# Check if server is running
echo "1. Checking server availability..."
if curl -s --max-time 2 "${SERVER_URL}?wsdl" > /dev/null 2>&1; then
    echo "✓ Server is available"
else
    echo "✗ Server is not available. Please start the server first:"
    echo "  go run cmd/lab3/server/main.go"
    exit 1
fi
echo ""

# Test WSDL
echo "2. Testing WSDL generation..."
WSDL=$(curl -s "${SERVER_URL}?wsdl")
if echo "$WSDL" | grep -q "FileService"; then
    echo "✓ WSDL is valid"
else
    echo "✗ WSDL is invalid"
    exit 1
fi
echo ""

# Test GetUptime
echo "3. Testing GetUptime method..."
cat > /tmp/test_uptime.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_uptime.xml)
if echo "$RESPONSE" | grep -q "uptime"; then
    UPTIME=$(echo "$RESPONSE" | grep -o '<uptime>[^<]*</uptime>' | sed 's/<[^>]*>//g')
    echo "✓ Server uptime: $UPTIME"
else
    echo "✗ GetUptime failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test File Upload
echo "4. Testing file upload..."
echo "This is a test file for SOAP upload" > /tmp/test_upload_file.txt
FILE_DATA=$(base64 /tmp/test_upload_file.txt)
cat > /tmp/test_upload.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>test_upload_file.txt</fileName>
      <fileData>${FILE_DATA}</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_upload.xml)
if echo "$RESPONSE" | grep -q "<success>true</success>"; then
    echo "✓ File uploaded successfully"
else
    echo "✗ File upload failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test GetLastFileInfo
echo "5. Testing GetLastFileInfo..."
cat > /tmp/test_lastfile.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <GetLastFileInfo xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_lastfile.xml)
if echo "$RESPONSE" | grep -q "test_upload_file.txt"; then
    FILESIZE=$(echo "$RESPONSE" | grep -o '<fileSize>[^<]*</fileSize>' | sed 's/<[^>]*>//g')
    echo "✓ Last file info retrieved: test_upload_file.txt (${FILESIZE} bytes)"
else
    echo "✗ GetLastFileInfo failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test GetFileListCSV
echo "6. Testing GetFileListCSV..."
cat > /tmp/test_filelist.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <GetFileListCSV xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_filelist.xml)
CSV_DATA=$(echo "$RESPONSE" | grep -o '<csvData>[^<]*</csvData>' | sed 's/<[^>]*>//g')
if [ -n "$CSV_DATA" ]; then
    echo "✓ File list retrieved (CSV):"
    echo "$CSV_DATA" | base64 -d
else
    echo "✗ GetFileListCSV failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test Validation: File with forbidden character Ж
echo "7. Testing validation: forbidden character 'Ж'..."
cat > /tmp/test_forbidden.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>testЖfile.txt</fileName>
      <fileData>VGVzdCBkYXRh</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_forbidden.xml)
if echo "$RESPONSE" | grep -q "forbidden character"; then
    echo "✓ Validation works: forbidden character detected"
else
    echo "✗ Validation failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test Validation: JSON file
echo "8. Testing validation: JSON file rejection..."
cat > /tmp/test_json.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>data.json</fileName>
      <fileData>eyJuYW1lIjoidGVzdCIsInZhbHVlIjoxMjN9</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_json.xml)
if echo "$RESPONSE" | grep -q "JSON data"; then
    echo "✓ Validation works: JSON file rejected"
else
    echo "✗ Validation failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test Validation: Empty file
echo "9. Testing validation: empty file rejection..."
cat > /tmp/test_empty.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>empty.txt</fileName>
      <fileData></fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_empty.xml)
if echo "$RESPONSE" | grep -q "empty"; then
    echo "✓ Validation works: empty file rejected"
else
    echo "✗ Validation failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

# Test Authentication Failure
echo "10. Testing authentication failure..."
cat > /tmp/test_auth_fail.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>invalid</username>
    <password>wrong</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST "$SERVER_URL" -H "Content-Type: text/xml" -d @/tmp/test_auth_fail.xml)
if echo "$RESPONSE" | grep -q "Authentication failed"; then
    echo "✓ Authentication validation works"
else
    echo "✗ Authentication validation failed"
    echo "$RESPONSE"
    exit 1
fi
echo ""

echo "================================"
echo "All tests passed successfully! ✓"
echo "================================"
