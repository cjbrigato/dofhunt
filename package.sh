#!/bin/bash

rm -rf build/
zip release clues.json dofhunt.exe
mkdir -p build
mv release.zip build
