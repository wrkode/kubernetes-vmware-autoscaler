# Copyright 2019 Frederic Boltz Author. All rights reserved
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#ARG BASEIMAGE=gcr.io/distroless/static:latest-amd64
FROM alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM

COPY out .
RUN ls / ; mv /$TARGETPLATFORM/vsphere-autoscaler /vsphere-autoscaler

FROM ubuntu:focal

LABEL maintainer="Frederic Boltz <frederic.boltz@gmail.com>"

COPY --from=builder /vsphere-autoscaler /usr/local/bin/vsphere-autoscaler

EXPOSE 5200

CMD ["/usr/local/bin/vsphere-autoscaler"]
