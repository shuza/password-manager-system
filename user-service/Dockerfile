FROM golang

RUN mkdir app
WORKDIR /app
COPY ./app .
ADD migrations ./migrations

ENTRYPOINT ["./app", "migrate", "&&", "./app", "server"]

EXPOSE 8080