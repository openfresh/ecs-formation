FROM openfresh/golang:1.9.2 AS build

WORKDIR /go/src/github.com/openfresh/ecs-formation
COPY . . 
RUN make deps
RUN make build

FROM gliderlabs/alpine:3.6
RUN apk --no-cache add ca-certificates openssl
COPY --from=build /go/src/github.com/openfresh/ecs-formation/bin/ecs-formation /usr/local/bin/ 
