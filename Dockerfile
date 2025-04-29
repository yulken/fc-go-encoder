FROM alpine:3.11 AS builder
ENV PATH="$PATH:/bin/bash" \
    PATH="$PATH:/opt/bento4/bin"

RUN apk add --update ffmpeg bash make

WORKDIR /tmp/bento4
ENV BENTO4_BASE_URL="http://zebulon.bok.net/Bento4/source/" \
    BENTO4_VERSION="1-6-0-641" \
    BENTO4_CHECKSUM="ed3e2603489f4748caadccb794cf37e5e779422e" \
    BENTO4_TARGET="" \
    BENTO4_PATH="/opt/bento4" \
    BENTO4_TYPE="SRC"
    # download and unzip bento4
RUN apk add --update --upgrade python unzip bash gcc g++ scons && \
    wget -q ${BENTO4_BASE_URL}/Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    sha1sum -b Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip | grep -o "^$BENTO4_CHECKSUM " && \
    mkdir -p ${BENTO4_PATH} && \
    unzip Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip -d ${BENTO4_PATH} && \
    rm -rf Bento4-${BENTO4_TYPE}-${BENTO4_VERSION}${BENTO4_TARGET}.zip && \
    apk del unzip && \
    cd ${BENTO4_PATH} && scons -u build_config=Release target=x86_64-unknown-linux

FROM golang:1.24.2-alpine3.21
ENV PATH="$PATH:/bin/bash" \
    PATH="$PATH:/opt/bento4/bin" \
    BENTO4_PATH="/opt/bento4" 

WORKDIR ${BENTO4_PATH}
COPY --from=builder ${BENTO4_PATH}/Build/Targets/x86_64-unknown-linux/Release ${BENTO4_PATH}/bin
COPY --from=builder ${BENTO4_PATH}/Source/Python/utils ${BENTO4_PATH}/utils
COPY --from=builder ${BENTO4_PATH}/Source/Python/wrappers/. ${BENTO4_PATH}/bin
    
RUN apk add --update --upgrade python3 ffmpeg bash make gcc build-base

WORKDIR /go/src

ENTRYPOINT ["tail", "-f", "/dev/null"]