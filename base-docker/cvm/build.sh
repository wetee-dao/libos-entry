# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../../

tag=`date "+%Y-%m-%d-%H_%M"`

go build -o ./base-docker/cvm/libos ./lib/cvm/main.go

cd $DIR

docker build -t registry.cn-hangzhou.aliyuncs.com/wetee_dao/cvm:$tag .
docker push registry.cn-hangzhou.aliyuncs.com/wetee_dao/cvm:$tag