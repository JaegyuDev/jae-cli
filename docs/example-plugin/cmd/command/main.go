package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
)

type Plugin struct {
	name string
}

func NewPlugin() Plugin {
	return Plugin{
		"Demo Plugin",
	}
}

func (p *Plugin) Setup(cmd *cobra.Command, vp *viper.Viper, l *slog.Logger) {
	myCmd := &cobra.Command{
		Use:   "example",
		Short: "example demo",
		Run: func(cmd *cobra.Command, args []string) {
			// logic here.
			// my only request is that you don't modify the base configs.
			// you can add your own inside of the programs config paths, which will be provided through viper.

			fmt.Println("HELLO WORLD!!")
		},
	}

	cmd.AddCommand(myCmd)
}

func (p *Plugin) GetName() string {
	return p.name
}

// dont use this, it will be ignored and nothing will happen
func main() {}
