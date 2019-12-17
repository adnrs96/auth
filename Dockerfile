FROM golang:alpine AS builder

COPY . /auth
WORKDIR /auth/cmd/server

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix nocgo -o /server .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /server ./

CMD ["/server"]

EXPOSE 3000
