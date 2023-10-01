# Build React app
FROM node:18-alpine AS react-build
WORKDIR /app/frontend
COPY ./frontend/package.json ./frontend/package-lock.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build Go app
FROM golang:1.20-alpine AS go-build
WORKDIR /app/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o fibonacci .

# Run the app
FROM alpine:3.14.2
EXPOSE 8080
EXPOSE 8081
COPY --from=react-build /app/frontend/build /app/frontend/build
COPY --from=go-build /app/backend/fibonacci /app
CMD ["./app/fibonacci"]

