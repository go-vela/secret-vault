# SPDX-License-Identifier: Apache-2.0

##############################################################################
##    docker build --no-cache --target certs -t secret-vault:certs .    ##
##############################################################################

FROM alpine@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as certs

RUN apk add --update --no-cache ca-certificates

###############################################################
##      docker build --no-cache -t secret-vault:local .      ##
###############################################################

FROM alpine:latest@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/secret-vault /bin/secret-vault

ENTRYPOINT ["/bin/secret-vault"]