# backend
FROM golang:1.16 as backend
WORKDIR /src
COPY ./backend/go.sum ./backend/go.mod ./
COPY ./backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app .

# frontend
FROM node:12 AS frontend
WORKDIR /app
COPY ./frontend/package.json ./
COPY ./frontend/package-lock.json ./
RUN npm install
COPY ./frontend ./
RUN npm run build

# Build final image
FROM alpine as final
WORKDIR /go
RUN mkdir -p /go/public
COPY --from=frontend /app/dist /go/public
COPY --from=backend /bin/app /go/app
EXPOSE 8000
ENTRYPOINT /go/app 