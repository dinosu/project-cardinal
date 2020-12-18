FROM golang:latest AS builder
RUN go build -o /go/bin/app main.go

FROM gcr.io/distroless/base-debian10
ENV PORT 8080
COPY --from=build /go/bin/app /
CMD ["/app"]
