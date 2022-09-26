FROM golang:1.19.1-alpine3.16 as build

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY testio /

CMD [ "/testio" ]