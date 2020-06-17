FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
ENV PORT 8080
COPY --from=builder /build/main /app/
COPY --from=builder /build/datastore/scib-data-account.json /app/scib-data-account.json
ENV GOOGLE_APPLICATION_CREDENTIALS /app/scib-data-account.json
WORKDIR /app
CMD ["./main"]
EXPOSE ${PORT}