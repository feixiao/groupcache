## groupcache快速入门

### 功能介绍
groupcache是go语言开发的缓存库。用于替换memcache的。
#### 系统框架
![](./2.png)


#### 代码框架
![](./1.png)



### 使用入门
```shell
cd example
go build
./example -addr=:8080 -pool=http://127.0.0.1:8080 

# 查询
curl localhost:8080/color?name=green
#00FF00
curl localhost:8080/color?name=red
#FF0000

```

### 源码目录
 ```shell script
.
├── byteview.go
├── byteview_test.go
├── consistenthash    // 实现一致性hash功能
├── groupcache.go     // grpc生成的代码，用于远程调用
├── groupcachepb      // grpc生成的代码，用于远程调用
├── http.go           // http相关的代码
├── lru               // 实现缓存的置换算法（最近最少使用）
├── peers.go          // 单个节点的一些接口实现
├── singleflight      // 实现多个同请求的合并，保证“同时”多个同参数的get请求只执行一次操作功能
├── sinks.go
└── testpb

```

### 分析目的
+ consistenthash(提供一致性哈希算法的支持)，
+ lru(提供了LRU方式清楚缓存的算法)，
+ singleflight(保证了多次相同请求只去获取值一次，减少了资源消耗)，
https://segmentfault.com/a/1190000018464029

### 参考资料
+ [《groupcache 设计原理剖析 》](https://www.dazhuanlan.com/2019/12/11/5df07fcb62cae/?__cf_chl_jschl_tk__=e5a47b230d1b9d89eb3887cab036b09f2e3ea621-1590370196-0-AYcPFk14NmbUvag0bCwvLEwPGpXssbJuZhDvEpan7iZiKQi123FXqUvH-LsRSQaov7ybpQtzh-615A-1ZEDC54TuWv_6ZTwsr3zoEwubtJbUbw2J8PTOnzfviGoQB4UWA9Y1ZVzP5QLQ2BCSNlSYxDlegJsosJAV1xJQf06FNkbXPBEAh0SCE29OAzUhpZx1qOKfiUjkI1NNltnexAUoGKVMymm9ocKiWwcq4y_CnUX3xNGz6wyOTmUjQ0RrS1qcQDN8Z-0Jrzn9z1VbzCbEc8R-bdwdkzo7hqaHZ3goA0AQMpxVWxzRjbsy4YIf7vHWEg)
+ [《GROUPCACHE EXAMPLE》](https://sconedocs.github.io/groupcacheUseCase/)
+ [Playing with groupcache](https://capotej.com/blog/2013/07/28/playing-with-groupcache/)


https://www.jianshu.com/p/5c3db568b8b8
https://juejin.im/entry/57c3ce697db2a200680ab024