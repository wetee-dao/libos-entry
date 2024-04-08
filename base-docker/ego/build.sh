# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

docker build -f ./Dockerfile.ego-ubuntu-20-04-deploy -t wetee/ego-ubuntu-deploy:22.04 .
docker push wetee/ego-ubuntu-deploy:22.04