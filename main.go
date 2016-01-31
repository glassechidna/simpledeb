package main

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type semaphore chan int

type Conf struct {
	ListenPort      string   `json:"listenPort"`
	RootRepoPath    string   `json:"rootRepoPath"`
	SupportArch     []string `json:"supportedArch"`
	ScanpackagePath string   `json:"scanpackagePath"`
}

var sem = make(semaphore, 1)
var configFile = flag.String("c", "conf.json", "config file location")
var config = &Conf{}

func main() {
	flag.Parse()
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic("unable to read config file, exiting...")
	}
	config = &Conf{}
	if err := json.Unmarshal(file, &config); err != nil {
		log.Panic("unable to marshal config file, exiting...")
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.RootRepoPath))))
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":"+config.ListenPort, nil)
}
func createPackagesTar(arch string) bool {
	cmd1 := exec.Command(config.ScanpackagePath, config.RootRepoPath+"/dists/stable/main/binary-"+arch, "/dev/null")
	outfile, err := os.Create(config.RootRepoPath + "/dists/stable/main/binary-" + arch + "/Packages.gz")
	if err != nil {
		log.Println("error creating packages.gz file")
		return false
	}
	defer outfile.Close()
	gzOut := gzip.NewWriter(outfile)
	defer gzOut.Close()
	cmd1.Stdout = gzOut
	err = cmd1.Run()
	if err != nil {
		log.Println("unable run scanpackages ", err)
		return false
	}
	gzOut.Flush()
	return true
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println("error parsing file: ", err)
		}
		log.Println("filename is: ", header.Filename)
		out, err := os.Create(config.RootRepoPath + "/dists/stable/main/binary-amd64/" + header.Filename)
		if err != nil {
			log.Println("error creating file: ", err)
		}
		hash := md5.New()
		_, err = io.Copy(out, io.TeeReader(file, hash))
		defer file.Close()
		if err != nil {
			log.Println("error saving file: ", err)
		}
		md5sum := hex.EncodeToString(hash.Sum(nil))
		log.Println("md5sum = ", md5sum)
		log.Println("grabbing lock...")
		sem.Lock()
		log.Println("got lock, updating package list...")
		if !createPackagesTar("amd64") {
			log.Println("unable to create Packages.gz")
		}
		sem.Unlock()
		log.Println("lock returned")
	} else {
		log.Println("not a POST")
	}
}

func (s semaphore) Lock() {
	s <- 1
}

func (s semaphore) Unlock() {
	<-s
}