package main

import (
	"fmt"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

func main() {
	client := pingdom.NewClient("username", "password", "api_key")

	// List all checks
	checks, _ := client.Checks.List()
	fmt.Println("All checks:", checks)

	// Create a new http check
	newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
	check, _ := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

	// Create a new ping check
	newPingCheck := pingdom.PingCheck{Name: "Test Ping", Hostname: "example.com", Resolution: 1}
	pingcheck, _ := client.Checks.Create(&newPingCheck)
	fmt.Println("Created check:", pingcheck) // {ID, Name}

	// Get details for a check
	details, _ := client.Checks.Read(check.ID)
	fmt.Println("Details:", details)

	// Update a check
	updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	upMsg, _ := client.Checks.Update(check.ID, &updatedCheck)
	fmt.Println("Modified check, message:", upMsg)

	// Delete a check
	delMsg, _ := client.Checks.Delete(check.ID)
	fmt.Println("Deleted check, message:", delMsg)

}
