package main

import (
	"github.com/JaegyuDev/jae-cli/core/registry"
	"github.com/JaegyuDev/jae-cli/core/root"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
	"plugin"
)

var (
	logger  = slog.New(slog.NewTextHandler(os.Stdout, nil))
	vp      = viper.New()
	rootCmd = root.CreateRootCmd()
	reg     = registry.New(logger)
)

func init() {
	// Initialize Viper
	vp.SetConfigName("config.toml")
	vp.AddConfigPath("/etc/jc/")
	vp.AddConfigPath("$HOME/.jc")
	vp.SetConfigType("toml")
	vp.AutomaticEnv()
}

func main() {
	// load from `/etc/jc/plugins` and ~/.jc/plugins

	initializePlugins(reg)

	// load registered plugins
	for _, _plugin := range reg.GetPlugins() {
		logger.Info("Loading plugin", _plugin, _plugin.GetName())
		err := _plugin.Setup(rootCmd, vp, logger)
		if err != nil {
			logger.Warn("Couldn't load plugin", "plugin", _plugin.GetName(), "error", err.Error())
			os.Exit(1)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Couldn't execute command", "error", err.Error())
	}
}

func initializePlugins(_reg *registry.Registry) {
	dirs := _reg.GetPluginDirs()
	logger.Info("checking for plugins")
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Warn("plugins folder doesn't exist. skipping.", "path", dir)
			} else {
				logger.Error("error stating plugin dir", "path", dir, "error", err.Error())
			}
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			logger.Error("couldn't read dir", "path", dir, "error", err.Error())
		}

		for _, entry := range entries {
			var path = filepath.Join(dir, entry.Name())
			if entry.IsDir() || (entry.Type() == os.ModeSymlink) {
				// skip for now, might add dir support later
				logger.Warn("skipping symlink/dir", "path", path)
				continue
			}

			open, err := plugin.Open(path)
			if err != nil {
				logger.Warn("error opening plugin", "path", path, "error", err.Error())
				continue
			}

			// Load the NewPlugin function
			newPluginSym, err := open.Lookup("NewPlugin")
			if err != nil {
				logger.Warn("error looking up 'NewPlugin' symbol in plugin", "path", path, "error", err.Error())
				continue
			}

			// fails here. Probably fucking up when doing the lookup.
			// Cast the loaded symbol to the NewPlugin function
			newPluginFunc, ok := newPluginSym.(func() registry.Plugin)
			if !ok {
				logger.Warn("loaded symbol 'NewPlugin' is not a valid function")
				continue
			}

			// Call the NewPlugin function to create an instance of your plugin
			loadedPlugin := newPluginFunc()

			// Register the loaded plugin
			_reg.RegisterPlugin(loadedPlugin)
		}
	}
}
