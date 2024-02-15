FROM golang:alpine as builder
ENV GO111MODULE=on
LABEL maintainer="dmilyano"
RUN apk update && apk add --no-cache git
WORKDIR /example1
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main ./cmd/. #

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /example1/main .
COPY --from=builder /example1/config.yml .
COPY --from=builder /example1/.env .
EXPOSE 8081
CMD ["./main"]