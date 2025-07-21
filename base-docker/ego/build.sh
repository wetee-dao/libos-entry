# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

# EGO_DEB=ego_1.7.2_amd64_ubuntu-24.04.deb
# wget https://github.com/edgelesssys/ego/releases/download/v1.7.2/$EGO_DEB

docker build -f ./Dockerfile.ego-ubuntu-22-04-deploy -t wetee/ego-ubuntu-24-04:1.7.2 .
docker push wetee/ego-ubuntu-24-04:1.7.2