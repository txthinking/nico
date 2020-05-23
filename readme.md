# Nico

A HTTP2 web server for reverse proxy and single page application, automatically apply for ssl certificate, zero-configuration.

### Install via [nami](https://github.com/txthinking/nami)

```
$ nami install github.com/txthinking/nico
```

### Reverse proxy

> Make sure your domains are already resolved to your server IP

```
$ nico 'domain.com http://127.0.0.1:2020'
```

### Static server, can be used for single page application

```
$ nico 'domain.com /path/to/web/root'
```

### All can be in one line command

```
$ nico 'domain1.com http://127.0.0.1:2020' 'domain2.com /path/to/web/root' 'domain3.com http://127.0.0.1:3030'
```

### Daemon

You may like [joker](https://github.com/txthinking/joker)

## Why

Nico is simple HTTP2 web server, but she is enough for me. If you want to use rewrite, load balancing, you need to consider nginx or others.

## License

Licensed under The GPLv3 License
