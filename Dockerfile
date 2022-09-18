FROM alpine:3.16 as builder

RUN apk --no-cache add go

ADD . /app

WORKDIR /app

RUN go build



FROM alpine:3.16 as final

COPY --from=builder /app/multiple-upload-to-s3 /app/muts

WORKDIR /app

CMD ["/app/muts"]
