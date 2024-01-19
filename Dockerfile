FROM golang:1.21-alpine3.19 AS Builder

WORKDIR /app
COPY . .

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -o decision-maker ./src/main.go

FROM node:20 AS FE-Builder
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

FROM scratch
WORKDIR /app

COPY --from=Builder /app/decision-maker /app
COPY --from=FE-Builder /app/views /app/views
COPY --from=FE-Builder /app/public /app/public

CMD ["./decision-maker"]
