FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build/
RUN go mod tidy
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -tags "mongo local" -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /build/config.yml /app/
COPY --from=builder /build/users.json /app/
WORKDIR /app
CMD ["./main"]
EXPOSE 8082
