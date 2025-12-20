package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Lab 3: SOAP Web Service for File Management")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("This lab implements a SOAP-based file upload and management system.")
	fmt.Println()
	fmt.Println("Components:")
	fmt.Println("  1. Server - SOAP web service (cmd/lab3/server)")
	fmt.Println("  2. Client - Interactive client application (cmd/lab3/client)")
	fmt.Println()
	fmt.Println("To run the server:")
	fmt.Println("  go run cmd/lab3/server/main.go")
	fmt.Println()
	fmt.Println("To run the client:")
	fmt.Println("  go run cmd/lab3/client/main.go")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("  - WSDL-based service description")
	fmt.Println("  - User authentication")
	fmt.Println("  - File upload with MTOM support (up to 3MB)")
	fmt.Println("  - Asynchronous upload notifications")
	fmt.Println("  - File validation (size, empty, forbidden chars, JSON content)")
	fmt.Println("  - Query methods: last file info, file list CSV, server uptime")
	fmt.Println()
	fmt.Println("Authentication credentials:")
	fmt.Println("  user1:pass1, user2:pass2, admin:admin")
	fmt.Println()

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("For more details, see cmd/lab3/README.md")
	}
}
