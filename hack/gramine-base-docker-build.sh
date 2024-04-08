# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../base-docker/gramine

cp $DIR/../bin/libos-entry ./libos-entry

docker build -f ./Dockerfile.gramine-ubuntu-20-04 -t wetee/gramine-ubuntu:22.04 .
docker push wetee/gramine-ubuntu:22.04