FROM golang:1.17.3-alpine as build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o build/app cmd/main/main.go

FROM alpine
WORKDIR /
COPY --from=build /app/build/app /app
ENTRYPOINT ["./app"]