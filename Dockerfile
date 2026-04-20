FROM golang:1.23 AS build
WORKDIR /app
COPY ./app/go.mod ./app/go.sum ./
RUN go mod download
COPY ./app .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ssrp .

FROM scratch
COPY --from=build /app/ssrp /ssrp
ENTRYPOINT ["/ssrp"]
