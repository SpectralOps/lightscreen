FROM alpine:3.18.4
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD lightscreen /lightscreen
ADD examples/spectral-notary/admission.yaml /lightscreen.yaml
ENTRYPOINT ["./lightscreen"]