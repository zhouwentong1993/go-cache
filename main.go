package main

import (
	"flag"
	"fmt"
	. "go-cache/gocache"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	var port int
	var isApi bool

	flag.IntVar(&port, "port", 8000, "cache server port")
	flag.BoolVar(&isApi, "api", false, "this server is API server?")
	flag.Parse()

	var db = map[string]string{
		"zhangsan": "张三",
		"lisi":     "李四",
		"wangwu":   "王五",
	}

	apiServerURL := "http://localhost:9099"
	cacheServerURLs := []string{
		"http://localhost:8000",
		"http://localhost:8001",
		"http://localhost:8002",
	}

	g := newGroup(db)
	if isApi {
		startAPIServer(apiServerURL, g)
	}
	startCacheServer(getCacheServerURL(cacheServerURLs, port), cacheServerURLs, g)
}

func getCacheServerURL(urls []string, port int) string {
	for _, url := range urls {
		if strings.Contains(url, strconv.Itoa(port)) {
			return url
		}
	}
	return "http://localhost:8000"
}

func newGroup(db map[string]string) *Group {
	return New("score", GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; !ok {
			return nil, fmt.Errorf("no value")
		} else {
			return []byte(v), nil
		}
	}), math.MaxUint64)
}

func startCacheServer(addr string, addrs []string, g *Group) {
	peers := NewHTTPPool(addr)
	peers.Set(addrs...)
	g.RegisterPeer(peers)
	log.Printf("cache server run at %s\n", addr)
	http.ListenAndServe(addr[7:], peers)
}

func startAPIServer(url string, g *Group) {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("api server recevice query %s\n", url)
		key := r.URL.Query().Get("key")
		value, err := g.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(value.ByteSlice())
	})
	log.Printf("api server start at %s\n", url)
	http.ListenAndServe(url[7:], nil)
}
