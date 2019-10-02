FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER Jeff McCormick
ENV RQLITE_VERSION=4.5.0
RUN microdnf install curl gzip tar && microdnf clean all
RUN curl -L https://github.com/rqlite/rqlite/releases/download/v${RQLITE_VERSION}/rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz -o rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz && \
 tar xvfz rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz && \
 cp rqlite-v${RQLITE_VERSION}-linux-amd64/rqlited /bin
RUN mkdir -p /rqlite/file
VOLUME /rqlite/file
EXPOSE 4001 4002
ADD docker-entrypoint.sh /bin/docker-entrypoint.sh
ENTRYPOINT ["/bin/docker-entrypoint.sh"]
CMD ["rqlited", "-http-addr", "0.0.0.0:4001", "-raft-addr", "0.0.0.0:4002", "/rqlite/file/data"]
