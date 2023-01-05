FROM golang:1.19.4-alpine
#to do: change it 1.15 one
COPY ./release/server /bin/
COPY config/file/ /config/file/
WORKDIR /
# needs a volume for /config files to run in a production mode