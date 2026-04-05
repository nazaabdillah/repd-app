# Stage 1: Build Engine (Menggunakan Alpine Linux yang sangat ringan)
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy semua file kode
COPY . .

# Download dependensi module
RUN go mod download

# Build binary Golang secara statis
RUN CGO_ENABLED=0 GOOS=linux go build -o repd-api .

# Stage 2: Production Image (Hanya berisi binary, tanpa source code)
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage 1
COPY --from=builder /app/repd-api .

# Buat folder uploads agar fungsi simpan foto tidak error/panic
RUN mkdir -p uploads/photos

# Ekspos port default
EXPOSE 8080

# Eksekusi binary
CMD ["./repd-api"]