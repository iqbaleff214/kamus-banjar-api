FROM golang:alpine3.20 AS build

WORKDIR /build
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -tags=fs -o server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /build/server .
# json files dictionary
COPY data/ data/
ENV PORT=":8001"
ENTRYPOINT ["/app/server"]
