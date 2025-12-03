FROM oven/bun:1 AS frontend-builder
WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

RUN bun run build:assets
RUN bun run build:css

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/static ./static

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o regex-checker .

FROM alpine:latest

WORKDIR /root/

COPY --from=backend-builder /app/regex-checker .

ARG port=80
ENV PORT=$port

EXPOSE $PORT

CMD ["./regex-checker"]
