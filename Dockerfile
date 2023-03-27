
FROM golang:1.20-alpine AS build_base


RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Unit tests
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
#RUN go build -o ./out/seqnum-apis .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seqnum-apis .


# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

WORKDIR /root/

COPY --from=build_base /app/seqnum-apis .
COPY --from=build_base /app/.env .

# This container exposes port 9000 to the outside world
EXPOSE 9000


# Run seqnum-apis
ENTRYPOINT ["./seqnum-apis"]
