# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

ertgo build -o libos-entry -buildmode=pie -buildvcs=false ../../lib/entry/main.go

docker build -f ./Dockerfile.gramine-ubuntu-20-04 -t wetee/gramine-ubuntu:20.04 .
docker push wetee/gramine-ubuntu:20.04