#!/bin/bash
# Copyright (C) 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# MacOS Package Script.
set -ex

if [ $# -ne 1 -o ! -f "$1/pkg/build.properties" ]; then
	echo Expected the build folder as an argument.
	exit 1
fi

function absname {
  echo $(cd "$1" && pwd)
}

BUILD_OUT=$1
SRC=$(absname "$(dirname "${BASH_SOURCE[0]}")")

rm -rf "$BUILD_OUT/dist"
mkdir -p "$BUILD_OUT/dist"
pushd "$BUILD_OUT/dist"
VERSION=$(awk -F= 'BEGIN {major=0; minor=0; micro=0}
                  /Major/ {major=$2}
                  /Minor/ {minor=$2}
                  /Micro/ {micro=$2}
                  END {print major"."minor"."micro}' ../pkg/build.properties)

# Combine package contents.
mkdir -p gapid/jre
cp -r ../pkg/* gapid/
cp -r ../current/java/gapic-osx.jar gapid/lib/gapic.jar
"$SRC/copy_jre.sh" gapid/jre

# Create a zip file.
zip -r gapid-$VERSION-macos.zip gapid/

# Create a .app package
mkdir -p GAPID.app/Contents/MacOS/
cp -r gapid/* GAPID.app/Contents/MacOS/
cp "$SRC/Info.plist" GAPID.app/Contents/

# Create the icon. TODO: need resolution up to 1024 (512@2x)
mkdir -p GAPID.iconset GAPID.app/Contents/Resources
# Ensure the icon has an alpha channel to make iconutil work, sigh.
pip install --user pypng
python -c '
import sys;import png;i=png.Reader(sys.stdin).asRGBA();
png.Writer(width=i[0],height=i[1],alpha=True).write(sys.stdout,i[2])'\
  < "$SRC/../../gapic/res/icons/logo@2x.png" > logo.png
for i in 128 64 32 16; do
  sips -z $i $i logo.png --out GAPID.iconset/icon_${i}x$i.png
  sips -z $((i*2)) $((i*2)) logo.png --out GAPID.iconset/icon_${i}x$i\@2x.png
done
iconutil -c icns -o GAPID.app/Contents/Resources/GAPID.icns GAPID.iconset

# Make a dmg file.
pip install --user dmgbuild pyobjc-framework-Quartz
cp "$SRC/background@2x.png" .
# Yes, height, then width.... sigh.
sips -z 480 640 background\@2x.png --out background.png
cp "$SRC/dmg-settings.py" .
~/Library/Python/2.7/bin/dmgbuild -s dmg-settings.py GAPID gapid-$VERSION-macos.dmg
