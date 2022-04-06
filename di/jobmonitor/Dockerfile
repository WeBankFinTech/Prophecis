FROM ffdl/ffdl-lcm:v0.1.2
# switch to root
USER root

ADD bin/main /
ADD certs/* /etc/ssl/dlaas/
#COPY charts/ /charts/
RUN chmod 755 /main

WORKDIR /

CMD ["/main"]
