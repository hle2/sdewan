# SPDX-license-identifier: Apache-2.0
##############################################################################
# Copyright (c) 2020
# All rights reserved. This program and the accompanying materials
# are made available under the terms of the Apache License, Version 2.0
# which accompanies this distribution, and is available at
# http://www.apache.org/licenses/LICENSE-2.0
##############################################################################

FROM golang:1.14.1

WORKDIR /go/src/github.com/central-controller
COPY ./ ./
RUN cd src/scc && make all

FROM ubuntu:16.04

WORKDIR /opt/scc
RUN groupadd -r scc && useradd -r -g scc scc
RUN chown scc:scc /opt/scc -R
COPY --chown=scc --from=0 /go/src/github.com/central-controller/src/scc/scc ./
COPY --chown=scc --from=0 /go/src/github.com/central-controller/build/entrypoint ./
RUN chmod +x entrypoint

USER scc

ENTRYPOINT ["./entrypoint"]