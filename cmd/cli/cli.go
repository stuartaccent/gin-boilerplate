package main

import (
	"context"
	"flag"
	"fmt"
	"gin.go.dev/internal/config"
	"gin.go.dev/internal/db"
	"gin.go.dev/internal/webx"
	"github.com/jackc/pgx/v5"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [options]")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  createuser")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createuser":
		createCmd := flag.NewFlagSet("createuser", flag.ExitOnError)
		createConfig := createCmd.String("app-config", "config.toml", "The path of the app config eg: config.toml")
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
		if *createConfig == "" {
			fmt.Println("The -app-config flag is required.")
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

		cfg, err := config.NewConfigFromPath(*createConfig)
		if err != nil {
			panic(err)
		}

		createUser(cfg, *createEmail, *createPassword, *createFName, *createLName)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func createUser(cfg *config.Config, email, password, firstName, lastName string) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"host=%s port=%v user=%s password=%s database=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Db,
		cfg.Database.SslMode,
	))
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
