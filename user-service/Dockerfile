FROM golang:1.24-alpine

# Install tools yang dibutuhkan
RUN apk add --no-cache git curl

# Install Air (live reload dev tool)
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum dulu (supaya cache Docker optimal)
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

# Copy semua source code
COPY . .

# Expose port aplikasi
EXPOSE 8081

# Jalankan Air
CMD ["air"]
