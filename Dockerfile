FROM golang:1.17

WORKDIR /go/src/app
COPY . .

# ENV GIN_MODE=release
RUN go get -d -v ./...
RUN go install -v ./...

RUN mkdir "/notes"
RUN ln -sfr /media media
VOLUME /notes

EXPOSE 8080

CMD ["mark-notes-server"]
