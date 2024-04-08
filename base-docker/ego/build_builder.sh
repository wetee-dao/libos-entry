# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

# EGO_DEB=ego_1.5.0_amd64_ubuntu-22.04.deb  && \
# wget https://github.com/edgelesssys/ego/releases/download/v1.5.0/$EGO_DEB  && \

docker build -f ./Dockerfile.ego-ubuntu-20-04-builder -t wetee/ego-ubuntu:22.04 .
docker push wetee/ego-ubuntu:22.04