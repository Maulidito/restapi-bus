FROM golang:alpine as restapibus
WORKDIR /usr/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
WORKDIR /usr/src/app/cmd/restapibus
RUN go build -v -o restapibus 
CMD [ "./restapibus" ]