FROM alpine:latest

ENV APP diskstat

COPY diskstat diskstat

RUN chmod +x diskstat

EXPOSE 8080

CMD ["./diskstat"]
