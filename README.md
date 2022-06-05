#  cdnlookup 
一个使用 Edns-Client-Subnet(ECS) 遍历智能DNS节点IP地址的工具

## 原理

#### 智能DNS
为了让用户有更好的网络体验，让用户连接到地理位置更近的服务器， NS服务器会根据客户端IP地址来判断用户所在区域及运营商，来返回距离较近的节点。

早期，NS服务器通常无法直接获取到客户端ip，只能获取到上级公共DNS服务器地址。

####  Edns Client Subnet(ECS)

ECS 是由Google提交的一份DNS扩展协议，主要作用是传递用户的IP地址给权威DNS服务器。
需要获取到用户ip地址 
[rfc7871](https://datatracker.ietf.org/doc/html/rfc7871) （2016 年 5 月）

遵循ECS标准的公共DNS，会将经遮罩脱敏后的客户端ip添加至DNS扩展区域( [EDNS rfc6891](https://datatracker.ietf.org/doc/html/rfc6891))传递至NS服务器 （通常为/24）

这样NS服务器就可以获取到模糊的客户端ip，这足以用于判断用户运营商和位置信息。

####  cdnlookup
这个工具会直接发送包含自定义IP的ECS数据的DNS请求，诱导NS服务器返回对应IP的解析结果。

经测试，国内大部分公共DNS都不支持自定义ECS。  谷歌DNS 8.8.8.8 可以正常使用

除了公共DNS，也可以直接将带有ECS数据的DNS请求发送到目标NS服务器，获取解析结果。

# 使用
````
-d  域名  (默认 www.taobao.com)

-i  只输出IP地址列表

-ip 客户端ip

-r  请求重复轮数

-s DNS服务器地址 (默认 8.8.8.8:53)
````

自定义客户端ip
````
cdnlookup.exe -d www.taobao.com  -ip 1.2.3.4
219.147.75.XXX
219.147.75.XXX
````

使用内置的实例ip列表 (内置列表可能会出现判断错误. 建议使用家宽ip段地址定位. )
````
cdnlookup.exe -d www.taobao.com

北京市 教育网
36.99.228.XXX
36.99.228.XXX
吉林 长春 移动
111.26.147.XXX
111.26.147.XXX
辽宁 沈阳 电信
59.47.225.XXX
59.47.225.XXX
......
````