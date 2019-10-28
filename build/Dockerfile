FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/rq \
    USER_UID=1001 \
    USER_NAME=rq

# install operator binary
COPY build/_output/bin/rqlite-operator ${OPERATOR}

COPY templates /usr/local/bin
COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
