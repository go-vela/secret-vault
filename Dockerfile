# SPDX-License-Identifier: Apache-2.0

##############################################################################
##    docker build --no-cache --target certs -t secret-vault:certs .    ##
##############################################################################

FROM alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##      docker build --no-cache -t secret-vault:local .      ##
###############################################################

FROM alpine:latest@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/secret-vault /bin/secret-vault

ENTRYPOINT ["/bin/secret-vault"]