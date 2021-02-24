package packagemanagers

func init() {
	AddManager("pacman", &Manager{
		commands: ManagerCommands{
			install: "pacman",
			update:  "pacman",
			upgrade: "pacman",
			remove:  "pacman",
			clean:   "pacman",
		},
		flags: ManagerFlags{
			install: []string{"-S", "--needed"},
			update:  []string{"-Syy"},
			upgrade: []string{"-Su"},
			remove:  []string{"-Rcs"},
			clean:   []string{"-Sc"},
			global:  []string{"--noconfirm"},
		},
	})
}
