FROM golang:latest as builder
WORKDIR /go/src/github.com/alekssaul/pcf2cloudDNS
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o broker ./cmd/broker
RUN mkdir /app && \
	mv broker /app && \
    mv configs/router.yaml /app

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app .
CMD /app/broker
EXPOSE 8080