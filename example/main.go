package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/groupcache"
)

var (
	peers_addrs = []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}
	rpc_addrs   = []string{"127.0.0.1:9001", "127.0.0.1:9002", "127.0.0.1:9003"}
	index       = flag.Int("index", 0, "peer index")
)

func main() {
	flag.Parse()
	peers_addrs := make([]string, 3)
	rpc_addrs := make([]string, 3)
	if len(os.Args) > 0 {
		for i := 1; i < 4; i++ {
			peers_addrs[i-1] = os.Args[i]
			rpcaddr := strings.Split(os.Args[i], ":")[1]
			port, _ := strconv.Atoi(rpcaddr)
			rpc_addrs[i-1] = ":" + strconv.Itoa(port+1000)
		}
	}
	if *index < 0 || *index >= len(peers_addrs) {
		fmt.Printf("peer_index %d not invalid\n", *index)
		os.Exit(1)
	}
	peers := groupcache.NewHTTPPool(addrToURL(peers_addrs[*index]))
	var stringcache = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			result, err := ioutil.ReadFile(key)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Printf("asking for %s from dbserver\n", key)
			dest.SetBytes([]byte(result))
			return nil
		}))

	peers.Set(addrsToURLs(peers_addrs)...)

	http.HandleFunc("/zk", func(rw http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Query().Get("key"))
		var data []byte
		k := r.URL.Query().Get("key")
		fmt.Printf("cli asked for %s from groupcache\n", k)
		stringcache.Get(nil, k, groupcache.AllocatingByteSliceSink(&data))
		rw.Write([]byte(data))
	})
	go http.ListenAndServe(rpc_addrs[*index], nil)
	rpcaddr := strings.Split(os.Args[1], ":")[1]
	log.Fatal(http.ListenAndServe(":"+rpcaddr, peers))
}

func addrToURL(addr string) string {
	return "http://" + addr
}

func addrsToURLs(addrs []string) []string {
	result := make([]string, 0)
	for _, addr := range addrs {
		result = append(result, addrToURL(addr))
	}
	return result
}
