FROM alpine:3.13
COPY redminebot /redminebot
RUN chmod +x /redminebot
ENTRYPOINT ["/redminebot"]
