package main

import (
	"flag"
	"fmt"
	"net/http"

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
		".",
		"root path of public",
	)
)

func main() {
	flag.Parse()
	http.HandleFunc("/mp4", router.Mp4play)
	http.HandleFunc("/upload", router.Upload)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Path := r.URL.Path
		var urlPath = fmt.Sprintf("%s", Path)
		fmt.Println(Path + "\t" + urlPath)
		http.ServeFile(w, r, *public+urlPath)
	}) //设置访问的路由
	fmt.Println("start")
	simpleUtil.CheckErr(http.ListenAndServe(*port, nil)) //设置监听的端口
}
