# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

# EGO_DEB=ego_1.5.3_amd64_ubuntu-$(lsb_release -rs).deb
# wget https://github.com/edgelesssys/ego/releases/download/v1.5.3/$EGO_DEB

docker build -f ./Dockerfile.ego-ubuntu-22-04-deploy -t wetee/ego-ubuntu-deploy:24.04 .
docker push wetee/ego-ubuntu-deploy:24.04