FROM ubuntu:22.04

RUN apt-get update && apt-get install -y python3.10

RUN ln -s /usr/bin/python3.10 /usr/bin/python3 && \
    ln -s /usr/bin/python3.10 /usr/bin/python

CMD ["python3", "--version"]