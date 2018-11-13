FROM golang:latest as builder
RUN mkdir /build 
ADD ./src /build/src
ENV GOPATH $GOPATH:/build
WORKDIR /build/src
#RUN apk add git 
RUN go get -d ./...
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/src/main /app/
COPY --from=builder /build/src/config.json /app/
WORKDIR /app
CMD ["./main"]