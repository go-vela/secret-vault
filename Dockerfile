# SPDX-License-Identifier: Apache-2.0

##############################################################################
##    docker build --no-cache --target certs -t secret-vault:certs .    ##
##############################################################################

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##      docker build --no-cache -t secret-vault:local .      ##
###############################################################

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/secret-vault /bin/secret-vault

ENTRYPOINT ["/bin/secret-vault"]