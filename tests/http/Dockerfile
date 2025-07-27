FROM wetee/ego-ubuntu-24-04:1.7.2
WORKDIR /

# Add the hello 
ADD hello  /hello
# Add the hello end

ADD keys   /keys
ADD enclave.json /enclave.json

EXPOSE 8999 65535


ENTRYPOINT ["/bin/sh", "-c" ,"ego run hello"]