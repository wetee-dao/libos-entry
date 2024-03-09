# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../bin

if [ -f "libos-entry" ] ; then
    rm "libos-entry"
fi

ertgo build -o libos-entry -buildmode=pie -buildvcs=false ../lib/entry/main.go