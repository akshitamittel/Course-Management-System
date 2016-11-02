#!/usr/bin/env bash

#copy all HTML template files to the GOBIN directory for use
rm -rf $GOBIN/templates
cp -r templates $GOBIN

#clear the files directory
rm -rf $GOBIN/files
mkdir $GOBIN/files
