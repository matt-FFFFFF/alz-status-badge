FROM golang:1.18-alpine AS build

WORKDIR /src
COPY go.* main.go /src/
RUN CGO_ENABLED=0 go build -o /bin/badge

###

FROM scratch

COPY --from=build /bin/badge /bin/badge
ENTRYPOINT [ "/bin/badge" ]
