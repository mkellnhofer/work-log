FROM --platform=$BUILDPLATFORM golang:1.25-alpine3.22 AS build

ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add curl

WORKDIR /wl

COPY . ./

RUN ./install-tools.sh

RUN SWAGGER_GENERATE_EXTENSION="false" go-swagger generate spec -o swagger.json

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o work-log cmd/init.go cmd/main.go


FROM alpine:3.22

LABEL maintainer="matthias@kellnhofer.com"

RUN addgroup --system --gid 1000 wl-group && \
  adduser --system --uid 1000 --home /app --ingroup wl-group wl-user

RUN mkdir -p /app && \
  chown -R wl-user:wl-group /app && \
  chmod -R 751 /app

WORKDIR /app

USER 1000

COPY resources/db ./resources/db
COPY resources/localizations ./resources/localizations
COPY resources/static ./resources/static
COPY --from=build /wl/swagger.json ./resources/static/swagger-ui
COPY --from=build /wl/work-log ./work-log

EXPOSE 8080/tcp

ENTRYPOINT ["./work-log"]