#!/bin/bash

## Build a very minimal .deb without installing devscripts

set -euo pipefail

VERSION=$1
APPNAME=jasper

mkdir -p dpkg/usr/bin
mkdir -p dpkg/lib/systemd/system
cp $APPNAME dpkg/usr/bin/
cp $APPNAME.service dpkg/lib/systemd/system/

mkdir -p dpkg/DEBIAN

sed -i "s/Version: *.*/Version: 1.0."$(date +%s)"/g" debian/control

cp debian/control debian/preinst debian/postinst debian/prerm debian/postrm dpkg/DEBIAN/

chmod +x dpkg/DEBIAN/preinst dpkg/DEBIAN/postinst dpkg/DEBIAN/prerm dpkg/DEBIAN/postrm
chmod +x dpkg/usr/bin/$APPNAME

dpkg-deb --build --root-owner-group dpkg 
mv dpkg.deb $APPNAME-$VERSION.deb

echo "Package built at $APPNAME-$VERSION.deb"

rm -rf dpkg/
