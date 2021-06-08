# Nico

[English](readme.md)

[![捐赠](https://img.shields.io/badge/%E6%94%AF%E6%8C%81-%E6%8D%90%E8%B5%A0-ff69b4.svg)](https://www.txthinking.com/opensource-support.html)
[![交流群](https://img.shields.io/badge/%E7%94%B3%E8%AF%B7%E5%8A%A0%E5%85%A5-%E4%BA%A4%E6%B5%81%E7%BE%A4-ff69b4.svg)](https://docs.google.com/forms/d/e/1FAIpQLSdzMwPtDue3QoezXSKfhW88BXp57wkbDXnLaqokJqLeSWP9vQ/viewform)

一个HTTP2 web server, 支持反向代理和单页应用, 自动TLS证书. 零配置.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### 用 [nami](https://github.com/txthinking/nami) 安装

```
$ nami install github.com/txthinking/nico
```

### 静态服务器, 支持单页应用

> 确保你的域名已经指向你的服务器, 并且防火墙已经开放服务器的80/443端口

```
$ nico "domain.com /path/to/web/root"
```

### 反向代理

```
$ nico "domain.com http://127.0.0.1:2020"
```

### 反向代理 https 网站

```
$ nico "domain.com https://reactjs.org"
```

### 根据路径分发

> Exact match: domain.com/ws<br/>
> Prefix match when / is suffix: domain.com/api/<br/>
> Default match: domain.com<br/>
> A special one: domain.com/ is exact match

```
$ nico "domain.com /path/to/web/root" "domain.com/ws http://127.0.0.1:9999" "domain.com/api/ http://127.0.0.1:2020"
```

### 多个域名

```
$ nico "domain0.com /path/to/web/root" "domain1.com /another/web/root" "domain1.com/ws http://127.0.0.1:9999" "domain1.com/api/ http://127.0.0.1:2020"
```

### 守护进程

你可能喜欢 [joker](https://github.com/txthinking/joker)

## 为什么

Nico 是一个简单的HTTP2 web server, 但是在很多时候她已经足够了. 如果你需要更多复杂的功能, 可以考虑nginx等

## 开源协议

基于 GPLv3 协议开源
