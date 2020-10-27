package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/liserjrqlxue/goUtil/simpleUtil"

	"github.com/liserjrqlxue/fileServer/router"
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
	flag.Parse()
	if *public == "" {
		*public, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		print(*public)
	}
	http.HandleFunc("/mp4", router.Mp4play)
	http.HandleFunc("/upload", router.Upload)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var path = *public + r.URL.Path
		fmt.Printf("[%s] -> [%s]\n", r.URL.Path, path)
		http.ServeFile(w, r, path)
	}) //设置访问的路由
	fmt.Println("start", "http://localhost"+*port)
	simpleUtil.CheckErr(http.ListenAndServe(*port, nil)) //设置监听的端口
}
