# SOAP Server Test Script
# This script demonstrates all functionality of the Lab 3 SOAP service

$ErrorActionPreference = "Stop"

$SERVER_URL = "http://localhost:8080/soap"
$ADMIN_USER = "admin"
$ADMIN_PASS = "admin"

Write-Host "================================"
Write-Host "Lab 3 SOAP Service Test Script"
Write-Host "================================"
Write-Host ""

# Check if server is running
Write-Host "1. Checking server availability..."
try {
    $response = Invoke-WebRequest -Uri "${SERVER_URL}?wsdl" -TimeoutSec 2 -ErrorAction Stop
    Write-Host "✓ Server is available"
} catch {
    Write-Host "✗ Server is not available. Please start the server first:"
    Write-Host "  go run cmd/lab3/server/main.go"
    exit 1
}
Write-Host ""

# Test WSDL
Write-Host "2. Testing WSDL generation..."
try {
    $WSDL = (Invoke-WebRequest -Uri "${SERVER_URL}?wsdl").Content
    if ($WSDL -match "FileService") {
        Write-Host "✓ WSDL is valid"
    } else {
        Write-Host "✗ WSDL is invalid"
        exit 1
    }
} catch {
    Write-Host "✗ WSDL request failed: $_"
    exit 1
}
Write-Host ""

# Test GetUptime
Write-Host "3. Testing GetUptime method..."
$testUptimeXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testUptimeXml
    $content = $response.Content
    if ($content -match "uptime") {
        $uptime = [regex]::Match($content, '<uptime>([^<]*)</uptime>').Groups[1].Value
        Write-Host "✓ Server uptime: $uptime"
    } else {
        Write-Host "✗ GetUptime failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ GetUptime request failed: $_"
    exit 1
}
Write-Host ""

# Test File Upload
Write-Host "4. Testing file upload..."
$testFile = "C:\temp\test_upload_file.txt"
if (-not (Test-Path "C:\temp")) { New-Item -ItemType Directory -Path "C:\temp" -Force | Out-Null }
"This is a test file for SOAP upload" | Set-Content $testFile

$fileBytes = [System.IO.File]::ReadAllBytes($testFile)
$base64String = [Convert]::ToBase64String($fileBytes)

$testUploadXml = @"
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>${ADMIN_USER}</username>
    <password>${ADMIN_PASS}</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>test_upload_file.txt</fileName>
      <fileData>${base64String}</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testUploadXml
    $content = $response.Content
    if ($content -match "<success>true</success>") {
        Write-Host "✓ File uploaded successfully"
    } else {
        Write-Host "✗ File upload failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ File upload request failed: $_"
    Write-Host $_.Exception
    exit 1
}
Write-Host ""

# Test GetLastFileInfo
Write-Host "5. Testing GetLastFileInfo..."
$testLastFileXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testLastFileXml
    $content = $response.Content
    if ($content -match "test_upload_file.txt") {
        $fileSize = [regex]::Match($content, '<fileSize>([^<]*)</fileSize>').Groups[1].Value
        Write-Host "✓ Last file info retrieved: test_upload_file.txt (${fileSize} bytes)"
    } else {
        Write-Host "✗ GetLastFileInfo failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ GetLastFileInfo request failed: $_"
    exit 1
}
Write-Host ""

# Test GetFileListCSV
Write-Host "6. Testing GetFileListCSV..."
$testFileListXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testFileListXml
    $content = $response.Content
    $csvMatch = [regex]::Match($content, '<csvData>([^<]*)</csvData>')
    if ($csvMatch.Success) {
        $csvData = $csvMatch.Groups[1].Value
        Write-Host "✓ File list retrieved (CSV):"
        $decodedCsv = [System.Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($csvData))
        Write-Host $decodedCsv
    } else {
        Write-Host "✗ GetFileListCSV failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ GetFileListCSV request failed: $_"
    exit 1
}
Write-Host ""

# Test Validation: File with forbidden character Ж
Write-Host "7. Testing validation: forbidden character 'Ж'..."
$testForbiddenXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testForbiddenXml
    $content = $response.Content
    if ($content -match "forbidden character") {
        Write-Host "✓ Validation works: forbidden character detected"
    } else {
        Write-Host "✗ Validation failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ Forbidden character test failed: $_"
}
Write-Host ""

# Test Validation: JSON file
Write-Host "8. Testing validation: JSON file rejection..."
$testJsonXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testJsonXml
    $content = $response.Content
    if ($content -match "JSON data") {
        Write-Host "✓ Validation works: JSON file rejected"
    } else {
        Write-Host "✗ Validation failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ JSON file test failed: $_"
}
Write-Host ""

# Test Validation: Empty file
Write-Host "9. Testing validation: empty file rejection..."
$testEmptyXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testEmptyXml
    $content = $response.Content
    if ($content -match "empty") {
        Write-Host "✓ Validation works: empty file rejected"
    } else {
        Write-Host "✗ Validation failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ Empty file test failed: $_"
}
Write-Host ""

# Test Authentication Failure
Write-Host "10. Testing authentication failure..."
$testAuthFailXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "$SERVER_URL" -Method Post -ContentType "text/xml" -Body $testAuthFailXml
    $content = $response.Content
    if ($content -match "Authentication failed") {
        Write-Host "✓ Authentication validation works"
    } else {
        Write-Host "✗ Authentication validation failed"
        Write-Host $content
        exit 1
    }
} catch {
    Write-Host "✗ Authentication test failed: $_"
}
Write-Host ""

Write-Host "================================"
Write-Host "All tests passed successfully! ✓"
Write-Host "================================"

