package main

import "github.com/spf13/cobra"

func rootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "url-shortener-api",
		Short: "REST API for shortening urls.",
		Run: func(cmd *cobra.Command, args []string) {
			serveHTTP()
		},
	}

	rootCmd.AddCommand(migrateCmd())

	return rootCmd
}
