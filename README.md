# HorribleNotifier
![GitHub](https://img.shields.io/github/license/pizza61/horriblenotifier?style=flat-square)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/pizza61/horriblenotifier?include_prereleases&style=flat-square)
![GitHub All Releases](https://img.shields.io/github/downloads/pizza61/horriblenotifier/total?style=flat-square)

Simple desktop app for notifying about new [HorribleSubs](https://horriblesubs.info) releases

Currently only Windows 10 supported. Linux support in progress

## Installation

[Downloads](https://github.com/pizza61/horriblenotifier/releases)

### Building
Install deps and

```
rice embed-go
rsrc -manifest horriblenotifier.manifest -ico icons\hn-128.ico -o rsrc.syso
go build -ldflags "-H windowsgui -s -w" -o HorribleNotifier.exe
```

## Screenshots
![Configuration panel](https://i.imgur.com/e9b0u3e.png)

![Notification](https://i.imgur.com/DTnsbpj.png)

## TODO
* Linux support
* Automatic updates