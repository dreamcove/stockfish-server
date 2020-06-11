FROM golang:latest

MAINTAINER Chris Watson (chris@dreamcove.com)

USER root

COPY go.mod ./src
COPY main.go ./src

RUN apt-get update
RUN apt-get -y install unzip

RUN curl https://stockfishchess.org/files/stockfish-11-linux.zip -o stockfish-11-linux.zip
RUN unzip stockfish-11-linux.zip stock*64
RUN cp stockfish-11-linux/Linux/*_x64 ./stockfish_x64
RUN rm -Rf stockfish-11-linux*
RUN chmod a+x ./stockfish_x64

ENV STOCKFISH_PATH=./stockfish_x64

EXPOSE 8081

RUN cd src && go build
RUN mv src/stockfish-server .


ENTRYPOINT ./stockfish-server
