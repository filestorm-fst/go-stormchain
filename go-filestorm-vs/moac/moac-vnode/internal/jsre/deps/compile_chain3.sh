#!/bin/bash
#This is to build the deps for the console interface
#regenerate the binary code from chain3.js using go-bindata 
echo "../go-bindata -nometadata -pkg deps -o bindata.go bignumber.js chain3.js"
#the go-bindata program is used to create the bindata.go with input javascript
#file
../go-bindata -nometadata -pkg deps -o bindata.go bignumber.js chain3.js
#go:generate go-bindata -nometadata -pkg deps -o bindata.go bignumber.js chain3.js
#Simplify the bindata.go
#go:generate gofmt -w -s bindata.go
gofmt -w -s bindata.go