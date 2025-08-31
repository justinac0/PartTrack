###############################################################################
# build the application binary from a full go image
###############################################################################
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# TODO: try to stop this step being re-run every time something in the static folder changes
COPY . .
RUN go tool templ generate
RUN go build -o parttrack ./cmd/app

###############################################################################
# build an optimized, minimal image for the application runtime
###############################################################################
FROM scratch AS production
WORKDIR /app

COPY --from=builder /app/parttrack /usr/local/bin/parttrack
COPY ./static /app/static

EXPOSE 8000
ENTRYPOINT [ "/usr/local/bin/parttrack" ]
