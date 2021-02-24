package packagemanagers

func init() {
	AddManager("apt", &Manager{
		commands: ManagerCommands{
			install: "apt",
			update:  "apt",
			upgrade: "apt",
			remove:  "apt",
			clean:   "apt",
		},
		flags: ManagerFlags{
			install: []string{"install"},
			update:  []string{"update"},
			upgrade: []string{"dist-upgrade"},
			remove:  []string{"remove", "--auto-remove"},
			clean:   []string{"clean"},
			global:  []string{"-y"},
		},
	})
}
