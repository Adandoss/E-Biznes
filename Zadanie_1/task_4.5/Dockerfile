FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y openjdk-8-jdk wget unzip && \
    rm -rf /var/lib/apt/lists/*

RUN wget -q https://github.com/JetBrains/kotlin/releases/download/v1.9.24/kotlin-compiler-1.9.24.zip && \
    unzip -q kotlin-compiler-1.9.24.zip -d /opt && \
    rm kotlin-compiler-1.9.24.zip

RUN wget -q https://services.gradle.org/distributions/gradle-8.12.1-bin.zip && \
    unzip -q gradle-8.12.1-bin.zip -d /opt && \
    rm gradle-8.12.1-bin.zip

ENV PATH="/opt/kotlinc/bin:/opt/gradle-8.12.1/bin:${PATH}"

WORKDIR /app

COPY build.gradle .
COPY src ./src

RUN gradle build --no-daemon

CMD ["gradle", "run", "--no-daemon"]
