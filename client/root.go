/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"main/protofile"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var addUserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "add user information",
	Long:  `client adduser <name> <email> <phone>`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(getTimeStamp())
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		handleError(err)
		defer conn.Close()
		c := protofile.NewUserServiceClient(conn)
		phone, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		UserData(c, args[0], args[1], phone)
	},
}
var addActCmd = &cobra.Command{
	Use:   "addActivity",
	Short: "add activities done by user",
	Long:  `client addActivity <email> <activitytype> <duration> <label>`,

	Run: func(cmd *cobra.Command, args []string) {

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		handleError(err)
		defer conn.Close()
		c := protofile.NewUserServiceClient(conn)
		duration, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		ActData(c, args[0], args[1], int32(duration), args[3])
	},
}
var getuserCmd = &cobra.Command{
	Use:   "getuser",
	Short: "To get the user details",
	Long: `To get the user details (name, email, phone) by taking email.
syntax:
	client getuser --email=<email>`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		handleError(err)
		defer conn.Close()
		c := protofile.NewUserServiceClient(conn)
		email, err2 := cmd.Flags().GetString("email")
		handleError(err2)
		GetUser(c, email)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var getactCmd = &cobra.Command{
	Use:   "getact",
	Short: "To get activity data of a user",
	Long: `To get activity data of a specific user.
syntax:
	client getact --email=<email>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		handleError(err)
		defer conn.Close()
		c := protofile.NewUserServiceClient(conn)
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			log.Fatal(err)
		}
		GetActivity(c, email)
	},
}
var remuserCmd = &cobra.Command{
	Use:   "remuser",
	Short: "To delete an existing user.",
	Long: `To delete an existing user by taking email.
	
Input:
	email
Example:
	client remuser --email=<email>
`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		handleError(err)
		defer conn.Close()
		c := protofile.NewUserServiceClient(conn)
		email, err := cmd.Flags().GetString("email")
		handleError(err)
		RemoveUser(c, email)

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getuserCmd.PersistentFlags().String("email", "", "User email")
	remuserCmd.PersistentFlags().String("email", "", "User email")
	getactCmd.PersistentFlags().String("email", "", "User email")
	rootCmd.AddCommand(addUserCmd)
	rootCmd.AddCommand(addActCmd)
	rootCmd.AddCommand(getuserCmd)
	rootCmd.AddCommand(getactCmd)
	rootCmd.AddCommand(remuserCmd)
}
