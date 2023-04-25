FROM alpine:3.14
RUN mkdir "/app"
WORKDIR "/app"
COPY Evo /app/app
COPY config.yml /app/config.yml
CMD ./app

