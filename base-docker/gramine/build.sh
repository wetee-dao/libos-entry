# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib ertgo build -o libos-entry -buildmode=pie -buildvcs=false ../../lib/entry/main.go

docker build -f ./Dockerfile.gramine-ubuntu-20-04 -t wetee/ubuntu:20.04 .
docker push wetee/ubuntu:20.04