package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/webx"
	"github.com/jackc/pgx/v5"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [options]")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  createuser")
		fmt.Println("  hexauthkey")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createuser":
		createCmd := flag.NewFlagSet("createuser", flag.ExitOnError)
		createDb := createCmd.String("db", "", "The database DNS: eg: postgres://postgres:password@localhost:5432/gin-boilerplate?sslmode=disable")
		createEmail := createCmd.String("email", "", "The email address of the user")
		createPassword := createCmd.String("password", "", "The password of the user")
		createFName := createCmd.String("firstname", "", "The first name of the user")
		createLName := createCmd.String("lastname", "", "The last name of the user")
		createHelp := createCmd.Bool("help", false, "Display help for the createuser command")

		createCmd.Parse(os.Args[2:])

		if *createHelp {
			createCmd.Usage()
			return
		}
		if *createDb == "" {
			fmt.Println("The -db flag is required.")
			flag.Usage()
			os.Exit(1)
		}
		if *createEmail == "" {
			fmt.Println("The -email flag is required.")
			flag.Usage()
			os.Exit(1)
		}
		if *createPassword == "" {
			fmt.Println("The -password flag is required.")
			flag.Usage()
			os.Exit(1)
		}
		if *createFName == "" {
			fmt.Println("The -firstname flag is required.")
			flag.Usage()
			os.Exit(1)
		}
		if *createLName == "" {
			fmt.Println("The -lastname flag is required.")
			flag.Usage()
			os.Exit(1)
		}

		createUser(*createDb, *createEmail, *createPassword, *createFName, *createLName)

	case "hexauthkey":
		hexAuthKeyCmd := flag.NewFlagSet("hexauthkey", flag.ExitOnError)
		hexAuthKeyCmdLength := hexAuthKeyCmd.Int("length", 0, "The length of the random key.")
		hexAuthKeyHelp := hexAuthKeyCmd.Bool("help", false, "Display help for the hexauthkey command")

		hexAuthKeyCmd.Parse(os.Args[2:])

		if *hexAuthKeyHelp {
			hexAuthKeyCmd.Usage()
			return
		}

		if *hexAuthKeyCmdLength == 0 {
			fmt.Println("The -length flag is required.")
			flag.Usage()
			os.Exit(1)
		}
		hexAuthKey(*hexAuthKeyCmdLength)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func createUser(dns, email, password, firstName, lastName string) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	passwordHash, err := webx.GeneratePassword([]byte(password))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:          email,
		HashedPassword: passwordHash,
		FirstName:      firstName,
		LastName:       lastName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User created!: %v\n", user.Email)
}

func hexAuthKey(length int) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Random key: %s\n", hex.EncodeToString(k))
}
