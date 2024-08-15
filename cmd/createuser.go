package cmd

import (
	"context"
	"fmt"
	"gin.go.dev/internal/crypt"
	"gin.go.dev/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"log"
)

var (
	createEmail, createPassword, createFirstName, createLastName string
)

var cmdCreateUser = &cobra.Command{
	Use:   "createuser",
	Short: "Creates a new user",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil {
			log.Fatalf("Config not initialized")
		}

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
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
		defer conn.Close(ctx)

		passwordHash, err := crypt.GeneratePassword([]byte(createPassword))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		queries := db.New(conn)
		user, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:          createEmail,
			HashedPassword: passwordHash,
			FirstName:      createFirstName,
			LastName:       createLastName,
		})
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		log.Printf("User created!: %v\n", user.Email)
	},
}

func init() {
	cmdCreateUser.Flags().StringVarP(&createEmail, "email", "e", "", "The email address of the user")
	cmdCreateUser.Flags().StringVarP(&createPassword, "password", "p", "", "The password of the user")
	cmdCreateUser.Flags().StringVarP(&createFirstName, "firstname", "f", "", "The first name of the user")
	cmdCreateUser.Flags().StringVarP(&createLastName, "lastname", "l", "", "The last name of the user")
	_ = cmdCreateUser.MarkFlagRequired("email")
	_ = cmdCreateUser.MarkFlagRequired("password")
	_ = cmdCreateUser.MarkFlagRequired("firstname")
	_ = cmdCreateUser.MarkFlagRequired("lastname")
}
