package packagemanagers

func init() {
	AddManager("pacman", &Manager{
		commands: ManagerCommands{
			install: "pacman",
			update:  "pacman",
			refresh: "pacman",
			remove:  "pacman",
			clean:   "pacman",
		},
		flags: ManagerFlags{
			install: []string{"-S", "--needed"},
			update:  []string{"-Su"},
			refresh: []string{"-Syy"},
			remove:  []string{"-Rcs"},
			clean:   []string{"-Sc"},
			global:  []string{"--noconfirm"},
		},
	})
}
