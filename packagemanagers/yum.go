package packagemanagers

func init() {
	AddManager("yum", &Manager{
		commands: ManagerCommands{
			install: "yum",
			update:  "yum",
			refresh: "yum",
			remove:  "yum",
			clean:   "yum",
		},
		flags: ManagerFlags{
			install: []string{"add"},
			update:  []string{"update"},
			refresh: []string{"makecache"},
			remove:  []string{"remove"},
			clean:   []string{"clean", "all"},
			global:  []string{"-y"},
		},
	})
}
