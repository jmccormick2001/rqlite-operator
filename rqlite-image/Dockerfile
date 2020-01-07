FROM registry.access.redhat.com/ubi8/ubi-minimal
MAINTAINER Jeff McCormick
ENV RQLITE_VERSION=5.0.0
RUN microdnf install curl procps hostname psmisc gzip tar && microdnf clean all
RUN curl -L https://github.com/rqlite/rqlite/releases/download/v${RQLITE_VERSION}/rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz -o rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz && \
 tar xvfz rqlite-v${RQLITE_VERSION}-linux-amd64.tar.gz && \
 cp rqlite-v${RQLITE_VERSION}-linux-amd64/* /bin
RUN mkdir -p /rqlite/file
VOLUME /rqlite/file
EXPOSE 4001
ADD rqlite-entrypoint.sh /bin/rqlite-entrypoint.sh
USER root
ENTRYPOINT ["/bin/rqlite-entrypoint.sh"]
