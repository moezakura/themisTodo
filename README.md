<!---
 Copyright (c) 2018 mox
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Themis Todo

## Description
- simple todo management
- multi project todo
- multi accounts

## Requirement
- Node.js
	- Webpack
	- Vue.js
	- babel
	- progressbar.js
	- sass
- golang
	- jwt-go
	- multitemplate
	- gin
	- mysql
	- imageupload
- Mysql or MariaDB

## Create database

```
mysql -u root -p < data.sql
```

## For Developers
### Develop Frontend
```
cd www/assets/
npm install #First time only
npm run start
```
Start up webpack debug server "localhost:8652".

### Develop Backend

#### Install library

```
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/gin-contrib/multitemplate
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/olahol/go-imageupload
```
#### modify file
Move config.go.example to config.go and modify your mysql config.  
※ config.go will be replaced in the future by environment.  

#### Statup debug server

```
go run *.go
```
Start up webpack debug server "localhost:31204".

## Install
### Build Frontend
```
cd www/assets/
npm install #First time only
npm run build
```

### Build Backend
#### Install library

```
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/gin-contrib/multitemplate
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/olahol/go-imageupload
```
#### modify file
Move config.go.example to config.go and modify your mysql config.  
※ config.go will be replaced in the future by environment.  

#### Build server

```
go build
```


## Author
Mox [Twitter](http://twitter.com/__MOX__) / [Github](https://github.com/moezakura)  
uryoya [Github](https://github.com/uryoya)  
nishi3 [Github](https://github.com/nishi3)
nitoling [Github](https://github.com/nitoling)

## Version
- 0.1 (Alpha)

## License
- Under the MIT License
- Copyright (c) 2018 mox
