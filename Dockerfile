FROM node:20-alpine as client-builder

RUN mkdir /app
ADD client/package.json /app/client/package.json
ADD client/package-lock.json /app/client/package-lock.json

WORKDIR /app/client
RUN npm install

ADD client /app/client
RUN npm run lint && npm run build



FROM golang:1.22-alpine as go-builder
RUN mkdir /app
ADD go.mod /app/go.mod
ADD go.sum /app/go.sum

WORKDIR /app
RUN go mod download

ADD . /app
RUN go test ./... && \
    go build -o /app/.bin/server /app/cmd/server && \
    go build -o /app/.bin/migrate /app/cmd/migrate


FROM alpine:3
RUN mkdir /app
WORKDIR /app

COPY --from=client-builder /app/client/dist /app/client/dist
COPY --from=go-builder /app/.bin/server /app/server
COPY --from=go-builder /app/.bin/migrate /app/migrate

EXPOSE 8000 8001 8081

RUN adduser -D app
USER app

CMD ["/app/server"]
