package cmd

import (
	"context"
	"fmt"
	"gin.go.dev/pkg/auth"
	"gin.go.dev/pkg/storage/db/dbx"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"os"
)

var (
	createEmail, createPassword, createFirstName, createLastName string
)

var cmdCreateUser = &cobra.Command{
	Use:   "createuser",
	Short: "Creates a new user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		conn, err := pgx.Connect(ctx, cfg.Database.URL().String())
		if err != nil {
			fmt.Printf("Error connecting to the database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(ctx)

		passwordHash, err := auth.GeneratePassword([]byte(createPassword))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		queries := dbx.New(conn)
		user, err := queries.CreateUser(ctx, dbx.CreateUserParams{
			Email:          createEmail,
			HashedPassword: passwordHash,
			FirstName:      createFirstName,
			LastName:       createLastName,
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("User created!: %v\n", user.Email)
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
