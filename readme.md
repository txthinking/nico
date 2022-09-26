# Nico

[中文](readme_zh.md)

[🗣 News](https://t.me/txthinking_news)
[💬 Chat](https://t.me/brookgroup)
[🩸 Youtube](https://www.youtube.com/txthinking) 
[❤️ Sponsor](https://github.com/sponsors/txthinking)

A HTTP2 web server for reverse proxy and single page application, automatically apply for ssl certificate, zero-configuration.

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install via [nami](https://github.com/txthinking/nami)

```
$ nami install nico
```

### Static server, can be used for single page application

> Make sure your domains are already resolved to your server IP and open 80/443 port

```
$ nico domain.com /path/to/web/root
```

### Reverse proxy

```
$ nico domain.com http://127.0.0.1:2020
```

### Reverse proxy https website

```
$ nico domain.com https://reactjs.org
```

### Dispatch according to path

> Exact match: domain.com/ws<br/>
> Prefix match when / is suffix: domain.com/api/<br/>
> Default match: domain.com<br/>
> A special one: domain.com/ is exact match

```
$ nico domain.com /path/to/web/root domain.com/ws http://127.0.0.1:9999 domain.com/api/ http://127.0.0.1:2020
```

### Multiple domains

```
$ nico domain0.com /path/to/web/root domain1.com /another/web/root domain1.com/ws http://127.0.0.1:9999 domain1.com/api/ http://127.0.0.1:2020
```

### Daemon

You may like [joker](https://github.com/txthinking/joker)

## Why

Nico is a simple HTTP2 web server, but she is enough in most cases. If you want to use rewrite, load balancing, you need to consider nginx or others.

## License

Licensed under The GPLv3 License
