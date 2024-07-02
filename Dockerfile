FROM golang:1.20.14-alpine3.19


# Add Maintainer Info
LABEL maintainer="Arief Maffrudin"

WORKDIR /app
COPY . .
RUN go build 

EXPOSE 8080

# Run the executable
CMD ["./main"]