## groupcache快速入门

### 功能介绍
groupcache是go语言开发的缓存库。用于替换memcache的。
#### 代码框架
![](./1.png)

### 使用入门
```shell
cd example
go build
./main 127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003
```

### 分析目的
+ consistenthash(提供一致性哈希算法的支持)，
+ lru(提供了LRU方式清楚缓存的算法)，
+ singleflight(保证了多次相同请求只去获取值一次，减少了资源消耗)，
https://segmentfault.com/a/1190000018464029

### 参考资料
+ [《groupcache 使用入门》](http://betazk.github.io/2014/12/groupcache%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8%E7%9A%84%E4%B8%80%E4%B8%AA%E7%AE%80%E5%8D%95%E4%BE%8B%E5%AD%90/)
+ [Playing with groupcache](https://capotej.com/blog/2013/07/28/playing-with-groupcache/)


https://www.jianshu.com/p/5c3db568b8b8
https://juejin.im/entry/57c3ce697db2a200680ab024