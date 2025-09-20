package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pax",
	Short: "PAX - The Oreon package manager",
	Long: `PAX is the official package manager for the Oreon 11.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PAX - The Oreon package manager")
		fmt.Println("Use 'pax --help' to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
