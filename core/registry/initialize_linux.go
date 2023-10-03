package registry

import (
	"os/user"
	"path/filepath"
)

// do this but better ig lol
// Maybe do a system/user config map + allow custom dirs?
func (r *Registry) GetPluginDirs() []string {
	thisUser, err := user.Current()
	// We probably don't need to evaluate this here, but differenciating
	// between when an error happens on windows vs linux could be helpful
	if err != nil {
		r.l.Error("couldn't determine user", "error", "user initalize error", "os", "unix")
	}

	pluginDirs := []string{
		`/etc/jc/plugins`,
		filepath.Join(thisUser.HomeDir, `.jc/plugins`),
	}

	return pluginDirs
}
