FROM golang:1.21-alpine3.19 AS build

RUN apk --no-cache add curl

WORKDIR /wl

COPY . ./

RUN ./install-tools.sh

RUN SWAGGER_GENERATE_EXTENSION="false" ./go-swagger generate spec -o swagger.json

RUN go build -o work-log cmd/init.go cmd/main.go


FROM alpine:3.19

LABEL maintainer="matthias@kellnhofer.com"

RUN addgroup --system --gid 1000 wl-group && \
  adduser --system --uid 1000 --home /app --ingroup wl-group wl-user

RUN mkdir -p /app && \
  chown -R wl-user:wl-group /app && \
  chmod -R 751 /app

WORKDIR /app

USER 1000

COPY config/localizations ./config/localizations
COPY scripts ./scripts
COPY static/resources ./static/resources
COPY static/swagger-ui ./static/swagger-ui
COPY web/templates ./web/templates
COPY --from=build /wl/swagger.json ./static/swagger-ui
COPY --from=build /wl/work-log ./work-log

EXPOSE 8080/tcp

ENTRYPOINT ["./work-log"]