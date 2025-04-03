# Step 1 : Build application
FROM golang:1.20 AS builder

# Define build workdir
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code in the container
COPY . .

# Build Go executable
RUN go build -o gqmqtt .

# Step 2 : Create final image for execution
FROM debian:testing-slim

# Create dialout group and a user goapp within this group
RUN useradd -m -g dialout -u 1000 goapp

# Install dependencies and tools
RUN apt-get update && apt-get install -y \
    udev libc6\
    && rm -rf /var/lib/apt/lists/*

# Copy executable
COPY --from=builder /app/gqmqtt /usr/local/bin/

# Défine workdir
WORKDIR /app

# Utiliser l'utilisateur non privilégié 'goapp' du groupe 'dialout'
USER goapp

# Lancer l'application
CMD ["gqmqtt"]
