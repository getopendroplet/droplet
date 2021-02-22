package packagemanagers

func init() {
	AddManager("zypper", &Manager{
		commands: ManagerCommands{
			install: "zypper",
			update:  "zypper",
			refresh: "zypper",
			remove:  "zypper",
			clean:   "zypper",
		},
		flags: ManagerFlags{
			install: []string{"install", "--allow-downgrade"},
			update:  []string{"update"},
			refresh: []string{"refresh"},
			remove:  []string{"remove"},
			clean:   []string{"clean", "-a"},
			global:  []string{"--non-interactive", "--gpg-auto-import-keys"},
		},
	})
}
