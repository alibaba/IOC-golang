package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "iocli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}
