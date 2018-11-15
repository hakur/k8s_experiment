FROM registry.3dvisioncloud.com/library/alpine:3.8
# fix library dependencies
# otherwise golang binary may encounter 'not found' error 
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY ./drone-k8-deploy /opt/drone-k8-deploy
COPY kubeconfig /opt
RUN chmod +x /opt/drone-k8-deploy
WORKDIR /opt
CMD ["/opt/drone-k8-deploy","-kubeconfig","/opt/kubeconfig"]