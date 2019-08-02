## BUILDER #####################################################################

FROM centos:7 as builder

RUN mkdir -p "/go/src" && chmod -R 777 "/go"

ENV GOPATH=/go

WORKDIR /go/src/github.com/essentialkaos/bibop

COPY . .

RUN yum -y -q install https://yum.kaos.st/kaos-repo-latest.el7.noarch.rpm && \
    yum -y -q install make golang git upx && \
    make deps && \
    make all && \
    upx bibop

## FINAL IMAGE #################################################################

FROM centos:7

LABEL name="Bibop Image on CentOS 7" \
      vendor="ESSENTIAL KAOS" \
      maintainer="Anton Novojilov" \
      license="EKOL" \
      version="2019.08.02"

COPY --from=builder /go/src/github.com/essentialkaos/bibop/bibop /usr/bin/
COPY --from=builder /go/src/github.com/essentialkaos/bibop/bibop-entrypoint /usr/bin/

RUN yum -y -q install https://yum.kaos.st/kaos-repo-latest.el7.noarch.rpm && \
    yum -y -q install epel-release

VOLUME /bibop
WORKDIR /bibop

ENTRYPOINT ["bibop-entrypoint"]

################################################################################
