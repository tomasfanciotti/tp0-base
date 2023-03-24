FROM subfuzion/netcat

WORKDIR /app/
COPY ./netcat/test-server.sh .
RUN chmod 111 test-server.sh
ENTRYPOINT ["/bin/sh"]