package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gin.go.dev/internal/config"
	"gin.go.dev/internal/crypt"
	"gin.go.dev/internal/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [options]")
		fmt.Println("")
		fmt.Println("Available commands:")
		fmt.Println("  createuser")
		return
	}

	command := os.Args[1]
	ctx := context.Background()

	switch command {
	case "createuser":
		commandCreateUser(ctx)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}

// commandCreateUser is the createuser command.
func commandCreateUser(ctx context.Context) {
	cmd := flag.NewFlagSet("createuser", flag.ExitOnError)
	help := cmd.Bool("help", false, "Display help for the createuser command")
	conf := cmd.String("app-config", "config.toml", "The path of the app config eg: config.toml")
	email := cmd.String("email", "", "The email address of the user")
	password := cmd.String("password", "", "The password of the user")
	firstName := cmd.String("firstname", "", "The first name of the user")
	lastName := cmd.String("lastname", "", "The last name of the user")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v\n", err)
	}

	if *help {
		cmd.Usage()
		return
	}
	if *conf == "" {
		flag.Usage()
		log.Fatalln("The -app-config flag is required.")
	}
	if *email == "" {
		flag.Usage()
		log.Fatalln("The -email flag is required.")
	}
	if *password == "" {
		flag.Usage()
		log.Fatalln("The -password flag is required.")
	}
	if *firstName == "" {
		flag.Usage()
		log.Fatalln("The -firstname flag is required.")
	}
	if *lastName == "" {
		flag.Usage()
		log.Fatalln("The -lastname flag is required.")
	}

	cfg, err := config.NewConfigFromPath(*conf)
	if err != nil {
		log.Fatalf("Invalid config: %v\n", err)
	}

	createUser(ctx, cfg, *email, *password, *firstName, *lastName)
}

// createUser creates a new user in the database.
func createUser(ctx context.Context, cfg *config.Config, email, password, firstName, lastName string) {
	conn, err := dbConn(ctx, cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer conn.Close(ctx)

	passwordHash, err := crypt.GeneratePassword([]byte(password))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	queries := db.New(conn)
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:          email,
		HashedPassword: passwordHash,
		FirstName:      firstName,
		LastName:       lastName,
	})
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	log.Printf("User created!: %v\n", user.Email)
}

// dbConn creates a new database connection.
func dbConn(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
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
		return nil, err
	}
	return conn, nil
}
