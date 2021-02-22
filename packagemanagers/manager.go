package packagemanagers

import (
	"fmt"
	"strings"
)

var (
	managers = map[string]*Manager{}
)

// ManagerCommands represents all commands.
type ManagerCommands struct {
	install string
	update  string
	refresh string
	remove  string
	clean   string
}

// ManagerFlags represents flags for all subcommands of a package manager.
type ManagerFlags struct {
	install []string
	update  []string
	refresh []string
	remove  []string
	clean   []string
	global  []string
}

// A Manager represents a package manager.
type Manager struct {
	Comment  string
	commands ManagerCommands
	flags    ManagerFlags
}

// Install installs packages to the system.
func (m Manager) Install(packages []string, flags []string) string {
	if len(m.flags.install) == 0 || packages == nil || len(packages) == 0 {
		return ""
	}

	args := append(m.flags.global, m.flags.install...)
	args = append(args, flags...)
	args = append(args, packages...)

	return fmt.Sprintf("%s %s", m.commands.install, strings.Join(args, " "))
}

// Update updates all packages.
func (m Manager) Update() string {
	if len(m.flags.update) == 0 {
		return ""
	}

	args := append(m.flags.global, m.flags.update...)

	return fmt.Sprintf("%s %s", m.commands.update, strings.Join(args, " "))
}

// Refresh refreshes the package database.
func (m Manager) Refresh() string {
	if len(m.flags.refresh) == 0 {
		return ""
	}

	args := append(m.flags.global, m.flags.refresh...)

	return fmt.Sprintf("%s %s", m.commands.refresh, strings.Join(args, " "))
}

// Remove removes packages from the system.
func (m Manager) Remove(packages []string, flags []string) string {
	if len(m.flags.remove) == 0 || packages == nil || len(packages) == 0 {
		return ""
	}

	args := append(m.flags.global, m.flags.remove...)
	args = append(args, flags...)
	args = append(args, packages...)

	return fmt.Sprintf("%s %s", m.commands.remove, strings.Join(args, " "))
}

// Clean cleans up cached files used by the package managers.
func (m Manager) Clean() string {
	if len(m.flags.clean) == 0 {
		return ""
	}

	args := append(m.flags.global, m.flags.clean...)

	return fmt.Sprintf("%s %s", m.commands.clean, strings.Join(args, " "))
}

// AddManager add manager.
func AddManager(name string, manager *Manager) bool {
	if ExistsManager(name) {
		return false
	}

	managers[name] = manager
	return true
}

// DeleteManager delete manager.
func DeleteManager(name string) bool {
	if !ExistsManager(name) {
		return false
	}

	delete(managers, name)
	return true
}

// ExistsManager is manager exists
func ExistsManager(name string) bool {
	_, ok := managers[name]
	return ok
}

// GetManager returns a Manager specified by name.
func GetManager(name string) *Manager {
	if !ExistsManager(name) {
		return nil
	}

	return managers[name]
}

// Managers returns all Managers.
func Managers() map[string]*Manager {
	return managers
}
