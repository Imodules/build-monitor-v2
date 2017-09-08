# build-monitor-v2

A Teamcity build radiator used to project passing / failing builds in
our dev shop.

*re-written from (https://github.com/Imodules/build-monitor) now with Go and Elm*

## Dependencies
1. Go (https://golang.org/)
2. Elm (http://elm-lang.org/)
3. Yarn (https://yarnpkg.com)
4. Go Dep (https://github.com/golang/dep)
5. Goconvey (http://goconvey.co/)

## Setup
1. ```cd ./server; ensure dep; cd ..```
2. ```cd ./client; yarn install; elm-package install; cd ..```

## Test Server
```cd ./server; goconvey```

## Run Dev
```cd ./client; yarn start```

## Build
```make```