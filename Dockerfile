FROM alpine:3.16
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD lightscreen /lightscreen
ADD examples/spectral-notary/admission.yaml /lightscreen.yaml
ENTRYPOINT ["./lightscreen"]