FROM golang:1.17 as build-env

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY *.go . 

RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static

COPY --from=build-env /app/aws-ecr-helper /
CMD ["/aws-ecr-helper"]