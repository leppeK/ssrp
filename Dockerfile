FROM golang AS build
WORKDIR /app
COPY ./app /app
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
WORKDIR /app
COPY --from=build /app/ssrp /app/ssrp
ENTRYPOINT ["./ssrp"]
