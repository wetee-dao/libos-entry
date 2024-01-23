# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../bin

rm libos-entry
CGO_CFLAGS=-I/opt/ego/include CGO_LDFLAGS=-L/opt/ego/lib ertgo build -o libos-entry -buildmode=pie -buildvcs=false ../lib/entry/main.go