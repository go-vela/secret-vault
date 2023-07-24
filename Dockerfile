# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

##############################################################################
##    docker build --no-cache --target certs -t secret-vault:certs .    ##
##############################################################################

FROM alpine@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1 as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##      docker build --no-cache -t secret-vault:local .      ##
###############################################################

FROM alpine:latest@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/secret-vault /bin/secret-vault

ENTRYPOINT ["/bin/secret-vault"]