# pingdom-go #

pingdom-go is a Go client library for the Pingdom API.

This currently supports working with basic HTTP (with specific details)
and ping checks. It can also manage contacts who will be alerted in case of downtime.

**Build Status:** [![Build Status](https://travis-ci.org/russellcardullo/go-pingdom.svg?branch=master)](https://travis-ci.org/russellcardullo/go-pingdom)

**Godoc:** https://godoc.org/github.com/russellcardullo/go-pingdom/pingdom

## Usage ##

### Client ###

Pingdom handles single-user and multi-user accounts differently.

Construct a new single-user Pingdom client:

```go
client := pingdom.NewClient("pingdom_username", "pingdom_password", "pingdom_api_key")
```

Construct a multi-user Pingdom client:

```go
client := pingdom.NewMultiUserClient("pingdom_username", "pingdom_password", "pingdom_api_key", "pingdom_account_email")
```

The `pingdom_account_email` variable is the email address of the owner of the multi-user account. This is passed in the `Account-Email` header to the Pingdom API.

Using a Pingdom client, you can access supported services.

### CheckService ###

This service manages pingdom Checks which are represented by the `Check` struct.
When creating or updating Checks you must specify at a minimum the `Name`, `Hostname`
and `Resolution`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Checks from Pingdom: https://www.pingdom.com/features/api/documentation/#ResourceChecks

Get a list of all checks:

```go
checks, err := client.Checks.List()
fmt.Println("Checks:", checks) // [{ID Name} ...]
```

Create a new HTTP check:

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new Ping check:
```go
newCheck := pingdom.PingCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Get details for a specific check:

```go
checkDetails, err := client.Checks.Read(12345)
```

For checks with detailed information, check the specific details in
the field `Type` (e.g. `checkDetails.Type.HTTP`).

Update a check:

```go
updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
msg, err := client.Checks.Update(12345, &updatedCheck)
```

Delete a check:

```go
msg, err := client.Checks.Delete(12345)
```

Create a notification contact:

```go
newContact := pingdom.Contact{ Name: "John Doe", Email: "john.doe@example.com", Defaultsmsprovider: "nexmo" }
contactResponse, err := client.Contacts.Create(&newContact)
```

Create a check with basic alert notification to a contact. Note that you must set ContactIds, UseLegacyNotifications, SendNotificationWhenDown and at least one of the SendTo* parameters:

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5, ContactIds: []int{contactResponse.ID}, UseLegacyNotifications: true, SendNotificationWhenDown: 2, SendToEmail: true}
checkResponse, err := client.Checks.Create(&newCheck)
```
