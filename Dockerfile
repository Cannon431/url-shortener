FROM golang

RUN go get gopkg.in/mgo.v2 && go get github.com/gorilla/mux && go get github.com/joho/godotenv && mkdir -p /go/src/app

WORKDIR /go/src/app

CMD [ "go", "run", "main.go" ]
