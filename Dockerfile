FROM golang:1.15 AS build
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM registry.il2.dso.mil/platform-one/devops/pipeline-templates/node-alpine AS prod
COPY --from=build /app .
CMD ["./app"]