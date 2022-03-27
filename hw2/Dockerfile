FROM golang:1.17 as build

WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go build -o ./out/k8_test cmd/k8_test/main.go

FROM centos:7
WORKDIR /
COPY --from=build /app/out /app
EXPOSE 8080
ENV PORT=:8080

ENTRYPOINT ["/app/k8_test"]



