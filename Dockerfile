FROM golang:1.19.1-alpine3.16

ADD testio /testio

ENTRYPOINT [ "/testio" ]