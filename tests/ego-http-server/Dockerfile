FROM wetee/ego-ubuntu-deploy:22.04:22.04
WORKDIR /

# Add the hello 
ADD hello  /hello
# Add the hello end

ADD keys   /keys
ADD enclave.json /enclave.json

EXPOSE 8999 65535


ENTRYPOINT ["/bin/sh", "-c" ,"ego sign hello && ego run hello"]