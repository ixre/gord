# GoRd
一个使用golang编写的短域名，及302跳转的小应用。

## 应用场景 ##
1. 做301跳转
2. 短域名，比如:to2.net/git 跳转到 http://github.com/jsix

## 配置 ##
配置文件为: *.conf，启动时会加载目录(默认为当前目录，可以通过dir参数指定)下的所有配置。

host项支持通配，如: *.to2.net 能匹配 to2.net的所有子域名

## 启动 ##

    ./gord
 
第一次运行会生成一个默认的配置示例：

    [
      {
       "host": "*.to2.net",
       "to": "http://www.to2.net/{path}{query}",
       "location": {
        "/1/2/3/": "http://a.com/{#0}-{#1}-{#2}",
        "/a": "http://a.com/{path}{query}",
        "/b/*": "http://b.com/t-{*}"
       }
      }
     ]
    
host多个域，用空格分割。



