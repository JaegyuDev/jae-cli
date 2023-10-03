package main

import (
	"github.com/JaegyuDev/jae-cli/core/registry"
	"github.com/JaegyuDev/jae-cli/core/root"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"unsafe"
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

			if plugin, ok := loadPlugin(path); ok {
				_reg.RegisterPlugin(plugin)

			}
		}
	}
}

func loadPlugin(path string) (registry.Plugin, bool) {
	lib, err := syscall.LoadLibrary(path)
	if err != nil {
		logger.Info("Failed to load file as a library", "path", path, "error", err)
	}
	defer syscall.FreeLibrary(lib)

	// we need to load the symbol for the plugin struct
	symPlugin, err := syscall.GetProcAddress(lib, "NewPlugin")
	if err != nil {
		logger.Info("Failed to load NewPlugin symbol from the library", "path", path, "error", err)
	}

	// we convert the symbol to a function pointer (?)
	pluginPtr := *(*unsafe.Pointer)(unsafe.Pointer(&symPlugin))
	newPluginFuncType := reflect.FuncOf([]reflect.Type{}, []reflect.Type{reflect.TypeOf((*registry.Plugin)(nil)).Elem()}, false)

	// This is where i want to change the pluginPtr to a function I can call:
	newPluginFunc := reflect.NewAt(newPluginFuncType, pluginPtr)

	// Call NewPlugin function to create an instance of Plugin
	pluginInstance := newPluginFunc.Call(nil)

	if plugin, ok := pluginInstance[0].Interface().(registry.Plugin); ok {
		return plugin, true
	}

	return nil, false
}

//
//// check if it implements the interface
//if funcValue.Type().Implements(reflect.TypeOf((*registry.Plugin)(nil)).Elem()) {
//return funcValue.Interface().(registry.Plugin), true
//}
