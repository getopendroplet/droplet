package packagemanagers

func init() {
	AddManager("apk", &Manager{
		commands: ManagerCommands{
			install: "apk",
			update:  "apk",
			refresh: "apk",
			remove:  "apk",
			clean:   "apk",
		},
		flags: ManagerFlags{
			install: []string{"add"},
			update:  []string{"upgrade"},
			refresh: []string{"update"},
			remove:  []string{"del", "--rdepends"},
			global:  []string{"--no-cache"},
		},
	})
}
