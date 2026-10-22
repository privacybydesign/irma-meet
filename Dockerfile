FROM golang:1.26.10 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# No cgo, so the binary runs on the libc-free distroless static image.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/irma-meet .

# -----------------------------------------------------

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=build /out/irma-meet /app/irma-meet

# Declared explicitly to satisfy static scanners; the :nonroot tag already
# defaults to this same user (UID 65532).
USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/app/irma-meet"]
