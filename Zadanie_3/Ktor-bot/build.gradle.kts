plugins {
    kotlin("jvm") version "2.3.10"
    application
}

application {
    mainClass.set("com.example.MainKt")
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation("dev.kord:kord-core:0.18.1")
    implementation("org.slf4j:slf4j-simple:2.0.16")
}
kotlin {
    jvmToolchain(25)
}

tasks.test {
    useJUnitPlatform()
}