package com.example

import dev.kord.core.Kord
import dev.kord.core.event.message.MessageCreateEvent
import dev.kord.core.on
import dev.kord.gateway.Intent
import dev.kord.gateway.PrivilegedIntent

suspend fun main() {
    val token = System.getenv("BOT_TOKEN") ?: throw IllegalArgumentException("Ustaw ten token")
    val kord = Kord(token)

    kord.on<MessageCreateEvent> {
        val tresc = message.content.lowercase()

        if (tresc == "!hello") {
            message.channel.createMessage("wolrd!")
        }
    }

    println("Działa chyba")
    kord.login {
        @OptIn(PrivilegedIntent::class)
        intents += Intent.MessageContent
    }
}