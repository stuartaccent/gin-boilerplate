package cmd

import (
	"context"
	"log"

	"gin.go.dev/db/dbx"
	"gin.go.dev/internal/crypt"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var (
	setPWEmail, setPWPassword string
)

var cmdSetPassword = &cobra.Command{
	Use:   "setpassword",
	Short: "Set a user's password",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		conn, err := pgx.Connect(ctx, cfg.Database.URL().String())
		if err != nil {
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
		defer conn.Close(ctx)

		hashed, err := crypt.GeneratePassword([]byte(setPWPassword))
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		queries := dbx.New(conn)
		if err = queries.SetUserPasswordByEmail(ctx, dbx.SetUserPasswordByEmailParams{
			Email:          setPWEmail,
			HashedPassword: hashed,
		}); err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		log.Printf("Password set for user: %v\n", setPWEmail)
	},
}

func init() {
	cmdSetPassword.Flags().StringVarP(&setPWEmail, "email", "e", "", "The email address of the user")
	cmdSetPassword.Flags().StringVarP(&setPWPassword, "password", "p", "", "The password of the user")
	_ = cmdSetPassword.MarkFlagRequired("email")
	_ = cmdSetPassword.MarkFlagRequired("password")
}
