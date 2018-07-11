package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/michaelnugent/docker-ipv6nat"
)

var (
	buildVersion string

	cleanup       bool
	retry         bool
	userlandProxy bool
	version       bool
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: docker-ipv6 [options]

Automatically configure IPv6 NAT for running docker containers

Options:`)
	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr, `
Environment Variables:
  DOCKER_HOST - default value for -endpoint
  DOCKER_CERT_PATH - directory path containing key.pem, cert.pem and ca.pem
  DOCKER_TLS_VERIFY - enable client TLS verification
`)

	fmt.Fprintln(os.Stderr, `For more information, see https://github.com/robbertkl/docker-ipv6nat`)
}

func initFlags() {
	flag.BoolVar(&cleanup, "cleanup", false, "remove rules when shutting down")
	flag.BoolVar(&retry, "retry", false, "keep retrying to reconnect after a disconnect")
	flag.BoolVar(&version, "version", false, "show version")

	flag.Usage = usage
	flag.Parse()
}

func main() {
	initFlags()

	if version {
		fmt.Println(buildVersion)
		return
	}

	if flag.NArg() > 0 {
		usage()
		os.Exit(1)
	}

	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	state, err := dockeripv6nat.NewState()
	if err != nil {
		return err
	}

	if cleanup {
		defer func() {
			if err := state.Cleanup(); err != nil {
				log.Printf("%v", err)
			}
		}()
	}

	watcher := dockeripv6nat.NewWatcher(client, state, retry)
	if err := watcher.Watch(); err != nil {
		return err
	}

	return nil
}
