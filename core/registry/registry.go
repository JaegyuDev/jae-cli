package registry

import "log/slog"

type Registry struct {
	l       *slog.Logger
	plugins []Plugin
}

func New(logger *slog.Logger) *Registry {
	return &Registry{
		l:       logger,
		plugins: make([]Plugin, 0),
	}
}

func (r *Registry) RegisterPlugin(plugin Plugin) {
	r.plugins = append(r.plugins, plugin)
}

func (r *Registry) GetPlugins() []Plugin {
	return r.plugins
}
