# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

EGO_DEB=ego_1.4.1_amd64_ubuntu-20.04.deb  && \
wget https://github.com/edgelesssys/ego/releases/download/v1.4.1/$EGO_DEB  && \

docker build -f ./Dockerfile.ego-ubuntu-20-04-builder -t wetee/ego-ubuntu:20.04 .
docker push wetee/ego-ubuntu:20.04