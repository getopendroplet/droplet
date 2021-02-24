package packagemanagers

func init() {
	AddManager("yum", &Manager{
		commands: ManagerCommands{
			install: "yum",
			update:  "yum",
			upgrade: "yum",
			remove:  "yum",
			clean:   "yum",
		},
		flags: ManagerFlags{
			install: []string{"add"},
			update:  []string{"update"},
			upgrade: []string{"makecache"},
			remove:  []string{"remove"},
			clean:   []string{"clean", "all"},
			global:  []string{"-y"},
		},
	})
}
