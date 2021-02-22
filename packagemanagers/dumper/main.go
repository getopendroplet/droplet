package main

import (
	"droplet/packagemanagers"
	"fmt"
)

func main() {
	for k, v := range packagemanagers.Managers() {
		fmt.Println(k)
		fmt.Println("\tInstall:", v.Install([]string{"a", "b", "c", "d"}, []string{"-y"}))
		fmt.Println("\tUpdate:", v.Update())
		fmt.Println("\tRefresh:", v.Refresh())
		fmt.Println("\tRemove:", v.Remove([]string{"a", "b", "c", "d"}, []string{"-y"}))
		fmt.Println("\tClean:", v.Clean())
	}
}
