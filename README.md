# GoURI
一个使用golang编写的短域名，及302跳转的小应用。

## 应用场景 ##
1. 做301跳转
2. 短域名，比如:z3q.net/blog 跳转到我的博客http://www.s1n1.com

## 配置文件 ##
配置文件存放于conf目录，后缀为*.conf，启动时会加载改目录下的所有配置。
host项支持通配，如: *.z3q.net 能匹配 z3q.net的所有子域名


第一次运行会生成一个默认的配置示例：

    [
      {
       "host": "*.z3q.net",
       "to": "http://www.z3q.net",
       "map": null
      },
      {
       "host": "s.z3q.net",
       "to": "",
       "map": {
        "blog": "http://www.s1n1.com"
       }
      }
     ]


