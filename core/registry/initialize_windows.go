package registry

import (
	"os/user"
	"path/filepath"
)

// do this but better ig lol
func (r *Registry) GetPluginDirs() []string {
	thisUser, err := user.Current()
	// We probably don't need to evaluate this here, but differenciating
	// between when an error happens on windows vs linux could be helpful
	if err != nil {
		r.l.Error("couldn't determine user", "error", "user initalize error", "os", "windows")
	}

	pluginDirs := []string{
		`C:\Program Files\jc\plugins`,
		filepath.Join(thisUser.HomeDir, `.jc\plugins`),
	}

	return pluginDirs
}
