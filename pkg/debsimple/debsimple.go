package debsimple

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fsnotify/fsnotify"
)

type Conf struct {
	RootRepoPath  string   `json:"rootRepoPath"`
	SupportArch   []string `json:"supportedArch"`
	Sections      []string `json:"sections"`
	DistroNames   []string `json:"distroNames"`
	EnableSigning bool     `json:"enableSigning"`
	PrivateKey    string   `json:"privateKey"`
}

func (c Conf) ArchPath(distro, section, arch string) string {
	return filepath.Join(c.RootRepoPath, "dists", distro, section, "binary-"+arch)
}

var (
	verbose   *bool
	mywatcher *fsnotify.Watcher

	// Now is a package level time function so we can mock it out
	Now = func() time.Time {
		return time.Now()
	}
)

func Main(parsedconfig Conf, verbosev *bool) {
	var err error
	verbose = verbosev

	// fire up filesystem watcher
	mywatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error creating fswatcher: ", err)
	}

	if err := createDirs(parsedconfig); err != nil {
		log.Println(err)
		log.Fatalf("error creating directory structure, exiting")
	}
}

func CopyDeb(fpath string, config Conf, distro, section, arch string) string {
	newpath := filepath.Join(config.ArchPath(distro, section, arch), filepath.Base(fpath))
	dst, err := os.Create(newpath)
	if err != nil {
		panic(err)
	}
	src, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}

	dst.Close()
	src.Close()
	return newpath
}

func CreateMetadata(name string, parsedconfig Conf) {
	distroArch := destructPath(name)
	if err := createPackagesGz(parsedconfig, distroArch[0], distroArch[1], distroArch[2]); err != nil {
		log.Printf("error creating package: %s", err)
	}
	if parsedconfig.EnableSigning {
		if err := createRelease(parsedconfig, distroArch[0]); err != nil {
			log.Printf("Error creating Release file: %s", err)
		}
	}
}

func GenerateSigningKey(keyName, keyEmail *string) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current working directory: %s", err)
	}
	fmt.Println("Generating new signing key pair..")
	fmt.Printf("Name: %s\n", *keyName)
	fmt.Printf("Email: %s\n", *keyEmail)
	createKeyHandler(workingDirectory, *keyName, *keyEmail)
	fmt.Println("Done.")
}

func destructPath(filePath string) []string {
	splitPath := strings.Split(filePath, "/")
	archFull := splitPath[len(splitPath)-2]
	archSplit := strings.Split(archFull, "-")
	distro := splitPath[len(splitPath)-4]
	section := splitPath[len(splitPath)-3]
	return []string{distro, section, archSplit[1]}
}

func createDirs(config Conf) error {
	for _, distro := range config.DistroNames {
		for _, arch := range config.SupportArch {
			for _, section := range config.Sections {
				if _, err := os.Stat(config.ArchPath(distro, section, arch)); err != nil {
					if os.IsNotExist(err) {
						log.Printf("Directory for %s (%s) does not exist, creating", distro, arch)
						if err := os.MkdirAll(config.ArchPath(distro, section, arch), 0755); err != nil {
							return fmt.Errorf("error creating directory for %s (%s): %s", distro, arch, err)
						}
					} else {
						return fmt.Errorf("error inspecting %s (%s): %s", distro, arch, err)
					}
				}
				log.Println("starting watcher for ", config.ArchPath(distro, section, arch))
				err := mywatcher.Add(config.ArchPath(distro, section, arch))
				if err != nil {
					return fmt.Errorf("error creating watcher for %s (%s): %s", distro, arch, err)
				}
			}
		}
	}
	return nil
}

func openDB() *bolt.DB {
	// open/create database for API keys
	db, err := bolt.Open("debsimple.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("unable to open database: ", err)
	}

	return db
}

func createAPIkey(db *bolt.DB) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	apiKey := base64.URLEncoding.EncodeToString(randomBytes)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("APIkeys"))
		if b == nil {
			return errors.New("Database bucket \"APIkeys\" does not exist")
		}

		err = b.Put([]byte(apiKey), []byte(apiKey))
		return err
	})
	if err != nil {
		return "", err
	}
	return apiKey, nil
}
