FROM gradle:jdk25 AS builder
WORKDIR /app
COPY . .
RUN gradle clean installDist --no-daemon

FROM eclipse-temurin:25-jre-alpine
WORKDIR /app
COPY --from=builder /app/build/install/Ktor-bot/ /app/
RUN chmod +x /app/bin/Ktor-bot
CMD ["/app/bin/Ktor-bot"]