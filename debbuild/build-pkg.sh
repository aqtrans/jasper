#!/bin/sh
cd ../
go build -o debbuild/jasper
cd debbuild/
debuild -us -uc -b
