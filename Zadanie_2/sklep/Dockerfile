FROM sbtscala/scala-sbt:eclipse-temurin-jammy-17.0.10_7_1.9.9_2.13.13

WORKDIR /app
COPY . .

RUN sbt stage

EXPOSE 9000

CMD ["target/universal/stage/bin/sklep", "-Dplay.http.secret.key=super_sekretny_klucz_do_zadania_12345", "-Dhttp.port=9000", "-Dplay.filters.hosts.allowed.0=."]