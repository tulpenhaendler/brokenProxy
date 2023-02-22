# Start from a base image
FROM golang:1.17.5-alpine3.15

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the app
RUN go build -o app .

# Expose the port on which the app will listen
EXPOSE 8080

# Start the app
CMD ["/app/app"]
