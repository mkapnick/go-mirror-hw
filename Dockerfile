FROM golang:1.8

# Create dir in the container and cd
RUN	mkdir -p /usr/go/go-mirror
WORKDIR /usr/go/go-mirror

# copy src code
COPY . .

# change to server dir
WORKDIR /usr/go/go-mirror/server

#RUN go build -o main .

EXPOSE 9000

CMD ["/usr/go/go-mirror/server/main"]
