# Use the official Golang image to build the binary
FROM golang AS build-stage

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the Go application source code
COPY *.go ./

# Download dependencies and build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /storm-exporter

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

ENV STORM_UI_HOST=localhost:8081

WORKDIR /

COPY --from=build-stage /storm-exporter /storm-exporter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/storm-exporter"]
