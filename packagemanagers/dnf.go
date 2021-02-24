package packagemanagers

func init() {
	AddManager("dnf", &Manager{
		commands: ManagerCommands{
			install: "dnf",
			update:  "dnf",
			upgrade: "dnf",
			remove:  "dnf",
			clean:   "dnf",
		},
		flags: ManagerFlags{
			install: []string{"install"},
			update:  []string{"makecache"},
			upgrade: []string{"upgrade"},
			remove:  []string{"remove"},
			clean:   []string{"clean", "all"},
			global:  []string{"-y"},
		},
	})
}
