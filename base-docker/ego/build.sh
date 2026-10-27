# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

docker build -f ./Dockerfile.ego-ubuntu-24-04-deploy -t wetee/ego-ubuntu-24-04:1.10.1 .
docker push wetee/ego-ubuntu-24-04:1.10.1
