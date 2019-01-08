package main

import (
	"encoding/json"
	"flag"
	"github.com/boltdb/bolt"
	"github.com/esell/deb-simple/pkg/debsimple"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	mutex              sync.Mutex
	configFile         = flag.String("c", "conf.json", "config file location")
	generateKey        = flag.Bool("g", false, "generate an API key")
	generateSigningKey = flag.Bool("k", false, "Generate a signing key pair")
	keyName            = flag.String("kn", "", "Name for the siging key")
	keyEmail           = flag.String("ke", "", "Email address")
	verbose            = flag.Bool("v", false, "Print verbose logs")
	parsedconfig       = debsimple.Conf{}
)

func populateDefaultConfig() {
	dir, _ := os.Getwd()

	parsedconfig = debsimple.Conf{
		ListenPort: "9090",
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

	var db *bolt.DB
	if parsedconfig.EnableAPIKeys || *generateKey {
		db = debsimple.CreateDb()
	}

	// generate API key and exit
	if *generateKey {
		debsimple.GenerateKey(db)
		os.Exit(0)
	}

	if *generateSigningKey {
		debsimple.GenerateSigningKey(keyName, keyEmail)
		os.Exit(0)
	}

	debsimple.Main(parsedconfig, verbose)
	go debsimple.KeepMetadataUpdated(mutex, verbose, parsedconfig)
	debsimple.ServeWeb(parsedconfig, db)
}
