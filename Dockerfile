From golang:1.13 AS BUILD

ADD . /app/backend
WORKDIR /app/backend

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .


FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=BUILD "/main" ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
