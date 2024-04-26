# SPDX-FileCopyrightText: Copyright (c) 2023 Intel Corporation
#
# SPDX-License-Identifier: BSD-3-Clause

# Use the specific version of ubuntu image
FROM ubuntu:22.04@sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f as builder

# Set the maintainer label
LABEL maintainer="qiang.han@intel.com"

# Set environment variables
ARG http_proxy="http://proxy-mu.intel.com:911"
ARG https_proxy="http://proxy-mu.intel.com:912"
ARG ftp_proxy="http://proxy-mu.intel.com:21"
ARG socks_proxy="http://proxy-mu.intel.com:1080"
ARG no_proxy="intel.com,.intel.com,.cluster.local,.local,localhost,10.0.0.0/8,192.168.0.0/16,127.0.0.0/8,172.16.0.0/12,134.134.0.0/16"

ENV DEBIAN_FRONTEND noninteractive
ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV ftp_proxy=$ftp_proxy
ENV socks_proxy=$socks_proxy
ENV no_proxy=$no_proxy

# Update and install package dependencies in a single RUN command to reduce layer size
RUN apt-get update && apt-get full-upgrade -y && \
    apt-get install --no-install-recommends -y apt-utils sudo build-essential cmake git ca-certificates curl libbsd-dev libssl-dev nlohmann-json3-dev clang-format nodejs npm && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set the MCM directory
ARG MCM_DIR="/opt/mcm"
ARG JWT_CPP_DIR="/opt/jwt-cpp"
ARG PREFIX_DIR="/usr/local"

# jwt-cpp requirement for catena-sdk
RUN git clone https://github.com/Thalhammer/jwt-cpp ${JWT_CPP_DIR}
WORKDIR ${JWT_CPP_DIR}/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} ..  && \
    make && \
    sudo make install

# Create a non-root user
ARG USER=docker
ARG UID=1000
RUN useradd -l -m -s /bin/bash -u ${UID} ${USER}

# Set default environment variables for the user
ARG HOME="/home/${USER}"
ENV PATH="${PREFIX_DIR}/bin:$HOME/.local/bin:$HOME/bin:$HOME/usr/bin:$PATH"
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib/x86_64-linux-gnu:${LD_LIBRARY_PATH}"

# gRPC
ARG GPRC_VERSION="v1.58.0"
RUN git clone --recurse-submodules -b ${GPRC_VERSION} --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc \
    && mkdir -p /tmp/grpc/cmake/build
WORKDIR /tmp/grpc/cmake/build
RUN cmake -DgRPC_INSTALL=ON \
        -DgRPC_BUILD_TESTS=OFF \
        -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} \
        ../.. && \
    make -j `nproc` && \
    make install && \
    rm -rf * && \
    cmake -DgRPC_INSTALL=ON \
        -DgRPC_BUILD_TESTS=ON \
        -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} \
        ../.. && \
    make -j `nproc` grpc_cli && \
    cp grpc_cli ${PREFIX_DIR}/bin/ && \
    rm -rf /tmp/*

# Copy the MCM repository and run npm install
COPY . ${MCM_DIR}
WORKDIR $MCM_DIR/catena-sdk
RUN npm install

# Build the catena-sdk
WORKDIR $MCM_DIR/catena-sdk/sdks/cpp/build
RUN cmake -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} .. && \
    make all

# Build MCM SDK
WORKDIR $MCM_DIR/sdk
RUN ./build.sh

## Re-build container for optimised runtime environment using clean Ubuntu release
FROM ubuntu:22.04@sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f

# Set the maintainer label
LABEL maintainer="qiang.han@intel.com"

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ARG http_proxy="http://proxy-mu.intel.com:911"
ARG https_proxy="http://proxy-mu.intel.com:912"
ARG ftp_proxy="http://proxy-mu.intel.com:21"
ARG socks_proxy="http://proxy-mu.intel.com:1080"
ARG no_proxy="intel.com,.intel.com,.cluster.local,.local,localhost,10.0.0.0/8,192.168.0.0/16,127.0.0.0/8,172.16.0.0/12,134.134.0.0/16"

# Update and install package dependencies in a single RUN command to reduce layer size
RUN apt-get update && apt-get full-upgrade -y && \
    apt-get install --no-install-recommends -y libbsd0 && \
    apt-get install --no-install-recommends -y sleepenh net-tools libssl3 wget ffmpeg && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ARG USER=docker
ARG UID=1000
RUN useradd -l -m -s /bin/bash -u ${UID} ${USER}

# Set default environment variables for the user
ARG PREFIX_DIR="/usr/local"
ARG HOME="/home/${USER}"
ENV PATH="${PREFIX_DIR}/bin:$HOME/.local/bin:$HOME/bin:$HOME/usr/bin:$PATH"

# Set the MCM directory
ARG MCM_DIR=/opt/mcm

# Copy the application build from the builder stage
COPY --from=builder $MCM_DIR/sdk/out/samples /opt/mcm-sdk/
COPY --from=builder $MCM_DIR/catena-sdk/sdks/cpp/build/common/examples/full_service/full_service /opt/catena-sdk/
COPY --from=builder $MCM_DIR/catena-sdk/example_device_models /opt/catena-sdk/example_device_models
COPY --from=builder ${PREFIX_DIR} ${PREFIX_DIR}
COPY app/recv.sh app/send.sh /opt/

WORKDIR /opt/assets
RUN wget https://www.larmoire.info/jellyfish/media/jellyfish-3-mbps-hd-hevc.mkv --no-check-certificate && \
    ffmpeg -i jellyfish-3-mbps-hd-hevc.mkv -vframes 8 -c:v rawvideo yuv420p8le_1080p.yuv && \
    ffmpeg -i jellyfish-3-mbps-hd-hevc.mkv -vframes 8 -c:v rawvideo -pix_fmt nv12 yuv_nv12_1080p.yuv \
    ffmpeg -s 1920x1080 -pix_fmt yuv420p -i yuv420p8le_1080p.yuv -pix_fmt yuv422p yuv422p8le_1080p.yuv

RUN ldconfig

USER ${USER}

# Set the working directory
WORKDIR /opt/mcm-sdk/

# Set the default command to run the application
ENTRYPOINT ["./recver_app"]
