FROM golang

RUN mkdir app
WORKDIR /app
COPY ./app .

ENTRYPOINT ["./app", "server"]

EXPOSE 8080