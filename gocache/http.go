package gocache

import (
	"fmt"
	ch "go-cache/gocache/consistenthash"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultPath    = "/_geecache/"
	defaultReplica = 10
)

type HTTPPool struct {
	self      string
	basePath  string
	mu        sync.Mutex
	peers     *ch.Map
	httpPeers map[string]PeerGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultPath,
	}
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = ch.New(defaultReplica, crc32.ChecksumIEEE)
	// Golang 要这么写，跟 Java 不一样。类型转换得不太对
	h.peers.AddNodes(peers...)
	for _, peer := range peers {
		h.httpPeers[peer] = &httpPeerGetter{baseUrl: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeerGetter(key string) (pg PeerGetter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.GetNode(key); peer != "" && peer != h.self {
		return h.httpPeers[peer], true
	} else {
		return
	}
}

// Log info with server name
func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log("%s %s", r.Method, r.URL)

	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}
	if v, err := group.Get(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(v.bytes)
	}

}

type httpPeerGetter struct {
	baseUrl string
}

// 通过 HTTP 方式获取数据
func (h *httpPeerGetter) Get(group string, key string) (data []byte, err error) {
	requestURL := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server response: %s", resp.Status)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
