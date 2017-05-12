# GoRd
一个使用golang编写的短域名，及302跳转的小应用。

## 应用场景 ##
1. 做301跳转
2. 短域名，比如:z3q.net/blog 跳转到我的博客http://www.s1n1.com

## 配置 ##
配置文件后缀为*.conf，启动时会加载改目录下的所有配置。 默认加载根目录下的配置。
host项支持通配，如: *.z3q.net 能匹配 z3q.net的所有子域名

## 启动 ##
./gord -dir=./ -port=8032

第一次运行会生成一个默认的配置示例：

    [
      {
       "host": "*.to2.net",
       "to": "http://www.to2.net/{path}{query}",
       "location": {
        "/1/2/3/": "http://a.com/{#0}-{#1}-{#2}",
        "/a": "http://a.com",
        "/a/b": "http://a.com/{path}{query}",
        "/b/*": "http://b.com/{*}"
       }
      }
     ]


