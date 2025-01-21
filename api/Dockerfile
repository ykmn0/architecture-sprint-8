FROM golang:1.23

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /app

EXPOSE 8000

# Run
CMD ["/app"]
