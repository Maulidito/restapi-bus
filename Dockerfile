FROM golang:alpine as restapibus
WORKDIR /usr/src/app

ENV PORT=:8080
ENV SECRET_KEY_AUTH = MY_SECRET_KEY_GO_REST_API_BUS
ENV USERNAME_DB = root
ENV PASSWORD_DB = root
ENV NAME_DB = db_bus
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o restapibus
CMD [ "./restapibus" ]