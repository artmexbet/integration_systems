#!/bin/bash

# Lab 3 Demo Script - Complete workflow demonstration
# This script demonstrates the full client-server interaction

set -e

echo "======================================="
echo "Lab 3: SOAP Service Complete Demo"
echo "======================================="
echo ""

# Check if server is running
if ! curl -s --max-time 2 http://localhost:8080/soap?wsdl > /dev/null 2>&1; then
    echo "Starting SOAP server..."
    cd /home/runner/work/integration_systems/integration_systems
    go run cmd/lab3/server/main.go &
    SERVER_PID=$!
    echo "Server PID: $SERVER_PID"
    sleep 3
    echo ""
fi

echo "Server is ready!"
echo ""
echo "WSDL is available at: http://localhost:8080/soap?wsdl"
echo ""

# Display available credentials
echo "Available credentials:"
echo "  - user1:pass1"
echo "  - user2:pass2"
echo "  - admin:admin"
echo ""

# Create test files
echo "Creating test files..."
mkdir -p /tmp/lab3_demo
echo "This is a valid text file for upload" > /tmp/lab3_demo/valid.txt
echo '{"name":"test","value":123}' > /tmp/lab3_demo/invalid.json
echo "" > /tmp/lab3_demo/empty.txt
echo "Test file with русский текст" > /tmp/lab3_demo/Жfile.txt
echo ""

echo "Test files created in /tmp/lab3_demo/:"
ls -lh /tmp/lab3_demo/
echo ""

echo "================================================"
echo "Demo Scenario 1: Successful file upload"
echo "================================================"
FILE_DATA=$(base64 /tmp/lab3_demo/valid.txt)
cat > /tmp/lab3_demo/upload_valid.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>valid.txt</fileName>
      <fileData>${FILE_DATA}</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/upload_valid.xml | xmllint --format -
echo ""

echo "================================================"
echo "Demo Scenario 2: Upload with forbidden character"
echo "================================================"
FILE_DATA=$(base64 /tmp/lab3_demo/Жfile.txt)
cat > /tmp/lab3_demo/upload_forbidden.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>Жfile.txt</fileName>
      <fileData>${FILE_DATA}</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/upload_forbidden.xml | xmllint --format -
echo ""

echo "================================================"
echo "Demo Scenario 3: Upload JSON file (rejected)"
echo "================================================"
FILE_DATA=$(base64 /tmp/lab3_demo/invalid.json)
cat > /tmp/lab3_demo/upload_json.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>data.json</fileName>
      <fileData>${FILE_DATA}</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
EOF
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/upload_json.xml | xmllint --format -
echo ""

echo "================================================"
echo "Demo Scenario 4: Get last file info"
echo "================================================"
cat > /tmp/lab3_demo/get_lastfile.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetLastFileInfo xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/get_lastfile.xml | xmllint --format -
echo ""

echo "================================================"
echo "Demo Scenario 5: Get file list (CSV)"
echo "================================================"
cat > /tmp/lab3_demo/get_filelist.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetFileListCSV xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
RESPONSE=$(curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/get_filelist.xml)
echo "$RESPONSE" | xmllint --format -
echo ""
echo "Decoded CSV:"
echo "$RESPONSE" | grep -o '<csvData>[^<]*</csvData>' | sed 's/<[^>]*>//g' | base64 -d
echo ""

echo "================================================"
echo "Demo Scenario 6: Get server uptime"
echo "================================================"
cat > /tmp/lab3_demo/get_uptime.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <GetUptime xmlns="http://tempuri.org/"/>
  </soap:Body>
</soap:Envelope>
EOF
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/get_uptime.xml | xmllint --format -
echo ""

echo "================================================"
echo "Demo Scenario 7: Authentication failure"
echo "================================================"
cat > /tmp/lab3_demo/auth_fail.xml << EOF
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
curl -s -X POST http://localhost:8080/soap -H "Content-Type: text/xml" -d @/tmp/lab3_demo/auth_fail.xml | xmllint --format -
echo ""

echo "======================================="
echo "Demo completed successfully!"
echo "======================================="
echo ""
echo "To interact with the service manually:"
echo "  1. Start server: go run cmd/lab3/server/main.go"
echo "  2. Start client: go run cmd/lab3/client/main.go"
echo ""
echo "Or run automated tests:"
echo "  bash cmd/lab3/test_server.sh"
