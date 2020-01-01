package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "tictl",
	Short: "tictl install the k8s service of Ti platform",
	Long:  "",

	Run:  func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo: " + strings.Join(args, " "))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
