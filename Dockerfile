FROM ubuntu:latest
LABEL authors="aszotov"

ENTRYPOINT ["top", "-b"]