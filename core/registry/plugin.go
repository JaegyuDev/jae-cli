package registry

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
)

type Plugin interface {
	Setup(
		cmd *cobra.Command,
		vp *viper.Viper,
		logger *slog.Logger,
	) error
	GetName() string
}
