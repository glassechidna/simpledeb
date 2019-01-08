package debsimple

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	// Now is a package level time function so we can mock it out
	Now = func() time.Time {
		return time.Now()
	}
)

func Main(parsedconfig Conf) {
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
						if err := os.MkdirAll(config.ArchPath(distro, section, arch), 0755); err != nil {
							return fmt.Errorf("error creating directory for %s (%s): %s", distro, arch, err)
						}
					} else {
						return fmt.Errorf("error inspecting %s (%s): %s", distro, arch, err)
					}
				}
			}
		}
	}
	return nil
}
