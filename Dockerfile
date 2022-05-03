FROM golang:1.18-alpine AS build

WORKDIR /src
COPY . /src/
RUN apk --no-cache add ca-certificates
RUN CGO_ENABLED=0 go build -o /bin/badge

###

FROM scratch

COPY --from=build /bin/badge /bin/badge
COPY --from=build /etc/ssl /etc/ssl/

ENTRYPOINT [ "/bin/badge" ]s
