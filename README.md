[![Build status](https://badge.buildkite.com/98081178d65983e047299f1eaabeaf336cefc5db496500332c.svg)](https://buildkite.com/datum/go-turso)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=datumforge_go-turso&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=datumforge_go-turso)

# Go-Turso

Golang client for interacting with the Turso Platform API. 

Currently supports the following endpoints:
1. `Organizations`: `List`
1. `Groups`: `List`, `Get`, `Create`, `Delete`
1. `Databases`: `List`, `Get`, `Create`, `Delete`
1. `Database Tokens`: `Create`

## Usage

```
go get github.com/datumforge/go-turso
```

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/datumforge/go-turso"
)

func main() {
	// setup the configuration
	apiToken := os.Getenv("TURSO_API_TOKEN")
	config := turso.Config{
		Token:   apiToken,
		BaseURL: "https://api.turso.tech",
		OrgName: "datum",
	}

	// create the Turso Client
	tc, err := turso.NewClient(config)
	if err != nil {
		log.Fatalf("failed to initialize turso client", "error", err)
	}

	// List All Groups
	groups, err := tc.Group.ListGroups(context.Background())
	if err != nil {
		fmt.Println("error listing groups:", err)

		return
	}

	for _, group := range groups.Groups {
		fmt.Println("Group:", group.Name)
	}

	// List All Databases
	databases, err := tc.Database.ListDatabases(context.Background())
	if err != nil {
		fmt.Println("error listing databases:", err)

		return
	}

	for _, database := range databases.Databases {
		fmt.Println("Database:", database.Name)
	}

	// List All Organizations
	orgs, err := tc.Organization.ListOrganizations(context.Background())
	if err != nil {
		fmt.Println("error listing organizations:", err)

		return
	}

	for _, org := range *orgs {
		fmt.Println("Organization:", org.Name, org.Slug)
	}

	// Create a new Group
	g := turso.CreateGroupRequest{
		Name:     "test-group",
		Location: "ord",
	}

	group, err := tc.Group.CreateGroup(context.Background(), g)
	if err != nil {
		fmt.Println("error creating group:", err)
	}

	fmt.Println("Group Created:", group.Group.Name, group.Group.Locations)

	// Delete the Group
	deletedGroup, err := tc.Group.DeleteGroup(context.Background(), g.Name)
	if err != nil {
		fmt.Println("error deleting group:", err)
	}

	fmt.Println("Group Deleted:", deletedGroup.Group.Name, deletedGroup.Group.Locations)

	// Create a new Database
	d := turso.CreateDatabaseRequest{
		Group:    "default",
		IsSchema: false,
		Name:     "test-database",
	}

	database, err := tc.Database.CreateDatabase(context.Background(), d)
	if err != nil {
		fmt.Println("error creating database:", err)
	}

	fmt.Println("Database Created:", database.Database.Name)

	// Delete the Database
	deletedDatabase, err := tc.Database.DeleteDatabase(context.Background(), d.Name)
	if err != nil {
		fmt.Println("error deleting database:", err)
	}

	fmt.Println("Database Deleted:", deletedDatabase.Database)
}
```

## References

1. [Turso Platform API](https://docs.turso.tech/api-reference/introduction)