package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	simpleUtil "github.com/liserjrqlxue/simple-util"

	"github.com/liserjrqlxue/fileServer/router"
)

var (
	port = flag.String(
		"port",
		":9091",
		"web server listen port",
	)
)

func main() {
	flag.Parse()
	err := os.MkdirAll("public", 0755)
	simpleUtil.CheckErr(err)
	err = os.MkdirAll(path.Join("public", "upload"), 0755)
	simpleUtil.CheckErr(err)
	http.HandleFunc("/mp4", router.Mp4play)
	http.HandleFunc("/upload", router.Upload)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Path := r.URL.Path
		urlPath := fmt.Sprintf("%s", Path)
		fmt.Println(Path)
		if strings.HasPrefix(r.URL.Path, "/public") {
			http.ServeFile(w, r, "./"+urlPath)
		} else {
			http.ServeFile(w, r, "./public/"+urlPath)
		}

	}) //设置访问的路由
	fmt.Println("start")
	err = http.ListenAndServe(*port, nil) //设置监听的端口
	simpleUtil.CheckErr(err)
}
