# RuleX

RuleX 是一个轻量级网关，支持多种数据接入以及数据流筛选，可以理解为一个数据路由器。

> 当前处于极其不稳定阶段,请勿尝试.

## 快速开始
### 构建
```sh
git clone https://github.com/wwhai/rulex.git
cd rulex
make # on windows: make windows
```
### 启动
```sh
./rulex run
2021/09/20 17:09:05 cfg.go:24: [info] Init rulex config 
2021/09/20 17:09:05 cfg.go:34: [info] Rulex config init success. 
2021/09/20 17:09:05 utils.go:71: [info] 
 -----------------------------------------------------------     
~~~/=====\       ██████╗ ██╗   ██╗██╗     ███████╗██╗  ██╗       
~~~||\\\||--->o  ██╔══██╗██║   ██║██║     ██╔════╝╚██╗██╔╝       
~~~||///||--->o  ██████╔╝██║   ██║██║     █████╗   ╚███╔╝        
~~~||///||--->o  ██╔══██╗██║   ██║██║     ██╔══╝   ██╔██╗        
~~~||\\\||--->o  ██║  ██║╚██████╔╝███████╗███████╗██╔╝ ██╗       
~~~\=====/       ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝
-----------------------------------------------------------
2021/09/20 17:09:05 utils.go:74: [info] rulex start successfully
2021/09/20 17:09:05 http_api_server.go:139: [info] Http server started on http://127.0.0.1:2580
2021/09/20 17:09:05 grpc_resource.go:92: [info] RulexRpc resource started on [::]:2581
2021/09/20 17:09:05 coap_resource.go:71: [info] Coap resource started on [udp]:2582
2021/09/20 17:09:05 http_resource.go:47: [info] HTTP resource started on [0.0.0.0]:2583
2021/09/20 17:09:05 udp_resource.go:50: [info] UDP resource started on [0.0.0.0]:2584
```

## Dashboard
```
http://127.0.0.1:2580
```

## HTTP API
```
plugin\http_server\openapi.yml
```
## 详细文档
<a href="https://wwhai.github.io/rulex_doc_html/inde core.html">查看文档</a>
