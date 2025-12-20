# Lab 3 Demo Script - Complete workflow demonstration
# This script demonstrates the full client-server interaction

$ErrorActionPreference = "Stop"

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "Lab 3: SOAP Service Complete Demo" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host ""

# Check if server is running
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap?wsdl" -TimeoutSec 2 -ErrorAction SilentlyContinue -UseBasicParsing
    Write-Host "Server is ready!" -ForegroundColor Green
} catch {
    Write-Host "Starting SOAP server..." -ForegroundColor Yellow
    # Start server in background - adjust path if needed
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd 'C:\Users\artem\GolandProjects\RIS'; go run cmd/lab3/server/main.go"
    Write-Host "Waiting for server to start..." -ForegroundColor Yellow
    Start-Sleep -Seconds 3
    Write-Host "Server is ready!" -ForegroundColor Green
}

Write-Host ""
Write-Host "WSDL is available at: http://localhost:8080/soap?wsdl"
Write-Host ""

# Display available credentials
Write-Host "Available credentials:"
Write-Host "  - user1:pass1"
Write-Host "  - user2:pass2"
Write-Host "  - admin:admin"
Write-Host ""

# Create test files
Write-Host "Creating test files..." -ForegroundColor Yellow
$tempDir = "C:\temp\lab3_demo"
if (-not (Test-Path $tempDir)) { New-Item -ItemType Directory -Path $tempDir -Force | Out-Null }

"This is a valid text file for upload" | Set-Content "$tempDir\valid.txt"
'{"name":"test","value":123}' | Set-Content "$tempDir\invalid.json"
"" | Set-Content "$tempDir\empty.txt"
"Test file with русский текст" | Set-Content "$tempDir\Жfile.txt"

Write-Host ""
Write-Host "Test files created in $tempDir`:"
Get-ChildItem -Path $tempDir -Force | Format-Table -AutoSize
Write-Host ""

# Helper function to format XML output
function Format-XML {
    param([string]$xml)
    try {
        $xmlDoc = New-Object System.Xml.XmlDocument
        $xmlDoc.LoadXml($xml)
        $stringWriter = New-Object System.IO.StringWriter
        $xmlWriter = New-Object System.Xml.XmlTextWriter($stringWriter)
        $xmlWriter.Formatting = [System.Xml.Formatting]::Indented
        $xmlDoc.WriteTo($xmlWriter)
        return $stringWriter.ToString()
    } catch {
        return $xml
    }
}

# Demo Scenario 1: Successful file upload
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 1: Successful file upload" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$fileBytes = [System.IO.File]::ReadAllBytes("$tempDir\valid.txt")
$fileData = [Convert]::ToBase64String($fileBytes)

$uploadXml = @"
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>valid.txt</fileName>
      <fileData>$fileData</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $uploadXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Green
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 2: Upload with forbidden character
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 2: Upload with forbidden character" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$fileBytes = [System.IO.File]::ReadAllBytes("$tempDir\Жfile.txt")
$fileData = [Convert]::ToBase64String($fileBytes)

$uploadXml = @"
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>Жfile.txt</fileName>
      <fileData>$fileData</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $uploadXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Yellow
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 3: Upload JSON file (rejected)
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 3: Upload JSON file (rejected)" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$fileBytes = [System.IO.File]::ReadAllBytes("$tempDir\invalid.json")
$fileData = [Convert]::ToBase64String($fileBytes)

$uploadXml = @"
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <username>admin</username>
    <password>admin</password>
  </soap:Header>
  <soap:Body>
    <UploadFile xmlns="http://tempuri.org/">
      <fileName>data.json</fileName>
      <fileData>$fileData</fileData>
      <callbackURL></callbackURL>
    </UploadFile>
  </soap:Body>
</soap:Envelope>
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $uploadXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Yellow
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 4: Get last file info
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 4: Get last file info" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$getLastFileXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $getLastFileXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Green
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 5: Get file list (CSV)
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 5: Get file list (CSV)" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$getFileListXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $getFileListXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Green

    # Extract and decode CSV data
    $csvMatch = [regex]::Match($response.Content, '<csvData>([^<]*)</csvData>')
    if ($csvMatch.Success) {
        Write-Host ""
        Write-Host "Decoded CSV:" -ForegroundColor Cyan
        $csvData = $csvMatch.Groups[1].Value
        $decodedCsv = [System.Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($csvData))
        Write-Host $decodedCsv -ForegroundColor White
    }
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 6: Get server uptime
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 6: Get server uptime" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$getUptimeXml = @"
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
"@

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $getUptimeXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Green
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Demo Scenario 7: Authentication failure
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "Demo Scenario 7: Authentication failure" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

$authFailXml = @"
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
    $response = Invoke-WebRequest -Uri "http://localhost:8080/soap" -Method Post -ContentType "text/xml" -Body $authFailXml -UseBasicParsing
    Write-Host (Format-XML $response.Content) -ForegroundColor Red
} catch {
    # Authentication error should return 500, which is expected
    if ($null -ne $_.Exception.Response -and $_.Exception.Response.StatusCode -eq 500) {
        Write-Host "Server correctly rejected invalid authentication (Status: 500)" -ForegroundColor Yellow
    } else {
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
}
Write-Host ""

Write-Host "=======================================" -ForegroundColor Cyan
Write-Host "Demo completed successfully!" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "To interact with the service manually:" -ForegroundColor Yellow
Write-Host "  1. Start server: go run cmd/lab3/server/main.go"
Write-Host "  2. Start client: go run cmd/lab3/client/main.go"
Write-Host ""

