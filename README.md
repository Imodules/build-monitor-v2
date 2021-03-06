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
6. And a MongoDb

## Setup
1. ```cd ./server; ensure dep; cd ..```
2. ```cd ./client; yarn install; elm-package install; cd ..```

## Test Server
```cd ./server; goconvey```

## Run Dev
```cd ./client; yarn start```

## Build
```make```

## Config
| Flag                  | Env                    | Default                                    |
|-----------------------|------------------------|--------------------------------------------|
| -db                   | BM_DB                  | mongodb://localhost:27017/build-monitor-v2 |
| -port                 | BM_PORT                | 3030                                       |
| -client-path          | BM_CLIENT_PATH         | ../client/dist                             |
| -allowed-origin       | BM_ALLOWED_ORIGIN      | *                                          |
| -password-salt        | BM_PASSWORD_SALT       | you-really-need-to-change-this             |
| -jwt-secret           | BM_JWT_SECRET          | you-really-need-to-change-this-one-also    |
| -tc-url               | BM_TC_URL              | http://localhost:3031                      |

## Other dependencies
```docker run --name dev-mongo -p 27017:27017 -d mongo```