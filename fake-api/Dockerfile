From golang:1.13
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .


# FROM alpine
# RUN apk --no-cache add ca-certificates
# COPY --from=BUILDFAKE "/main" ./
# RUN chmod +x ./main
# WORKDIR .
# ENTRYPOINT ["./main"]

ENTRYPOINT ["./main"]
EXPOSE 8080
