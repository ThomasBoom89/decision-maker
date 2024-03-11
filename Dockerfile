FROM golang:1.22-alpine3.19 AS Builder

WORKDIR /app
COPY . .

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -o decision-maker ./src/main.go

FROM node:20 AS FE-Builder
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

FROM scratch AS Production

USER 1000

WORKDIR /app

COPY --from=Builder /app/decision-maker /app
COPY --from=FE-Builder /app/views /app/views
COPY --from=FE-Builder /app/public /app/public

EXPOSE 3000

CMD ["./decision-maker"]

FROM golang:1.22-alpine3.19 AS Development

RUN go install github.com/cosmtrek/air@v1.51.0
RUN go install github.com/go-delve/delve/cmd/dlv@v1.22.1

WORKDIR /app

CMD ["air"]

FROM node:20.11.0-slim AS Frontend-Development

WORKDIR /app

CMD ["npm", "run", "dev"]
