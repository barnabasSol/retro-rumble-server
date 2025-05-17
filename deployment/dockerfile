FROM alpine:latest

RUN mkdir /app

COPY bin/retro-rumble.exe /app/retro-rumble

WORKDIR /app

RUN chmod +x /app/retro-rumble

CMD ["/app/retro-rumble"]