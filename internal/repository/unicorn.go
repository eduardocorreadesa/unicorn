package repository

import (
	"bufio"
	"log"
	"os"
	"unicorn/internal/configuration"
)

type UnicornRepo struct {
	UnicornNames        []string
	UnicornAdj          []string
	UnicornCapabilities []string
}

func UnicornData(config configuration.Config) *UnicornRepo {

	unicornNames, err := os.Open(config.Server.PathPetNames)
	if err != nil {
		log.Fatal("Fail to open petnames.txt", err)
		panic(err)
	}

	var names []string
	var scanner = bufio.NewScanner(unicornNames)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}

	unicornAdj, err := os.Open(config.Server.PathPetAdj)
	if err != nil {
		log.Fatal("Fail to open adj.txt", err)
		panic(err)
	}

	var adj []string
	scanner = bufio.NewScanner(unicornAdj)
	for scanner.Scan() {
		adj = append(adj, scanner.Text())
	}

	unicornCapababilities, err := os.Open(config.Server.PathPetCapabilities)
	if err != nil {
		log.Fatal("Fail to open capabilities.txt", err)
		panic(err)
	}

	var cap []string
	scanner = bufio.NewScanner(unicornCapababilities)
	for scanner.Scan() {
		cap = append(cap, scanner.Text())
	}

	return &UnicornRepo{
		UnicornNames:        names,
		UnicornAdj:          adj,
		UnicornCapabilities: cap,
	}
}
