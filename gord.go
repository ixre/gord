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
	"time"
)

var (
	port  int
	dir   string
	debug bool = false
)

func main() {
	flag.IntVar(&port, "port", 8302, "port")
	flag.StringVar(&dir, "dir", ".", "config file")
	flag.BoolVar(&debug, "log", false, "print log")
	flag.Parse()
	app := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: getHandler(dir),
	}
	log.Printf("[ Gord][ Service] - running on port %s\n", app.Addr)
	if err := app.ListenAndServe(); err != nil {
		log.Fatalf(" Server aborted!reason:%s\n", err.Error())
	}
}

func debugLog(v ...interface{}) {
	if debug {
		v = append([]interface{}{"[ Gord][ Log]"}, v...)
		log.Println(v...)
	}
}

func getHandler(confPath string) http.Handler {
	r := &HttpHandler{
		itemManager: &ItemManager{
			confPath: confPath,
		},
	}
	r.itemManager.Load()
	return r
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
	confPath string
	items    map[string]*Item
}

// 检查目录，并初始化
func (i *ItemManager) checkDir(path string) {
	_, err := os.Stat(path)
	//创建目录
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		i.initExample(path)
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
			i.initExample(path)
		}
	}
}

func (i *ItemManager) initExample(path string) {
	var defaultItems []*Item = []*Item{
		{
			Host: "localhost localhost:8302 *.to2.net",
			To:   "http://www.to2.net/{path}{query}",
			Location: map[string]string{
				"/a":      "http://a.com/{path}{query}{timestamp}",
				"/a/*":    "http://a.com/t-{*}",
				"/1/2/3/": "http://a.com/{#0}-{#1}-{#2}",
			},
		},
	}

	//创建文件
	bytes, _ := json.MarshalIndent(defaultItems, " ", " ")
	f, err := os.Create(path + "/default.conf")
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
	i.checkDir(i.confPath)
	filepath.Walk(i.confPath, func(path string, info os.FileInfo, err error) error {
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
					log.Println(fmt.Sprintf("error config file ：%s; host: %s; error:%s",
						path, v.Host, err.Error()))
					os.Exit(1)
				}
			}
			debugLog("config file ", path, " load success.")
			return items
		}
	}
	log.Println("load config " + path + " error ：" + err.Error())
	os.Exit(1)
	return nil
}

// 增加配置项
func (i *ItemManager) Append(items []*Item) {
	if items == nil {
		return
	}
	if i.items == nil {
		i.items = make(map[string]*Item, 0)
	}

	for _, v := range items {
		hostArr := strings.Split(v.Host, " ")
		for _, host := range hostArr {
			if _, ok := i.items[host]; ok {
				log.Println("host " + host + " already exists ")
				os.Exit(0)
				break
			}
			i.items[host] = v
		}
	}
}

// 根据主机名获取相应的配置,如果无匹配，则默认使用localhost
func (i *ItemManager) GetItemByHost(host string) *Item {
	for k, v := range i.items {
		if i.matchHost(k, host) {
			return v
		}
	}
	return i.items["localhost"]
}

// 匹配主机
func (i *ItemManager) matchHost(cfgHost, host string) bool {
	if host == cfgHost {
		return true
	}
	// 判断是否泛解析
	if strings.HasPrefix(cfgHost, "*.") {
		return strings.HasSuffix(host, cfgHost[2:])
	}
	return false
}

var _ http.Handler = new(HttpHandler)

type HttpHandler struct {
	itemManager *ItemManager
}

func (r *HttpHandler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	host := req.Host
	debugLog("[ Request]: source host ", host)
	var item *Item = r.itemManager.GetItemByHost(host)
	if item != nil {
		if location, b := r.getLocation(req, item); b {
			rsp.Header().Add("Location", location)
			rsp.WriteHeader(302)
			return
		}
	}
	rsp.Write([]byte("Not match any host"))
}

// 获取目标路径
func (r *HttpHandler) getLocation(req *http.Request, item *Item) (string, bool) {
	path := req.URL.Path
	query := req.URL.RawQuery
	concat := ""
	if len(query) != 0 {
		concat = "?"
	}
	//查找匹配
	target := item.To
	anyMatchPos := -1
	for k, v := range item.Location {
		debugLog("[ Compare]:对比相同，path:", path, "; key:", k)
		//判断路径是否相同
		if path == k {
			target = v
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
				target = v
				break
			}
		}
	}
	//无匹配
	target = strings.TrimSpace(target)
	if target == "" {
		return "", false
	}

	//全局请求跳转路径,{path}表示完整的路径；
	if strings.Contains(target, "{path}") {
		target = strings.Replace(target, "{path}", path[1:], -1)
	}
	// {query}表示查询条件
	qt := strings.Contains(target, "{query}")
	if qt {
		target = strings.Replace(target, "{query}", concat+query, -1)
	}
	// {timestamp}会返回时间戳
	if strings.Contains(target, "{timestamp}") {
		unixTime := strconv.Itoa(int(time.Now().UnixNano()))
		if !qt || concat == "" {
			unixTime = "?_stamp=" + unixTime
		} else {
			unixTime = "&_stamp=" + unixTime
		}
		target = strings.Replace(target, "{timestamp}", unixTime, -1)
	}
	//路径通配
	if strings.Contains(target, "{*}") && anyMatchPos != -1 {
		target = strings.Replace(target, "{*}", path[anyMatchPos:], -1)
	}
	//匹配含有路径片段的URL,{#序号}表示指定的路径片段
	if strings.Contains(target, "{#") {
		segments := strings.Split(path[1:], "/")
		for i, l := 0, len(segments); i < l; i++ {
			target = strings.Replace(target, "{#"+strconv.Itoa(i)+"}",
				segments[i], -1)
		}
	}
	debugLog("--- origin:", path, "; target:", target)
	return target, true
}
