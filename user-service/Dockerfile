FROM golang

RUN mkdir app
WORKDIR /app
COPY ./app .
ADD migrations ./migrations
ADD Makefile .

CMD ["./app", "server"]

EXPOSE 8080