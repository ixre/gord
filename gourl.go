/**
 * Copyright 2015 @ S1N1 Team.
 * name : gourl.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package main
import (
	"net/http"
	"flag"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"encoding/json"
	"bufio"
)

var (
	ConfBasePath = "conf/"
)

func main(){
	var port int
	var dir string
	flag.IntVar(&port,"port",8302,"port")
	flag.StringVar(&dir,"dir","conf/","config file")
	ConfBasePath = dir

	app := http.Server{
		Addr:fmt.Sprintf(":%d",port),
		Handler:getHandler(),
	}

	log.Printf("[ Service] - running on port %s\n",app.Addr)
	if err := app.ListenAndServe();err != nil{
		log.Fatalf(" Server aborted!reason:%s\n",err.Error())
	}
}

type Item struct {
	Host        string                `json:"host"`
	AllLocation string                    `json:"to"`
	Map         map[string]string    `json:"map"`
}

type ItemManager struct{
	items map[string]*Item
	basePath string
}

func (i *ItemManager) checkExists(path string){
	_,err := os.Stat(path)
	if os.IsNotExist(err){
		i.initExample()
	}
}

func (i *ItemManager) initExample(){
	if err := os.MkdirAll(i.basePath,os.ModePerm);err != nil{
		panic(err)
	}
	var defaultItems []*Item =[]*Item{
		&Item{
			Host:"*.z3q.net",
			AllLocation:"http://www.z3q.net",
		},
		&Item{
			Host:"s.z3q.net",
			Map:map[string]string{
				"blog":"http://www.s1n1.com",
			},
		},
	}
	bytes,_ := json.MarshalIndent(defaultItems," "," ")

	f,_ := os.Create(i.basePath+"gourl.conf");
	defer f.Close()
	wr := bufio.NewWriter(f)
	wr.Write(bytes)
	wr.Flush()
}

func (this *ItemManager) Load(){
	this.checkExists(this.basePath)
	filepath.Walk(this.basePath,func(path string, info os.FileInfo, err error)error{
		if err == nil && !info.IsDir() && strings.HasSuffix(path,".conf"){
			items := this.GetItemsFromFile(path)
			this.Append(items)
		}
		return nil
	})
}

func (this *ItemManager) GetItemsFromFile(path string)[]*Item {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	var items []*Item=make([]*Item,0)
	err = json.Unmarshal(bytes,&items)
	if err != nil{
		panic(path+" - "+err.Error())
	}
	return items
}

func (this *ItemManager) Append(items []*Item){
	if this.items == nil{
		this.items = make(map[string]*Item,0)
	}
	if items!= nil {
		for _,v := range items {
			if _, ok := this.items[v.Host]; ok {
				panic("has exists host "+v.Host)
			}
			this.items[v.Host] = v
		}
	}
}

func (this *ItemManager) GetItemByHost(host string)*Item{
	for k,v := range this.items {
		if this.matchHost(k,host){
			return v
		}
	}
	return nil
}

func (this *ItemManager) matchHost(key,host string)bool{
	if host == key{
		return true
	}
	if strings.HasPrefix(key, "*."){
		return strings.HasSuffix(host,key[2:])
	}
	return false
}

var _ http.Handler = new(redirectHandler)
type redirectHandler struct{
	itemManager *ItemManager
}

func (r *redirectHandler) ServeHTTP(rsp http.ResponseWriter,req *http.Request) {
	host := req.Host; // "s.z3q.net" use for test
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

func (r *redirectHandler) getLocation(rsp http.ResponseWriter,req *http.Request,item *Item)(string,bool) {
	path := req.URL.Path
	query := req.URL.RawQuery
	var con string
	if len(query) != 0 {
		con = "?"
	}
	if len(item.AllLocation) != 0 {
		return fmt.Sprintf("%s%s%s%s", item.AllLocation, path, con, query), true
	}
	if v,ok := item.Map[path[1:]]; ok {
		return fmt.Sprintf("%s%s%s", v, con, query), true
	}
	return "", false
}


func getHandler()http.Handler {
	re := &redirectHandler{
		itemManager: &ItemManager{
			basePath:ConfBasePath,
		},
	}
	re.itemManager.Load()
	return re
}
