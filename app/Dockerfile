FROM docker.io/library/golang:1.16 as builder

WORKDIR /build
COPY go.mod main.go ./
RUN go build -o /build/configtest

FROM gcr.io/distroless/base-debian10
COPY --from=builder /build/configtest /
COPY langs.json /
ENTRYPOINT ["/configtest"]
