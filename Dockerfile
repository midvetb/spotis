FROM golang:1.26-alpine3.23 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /go/bin/spotis ./cmd/app/

FROM gcr.io/distroless/static-debian13
COPY --from=build /go/bin/spotis /spotis
CMD ["/spotis"]
