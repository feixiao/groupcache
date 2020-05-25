## groupcache快速入门

### 功能介绍
groupcache是go语言开发的缓存库。用于替换memcache的。
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

### 分析目的
+ consistenthash(提供一致性哈希算法的支持)，
+ lru(提供了LRU方式清楚缓存的算法)，
+ singleflight(保证了多次相同请求只去获取值一次，减少了资源消耗)，
https://segmentfault.com/a/1190000018464029

### 参考资料
+ [《GROUPCACHE EXAMPLE》](https://sconedocs.github.io/groupcacheUseCase/)
+ [Playing with groupcache](https://capotej.com/blog/2013/07/28/playing-with-groupcache/)


https://www.jianshu.com/p/5c3db568b8b8
https://juejin.im/entry/57c3ce697db2a200680ab024