FROM golang:1.20.4-alpine as build-backend

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && \
    go mod download -x

COPY . .

# Run tests
RUN go test ./...

RUN CGO_ENABLED=0 go build -o server cmd/*.go

FROM scratch
WORKDIR /app
COPY --from=build-backend /app/server .
COPY --from=build-backend /app/docs/ /app/docs/
COPY --from=build-backend /app/data/ /app/data/

EXPOSE 8080

ENTRYPOINT ["./server"]