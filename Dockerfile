FROM golang:latest as build
WORKDIR /build
COPY . .
RUN go build -o go-greeting ./app 
RUN ls

FROM ubuntu:latest as run
WORKDIR /app
COPY --from=build /build/go-greeting .
ENTRYPOINT [ "/app/go-greeting" ]