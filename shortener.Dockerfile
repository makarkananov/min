FROM ubuntu:latest
LABEL authors="portisch_barrel"

ENTRYPOINT ["top", "-b"]