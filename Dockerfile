FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /hostname

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /hostname /hostname

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/hostname"]
EXPOSE 8000

CMD [ "/hostname" ]
