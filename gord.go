/**
 * Copyright 2015 @ at3.net.
 * name : gord.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	port  int
	dir   string
	debug bool = false
)

func main() {
	flag.IntVar(&port, "port", 8302, "port")
	flag.StringVar(&dir, "dir", "./", "config file")
	flag.BoolVar(&debug, "log", false, "print log")
	flag.Parse()
	baseDir := dir
	app := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: getHandler(baseDir),
	}
	log.Printf("[ GoRd][ Service] - running on port %s\n", app.Addr)
	if err := app.ListenAndServe(); err != nil {
		log.Fatalf(" Server aborted!reason:%s\n", err.Error())
	}
}

func debugLog(v ...interface{}) {
	if debug {
		v = append([]interface{}{"[ GoRd][ Log]"}, v...)
		log.Println(v...)
	}
}

func getHandler(baseDir string) http.Handler {
	re := &redirectHandler{
		itemManager: &ItemManager{
			basePath: baseDir,
		},
	}
	re.itemManager.Load()
	return re
}

type Item struct {
	//主机头，*表示通配
	Host string `json:"host"`
	//全局请求跳转路径,{path}表示完整的路径；
	//{#序号}表示路径片段的序号
	To string `json:"to"`
	//如果未设定全局请求跳转路径，那么将启用路径字典
	//如果{"a/b/c":"http://abc.com"}，访问/a/b/c将跳转
	//到"http://abc.com"
	Location map[string]string `json:"location"`
}

type ItemManager struct {
	items    map[string]*Item
	basePath string
}

// 检查目录，并初始化
func (i *ItemManager) checkDir(path string) {
	_, err := os.Stat(path)
	//创建目录
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		i.initExample()
	} else {
		//是否存在.conf文件,不存在，则初始化
		fi, _ := os.Open(path)
		exits := false
		list, _ := fi.Readdirnames(-1)
		for _, v := range list {
			if strings.HasSuffix(v, ".conf") {
				exits = true
			}
		}
		if !exits {
			i.initExample()
		}
	}
}

func (i *ItemManager) initExample() {
	var defaultItems []*Item = []*Item{
		{
			Host: "*.to2.net",
			To:   "http://www.to2.net/{path}{query}",
			Location: map[string]string{
				"/a":      "http://a.com",
				"/b/*":    "http://b.com/t-{*}",
				"/a/b":    "http://a.com/{path}{query}",
				"/1/2/3/": "http://a.com/{#0}-{#1}-{#2}",
			},
		},
	}

	//创建文件
	bytes, _ := json.MarshalIndent(defaultItems, " ", " ")
	f, err := os.Create("example.conf")
	if err == nil {
		wr := bufio.NewWriter(f)
		wr.Write(bytes)
		err = wr.Flush()
		f.Close()
	}
	if err != nil {
		log.Println(" init example config fail ! error :", err.Error())
	}
}

// 加载配置
func (i *ItemManager) Load() {
	i.checkDir(i.basePath)
	filepath.Walk(i.basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".conf") {
			items := i.GetItemsFromFile(path)
			i.Append(items)
		}
		return nil
	})
}

func (i *ItemManager) checkItem(item *Item) error {
	if len(item.Host) < 2 {
		return errors.New("主机头长度不正确")
	}
	return nil
}

// 从文件中加载配置项目
func (i *ItemManager) GetItemsFromFile(path string) []*Item {
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		//从文件中反序列化
		items := make([]*Item, 0)
		err = json.Unmarshal(bytes, &items)
		//检查配置是否正确
		if err == nil {
			for _, v := range items {
				if err = i.checkItem(v); err != nil {
					log.Println(fmt.Sprintf("配置有误! 文件：%s; host: %s; error:%s",
						path, v.Host, err.Error()))
					os.Exit(1)
				}
			}
			debugLog("[ Load]: conf file:", path, "load ok!")
			return items
		}
	}
	log.Println("加载配置:" + path + "出错：" + err.Error())
	os.Exit(1)
	return nil
}

// 增加配置项
func (i *ItemManager) Append(items []*Item) {
	if i.items == nil {
		i.items = make(map[string]*Item, 0)
	}
	if items != nil {
		for _, v := range items {
			if _, ok := i.items[v.Host]; ok {
				panic("has exists host " + v.Host)
			}
			i.items[v.Host] = v
		}
	}
}

// 根据主机名获取相应的配置
func (i *ItemManager) GetItemByHost(host string) *Item {
	for k, v := range i.items {
		if i.matchHost(k, host) {
			return v
		}
	}
	return nil
}

// 匹配主机
func (i *ItemManager) matchHost(hostKey, host string) bool {
	if host == hostKey {
		return true
	}
	// 判断是否泛解析
	if strings.HasPrefix(hostKey, "*.") {
		return strings.HasSuffix(host, hostKey[2:])
	}
	return false
}

var _ http.Handler = new(redirectHandler)

type redirectHandler struct {
	itemManager *ItemManager
}

func (r *redirectHandler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	host := req.Host
	// host = "www.to2.net" // "z3q.net" use for test
	var item *Item = r.itemManager.GetItemByHost(host)
	if item != nil {
		if location, b := r.getLocation(rsp, req, item); b {
			rsp.Header().Add("Location", location)
			rsp.WriteHeader(302)
			return
		}
	}
	rsp.Write([]byte("Not match any host"))
}

func (r *redirectHandler) getLocation(rsp http.ResponseWriter,
	req *http.Request, item *Item) (string, bool) {
	path := req.URL.Path
	query := req.URL.RawQuery
	concat := ""
	if len(query) != 0 {
		concat = "?"
	}
	//查找匹配
	to := item.To
	anyMatchPos := -1
	for k, v := range item.Location {
		debugLog("[ Compare]:对比相同，path:", path, "; key:", k)
		//判断路径是否相同
		if path == k {
			to = v
			break
		}
		//匹配如：/d/* 含通配符的路径
		anyMatch := strings.HasSuffix(k, "*")
		debugLog("[ Compare]:包含通配:", anyMatch)
		if anyMatch {
			anyMatchPos = len(k) - 1 //通配符所在的索引位置
			anyMatch = strings.HasPrefix(path, k[:anyMatchPos])
			debugLog("[ Compare]:判断通配:", anyMatch, k[:anyMatchPos])
			if anyMatch {
				to = v
				break
			}
		}
	}
	//无匹配
	to = strings.TrimSpace(to)
	if to == "" {
		return "", false
	}

	//全局请求跳转路径,{path}表示完整的路径；
	if strings.Contains(to, "{path}") {
		to = strings.Replace(to, "{path}", path[1:], -1)
	}
	if strings.Contains(to, "{query}") {
		to = strings.Replace(to, "{query}", concat+query, -1)
	}
	//路径通配
	if strings.Contains(to, "{*}") && anyMatchPos != -1 {
		to = strings.Replace(to, "{*}", path[anyMatchPos:], -1)
	}
	//匹配含有路径片段的URL,{#序号}表示指定的路径片段
	if strings.Contains(to, "{#") {
		segments := strings.Split(path[1:], "/")
		for i, l := 0, len(segments); i < l; i++ {
			to = strings.Replace(to, "{#"+strconv.Itoa(i)+"}",
				segments[i], -1)
		}
	}
	debugLog("--- origin:", path, "; target:", to)
	return to, true
}
