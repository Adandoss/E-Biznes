FROM ubuntu:24.04

RUN apt-get update && apt-get install -y openjdk-8-jdk wget unzip

RUN wget -q https://github.com/JetBrains/kotlin/releases/download/v1.9.24/kotlin-compiler-1.9.24.zip && \
    unzip -q kotlin-compiler-1.9.24.zip -d /opt && \
    rm kotlin-compiler-1.9.24.zip

ENV PATH="/opt/kotlinc/bin:$PATH"

CMD ["/bin/bash", "-c", "java -version && kotlinc -version"]