package packagemanagers

func init() {
	AddManager("brew", &Manager{
		commands: ManagerCommands{
			install: "brew",
			update:  "brew",
			upgrade: "brew",
			remove:  "brew",
			clean:   "brew",
		},
		flags: ManagerFlags{
			install: []string{"install"},
			update:  []string{"update"},
			upgrade: []string{"upgrade"},
			remove:  []string{"remove"},
			clean:   []string{"cleanup"},
			global:  []string{"-f"},
		},
	})
}
