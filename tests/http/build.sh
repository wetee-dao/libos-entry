# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR

tag=`date "+%Y-%m-%d-%H-%M"`

ego-go build -o ./hello ./hello.go

# docker run --device /dev/sgx/enclave --device /dev/sgx/provision \
#     -v ${PWD}:/srv wetee/ego-ubuntu-builder:22.04 \
#     bash -c "cd /srv && ego-go build -o ./hello ./hello.go"

docker build -t registry.cn-hangzhou.aliyuncs.com/wetee_dao/ego-hello:$tag .
docker push registry.cn-hangzhou.aliyuncs.com/wetee_dao/ego-hello:$tag