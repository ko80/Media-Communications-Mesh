# SPDX-FileCopyrightText: Copyright (c) 2023 Intel Corporation
#
# SPDX-License-Identifier: BSD-3-Clause

FROM ger-is-registry.caas.intel.com/nex-vs-cicd-automation/mcm/ubuntu-2204:latest as builder
LABEL maintainer="qiang.han@intel.com"

ARG http_proxy="http://proxy-mu.intel.com:911"
ARG https_proxy="http://proxy-mu.intel.com:912"
ARG ftp_proxy="http://proxy-mu.intel.com:21"
ARG socks_proxy="http://proxy-mu.intel.com:1080"
ARG no_proxy="intel.com,.intel.com,.cluster.local,.local,localhost,10.0.0.0/8,192.168.0.0/16,127.0.0.0/8,172.16.0.0/12,134.134.0.0/16"

ARG USER=docker
ARG UID=1001

ENV DEBIAN_FRONTEND noninteractive

SHELL ["/bin/bash", "-ex", "-o", "pipefail", "-c"]
# Install package dependencies
RUN apt-get update --fix-missing && \
    apt-get full-upgrade -y && \
    apt-get install --no-install-recommends -y apt-utils \
    software-properties-common ca-certificates \
    libtool autoconf automake pkg-config \
    sudo jq \
    python3 python3-pip \
    cmake nasm \
    build-essential gcc meson \
    libnuma-dev libjson-c-dev libpcap-dev libgtest-dev \
    libsdl2-dev libsdl2-ttf-dev libssl-dev nlohmann-json3-dev \
    libbsd-dev \
    python3-pyelftools libgrpc++-dev protobuf-compiler \
    curl git git-lfs \
    vim openssh-server clang-format \
    gdb psmisc nodejs npm  \
    pciutils iproute2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -l -m -s /bin/bash -u ${UID} ${USER}

ARG PREFIX_DIR="/usr/local"
ARG MCM_DIR=/opt/mcm
ARG MTL_DIR=/opt/mcm/mtl
ARG DPDK_DIR=/opt/dpdk
ARG JPEG_XS_PATH=/opt/mcm/svt-jpeg-xs
ARG JWT_CPP_DIR="/opt/jwt-cpp"
ENV DPDK_VER=23.11

# jwt-cpp requirement for catena-sdk
RUN git clone https://github.com/Thalhammer/jwt-cpp ${JWT_CPP_DIR}
WORKDIR ${JWT_CPP_DIR}/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} ..  && \
    make && \
    make install

# Set default environment variables for the user
ARG HOME="/home/${USER}"
ENV PATH="${PREFIX_DIR}/bin:$HOME/.local/bin:$HOME/bin:$HOME/usr/bin:$PATH"
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib/x86_64-linux-gnu:${LD_LIBRARY_PATH}"

# Build DPDK/IMTL
COPY mtl ${MTL_DIR}/
WORKDIR ${DPDK_DIR}
RUN git clone --depth 1 --branch v${DPDK_VER} https://github.com/DPDK/dpdk.git . && \
    find ${MTL_DIR}/patches/dpdk/${DPDK_VER}/ -maxdepth 1 -name '*.patch' -print0 | sort -z | xargs -t -0 -n 1 patch -d ${DPDK_DIR} -p1 -i && \
    meson build && \
    ninja -C build && \
    ninja -C build install && \
    pkg-config --cflags libdpdk && \
    pkg-config --libs libdpdk && \
    pkg-config --modversion libdpdk

# Build IMTL
WORKDIR ${MTL_DIR}
RUN ./build.sh

# gRPC
ARG GPRC_VERSION="v1.58.0"
RUN git clone --recurse-submodules -b ${GPRC_VERSION} --depth 1 --shallow-submodules https://github.com/grpc/grpc /tmp/grpc && \
    mkdir -p /tmp/grpc/cmake/build
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

# Build and install JPEG-XS
COPY svt-jpeg-xs ${JPEG_XS_PATH}/
WORKDIR ${JPEG_XS_PATH}/Build/linux
RUN ./build.sh install

WORKDIR ${JPEG_XS_PATH}/imtl-plugin
RUN ./build.sh && \
    mkdir -p /opt/tools/ && \
    find ${MTL_DIR}/build/app -maxdepth 1 -type f -perm /a+x -exec cp '{}' -t /opt/tools/ \; && \
    cp ${JPEG_XS_PATH}/imtl-plugin/kahawai.json /opt/tools/

# Build the catena-sdk
COPY catena-sdk ${MCM_DIR}/catena-sdk
WORKDIR ${MCM_DIR}/catena-sdk/sdks/cpp/build
RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${PREFIX_DIR} .. && \
    make all

# Install Conan using pip3 package installer.
COPY media-proxy ${MCM_DIR}/media-proxy/
RUN python3 -m pip install --no-cache-dir --require-hashes -r /opt/mcm/media-proxy/requirements.txt

# Build MCM SDK
COPY sdk ${MCM_DIR}/sdk
WORKDIR ${MCM_DIR}/sdk
RUN ./build.sh

# Build MCM
COPY build* CMake* ${MCM_DIR}/
WORKDIR ${MCM_DIR}
RUN rm -rf out && \
    ./build.sh

# Re-build container for optimised runtime environment using clean Ubuntu release
FROM ger-is-registry.caas.intel.com/nex-vs-cicd-automation/mcm/ubuntu-2204:latest
LABEL maintainer="milosz.linkiewicz@intel.com"

ENV KAHAWAI_CFG_PATH=/opt/tools/kahawai.json
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Warsaw

ARG USER=ubuntu
ARG UID=1000

ARG http_proxy="http://proxy-mu.intel.com:911"
ARG https_proxy="http://proxy-mu.intel.com:912"
ARG ftp_proxy="http://proxy-mu.intel.com:21"
ARG socks_proxy="http://proxy-mu.intel.com:1080"
ARG no_proxy="intel.com,.intel.com,.cluster.local,.local,localhost,10.0.0.0/8,192.168.0.0/16,127.0.0.0/8,172.16.0.0/12,134.134.0.0/16"

# Set default environment variables for the user
ARG PREFIX_DIR="/usr/local"
ENV PATH="${PREFIX_DIR}/bin:$HOME/.local/bin:$HOME/bin:$HOME/usr/bin:$PATH"

# Set the MCM directory
ARG MCM_DIR=/opt/mcm

WORKDIR /opt/assets
RUN apt-get update --fix-missing && \
    apt-get full-upgrade -y && \
    apt-get install -y --no-install-recommends \
        libsdl2-2.0-0 libsdl2-ttf-2.0-0 jq \
        libnuma1 libjson-c5 libpcap0.8 libssl3 \
        libatomic1 libbsd0 ffmpeg pciutils iproute2 \
        wget sleepenh net-tools vim sudo && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    wget https://www.larmoire.info/jellyfish/media/jellyfish-3-mbps-hd-hevc.mkv --no-check-certificate && \
    ffmpeg -i jellyfish-3-mbps-hd-hevc.mkv -vframes 8 -c:v rawvideo -pix_fmt yuv420p yuv420p8le_1080p.yuv && \
    ffmpeg -i jellyfish-3-mbps-hd-hevc.mkv -vframes 8 -c:v rawvideo -pix_fmt nv12 yuv_nv12_1080p.yuv && \
    ffmpeg -s 1920x1080 -pix_fmt yuv420p -i yuv420p8le_1080p.yuv -pix_fmt yuv422p yuv422p8le_1080p.yuv

# Copy the application build from the builder stage
COPY --from=builder /opt/tools /opt/tools/
COPY --from=builder $MCM_DIR/sdk/out/samples /opt/mcm-sdk/
COPY --from=builder $MCM_DIR/catena-sdk/sdks/cpp/build/common/examples/full_service/full_service /opt/catena-sdk/
COPY --from=builder $MCM_DIR/catena-sdk/example_device_models /opt/catena-sdk/example_device_models
COPY --from=builder ${PREFIX_DIR} ${PREFIX_DIR}
COPY app/tx_r02s15_1port_2v.json /opt/tools/
COPY media-proxy/imtl.json /opt/

RUN ldconfig

# Expose correct default ports to allow quick publishing
EXPOSE 8001 8002

USER ${USER}
WORKDIR /opt

ENTRYPOINT ["media_proxy"]
