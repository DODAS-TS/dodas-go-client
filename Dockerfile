FROM golang:alpine as BUILD

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  -o dodas .

FROM alpine as APP

COPY --from=0 /app/dodas /usr/bin/dodas

ENTRYPOINT ["dodas"]