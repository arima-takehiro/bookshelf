#FROM golang AS builder

#WORKDIR /app

#COPY . .

#RUN CGO_ENABLED=0 go build -o bookshelf .

# ------------------------

#FROM alpine AS bookshelf

#WORKDIR /app

#COPY --from=builder /app/bookshelf .

#ENV GOOGLE_CLOUD_PROJECT ca-willsbooster-test
#ENV GOOGLE_APPLICATION_CREDENTIALS /app/ca-willsbooster-test-db2897beab98.json

#ENTRYPOINT ["./bookshelf"]

