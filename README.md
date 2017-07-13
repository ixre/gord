# GoRd
一个使用golang编写的短域名，及302跳转的小应用。

## 应用场景 ##
1. 做301跳转
2. 短域名，比如:to2.net/git 跳转到 http://github.com/jsix

## 配置 ##
配置文件为: *.conf，启动时会加载目录(默认为当前目录，可以通过dir参数指定)下的所有配置。

host项支持通配，如: *.to2.net 能匹配 to2.net的所有子域名; 如host包含多个域，用空格分开。

## 启动 ##

    ./gord
 
第一次运行会生成一个默认的配置示例：

    [
      {
       "host": "*.to2.net",
       "to": "http://www.to2.net/{path}{query}",
       "location": {
       	"/a":      "http://a.com/{path}{query}{timestamp}",
       	"/a/*":    "http://a.com/t-{*}",
        "/1/2/3/": "http://a.com/{#0}-{#1}-{#2}",
       }
      }
     ]
    
    


##  高级应用 ##
### 防止SP缓存 ###
如APP的更新服务器和更新包，直接使用地址可能会被SP强制缓存，使之无论如何无法返回正确的信息。

可以添加{timestamp}来为每个请求自动添加时间戳简单解决。

    [
      {
       "host": "*.to2.net",
       "location": {
          "a.apk":"http://a.com/a.apk{timestamp}"
       }
      }
    ]
   
   
### 整个目录URL进行302跳转 ###

当网页或者资源，目录名称发生变化，可以通过对整个目录进行302跳转，平稳过渡。

    
    [
      {
       "host": "*.to2.net",
       "location": {
          "/a/*":    "http://a.com/b/{*}",
       }
      }
    ]

### 识别URL片段进行302跳转 ###

如博客有地址：/2001/12/30/1.html,现对URL进行优化，将跳转到/2001/12-30/1.html


    [
      {
       "host": "*.to2.net",
       "location": {
		"/*/*/*/*.html": "http://a.com/{#0}/{#1}-{#2}/{#3}.html",
       }
      }
    ]

