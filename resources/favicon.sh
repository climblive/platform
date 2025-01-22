#!/bin/bash
set -e

for app in admin scorecard scoreboard; do
    BASE="../web/$app/public"
    cp favicon.svg $BASE
    inkscape -w 192 -h 192 favicon.svg -o $BASE/android-chrome-192x192.png
    inkscape -w 512 -h 512 favicon.svg -o $BASE/android-chrome-512x512.png
    inkscape -w 180 -h 180 favicon.svg -o $BASE/apple-touch-icon.png
    inkscape -w 16 -h 16 favicon.svg -o $BASE/favicon-16x16.png
    inkscape -w 32 -h 32 favicon.svg -o $BASE/favicon-32x32.png
    inkscape -w 48 -h 48 favicon.svg -o favicon-48x48.png
    inkscape -w 96 -h 96 favicon.svg -o $BASE/favicon-96x96.png
    convert $BASE/favicon-16x16.png $BASE/favicon-32x32.png favicon-48x48.png $BASE/favicon.ico
    rm favicon-48x48.png
done 
