FROM golang:1.21.5 as build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal

RUN CGO_ENABLED=0 go build -o ords ./cmd/main.go

FROM scratch

COPY --from=build /build/ords /ords

ENTRYPOINT ["/ords"]
