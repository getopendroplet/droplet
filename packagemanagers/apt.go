package packagemanagers

func init() {
	AddManager("apt", &Manager{
		commands: ManagerCommands{
			install: "apt-get",
			update:  "apt-get",
			refresh: "apt-get",
			remove:  "apt-get",
			clean:   "apt-get",
		},
		flags: ManagerFlags{
			install: []string{"add"},
			update:  []string{"dist-upgrade"},
			refresh: []string{"update"},
			remove:  []string{"remove", "--auto-remove"},
			clean:   []string{"clean"},
			global:  []string{"--no-cache"},
		},
	})
}
