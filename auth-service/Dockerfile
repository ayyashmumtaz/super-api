# Gunakan base image golang alpine
FROM golang:1.24-alpine

# Install tools yang dibutuhkan
RUN apk add --no-cache git curl

# Install Air (live reload tool)
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

# Set working directory
WORKDIR /app

# Copy file go.mod dan go.sum dulu (biar caching efisien)
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

# Copy semua source code
COPY . .

# Expose port yang digunakan oleh auth service
EXPOSE 8080

# Jalankan pakai Air
CMD ["air"]
