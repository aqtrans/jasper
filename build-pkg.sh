#!/bin/sh
go build -o jasper
debuild -us -uc -b
