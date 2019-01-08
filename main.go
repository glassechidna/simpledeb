package main

import (
	"encoding/json"
	"flag"
	"github.com/esell/deb-simple/pkg/debsimple"
	"io/ioutil"
	"log"
	"os"
)

var (
	configFile         = flag.String("c", "conf.json", "config file location")
	generateSigningKey = flag.Bool("k", false, "Generate a signing key pair")
	keyName            = flag.String("kn", "", "Name for the siging key")
	keyEmail           = flag.String("ke", "", "Email address")
	verbose            = flag.Bool("v", false, "Print verbose logs")
	parsedconfig       = debsimple.Conf{}
)

func populateDefaultConfig() {
	dir, _ := os.Getwd()

	parsedconfig = debsimple.Conf{
		RootRepoPath: dir,
		SupportArch: []string{"i386", "amd64"},
		DistroNames: []string{"stable"},
		Sections: []string{"main"},
	}
}

func main() {
	populateDefaultConfig()

	flag.Parse()
	if _, err := os.Stat(*configFile); !os.IsNotExist(err) {
		file, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatal("unable to read config file, exiting...")
		}
		if err := json.Unmarshal(file, &parsedconfig); err != nil {
			log.Fatal("unable to marshal config file, exiting...")
		}
	}

	if *generateSigningKey {
		debsimple.GenerateSigningKey(keyName, keyEmail)
		os.Exit(0)
	}

	debsimple.Main(parsedconfig, verbose)
	newpath := debsimple.CopyDeb("/Users/aidan.steele/Downloads/sshst_v0.0.5-next_linux_amd64.deb", parsedconfig, "stable", "main", "amd64")
	debsimple.CreateMetadata(newpath, parsedconfig)
}
