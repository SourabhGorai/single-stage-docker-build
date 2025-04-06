# Single-stage Dockerfile
FROM ubuntu AS build

RUN apt-get update && apt-get install -y golang-go

ENV GO111MODULE=off

COPY . .

RUN CGO_ENABLED=0 go build -o /myapp .

# Run the binary when the container starts
ENTRYPOINT ["./myapp"]

