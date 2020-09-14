package gocache

import (
	"fmt"
	ch "go-cache/gocache/consistenthash"
	"go-cache/gocache/pb"
	"google.golang.org/protobuf/proto"
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

// 这是什么抽象？
type HTTPPool struct {
	self      string
	basePath  string
	mu        sync.Mutex
	peers     *ch.Map
	httpPeers map[string]*httpPeerGetter
}

func (h *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if pg, ok := h.httpPeers[key]; ok {
		return pg, true
	}
	return
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
	h.httpPeers = make(map[string]*httpPeerGetter, len(peers))
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

var _ PeerPicker = (*HTTPPool)(nil)

type httpPeerGetter struct {
	baseUrl string
}

// 通过 HTTP 方式获取数据
func (h *httpPeerGetter) Get(req *pb.Request, resp *pb.Response) (err error) {
	requestURL := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(req.GetGroup()), url.QueryEscape(req.GetKey()))
	fmt.Printf("[slowdb]get key:%s from group:%s", req.GetKey(), req.GetGroup())
	resp1, err := http.Get(requestURL)
	if err != nil {
		return err
	}
	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf("server response: %s", resp1.Status)
	}

	respData, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return err
	}
	if err = proto.Unmarshal(respData, resp); err != nil {
		return err
	}
	return nil
}

var _ PeerGetter = (*httpPeerGetter)(nil)
