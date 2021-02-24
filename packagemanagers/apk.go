package packagemanagers

func init() {
	AddManager("apk", &Manager{
		commands: ManagerCommands{
			install: "apk",
			update:  "apk",
			upgrade: "apk",
			remove:  "apk",
			clean:   "apk",
		},
		flags: ManagerFlags{
			install: []string{"add"},
			update:  []string{"update"},
			upgrade: []string{"upgrade"},
			remove:  []string{"del", "--rdepends"},
			global:  []string{"--no-cache"},
		},
	})
}
