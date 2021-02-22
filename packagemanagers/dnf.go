package packagemanagers

func init() {
	AddManager("dnf", &Manager{
		commands: ManagerCommands{
			install: "dnf",
			update:  "dnf",
			refresh: "dnf",
			remove:  "dnf",
			clean:   "dnf",
		},
		flags: ManagerFlags{
			install: []string{"install"},
			update:  []string{"upgrade"},
			refresh: []string{"makecache"},
			remove:  []string{"remove"},
			clean:   []string{"clean", "all"},
			global:  []string{"-y"},
		},
	})
}
