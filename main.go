package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/liserjrqlxue/goUtil/simpleUtil"

	"github.com/liserjrqlxue/fileServer/router"
)

var (
	cwd, _       = os.Getwd()
	ex, _        = os.Executable()
	exPath       = filepath.Dir(ex)
	templatePath = filepath.Join(exPath, "template")
)

var (
	port = flag.String(
		"port",
		":9091",
		"web server listen port",
	)
	public = flag.String(
		"public",
		"",
		"root path of public, default is current workdir",
	)
)

var err error

func main() {
	simpleUtil.CheckErr(err)

	flag.Parse()
	if *public == "" {
		*public = cwd
	}

	router.PublicPath = *public
	router.TemplatePath = templatePath
	http.HandleFunc("/mp4", router.Mp4play)
	http.HandleFunc("/upload", router.Upload)
	http.HandleFunc("/", router.Download)

	log.Println("Start", "http://localhost"+*port)
	simpleUtil.CheckErr(http.ListenAndServe(*port, nil)) //设置监听的端口
}
