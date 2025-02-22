#!/usr/bin/env bash

set -e

TOPDIR=$PWD
SERVER_DIR=$TOPDIR/build/rest_server
CVLDIR=$TOPDIR/src/cvl

# LD_LIBRARY_PATH for CVL
[ -z $LD_LIBRARY_PATH ] && export LD_LIBRARY_PATH=/usr/local/lib

# Setup CVL schema directory
if [ -z $CVL_SCHEMA_PATH ]; then
    export CVL_SCHEMA_PATH=$CVLDIR/schema
fi

echo "CVL schema directory is $CVL_SCHEMA_PATH"
if [ $(find $CVL_SCHEMA_PATH -name *.yin | wc -l) == 0 ]; then
    echo "WARNING: no yin files at $CVL_SCHEMA_PATH"
    echo ""
fi

EXTRA_ARGS="-ui $SERVER_DIR/dist/ui -logtostderr"
HAS_CRTFILE=
HAS_KEYFILE=

for V in $@; do
    case $V in
    -cert|--cert|-cert=*|--cert=*) HAS_CRTFILE=1;;
    -key|--key|-key=*|--key=*) HAS_KEYFILE=1;;
    esac
done

cd $SERVER_DIR

##
# Setup TLS server cert/key pair
if [ -z $HAS_CRTFILE ] && [ -z $HAS_KEYFILE ]; then
    if [ -f cert.pem ] && [ -f key.pem ]; then
        echo "Reusing existing cert.pem and key.pem ..."
    else 
        echo "Generating temporary server certificate ..."
        ./generate_cert --host localhost
    fi

    EXTRA_ARGS+=" -cert cert.pem -key key.pem"
fi

##
# Start server
./main $EXTRA_ARGS $* 

