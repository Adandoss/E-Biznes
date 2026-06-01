package com.example

import dev.kord.core.Kord
import dev.kord.core.event.message.MessageCreateEvent
import dev.kord.core.on
import dev.kord.gateway.Intent
import dev.kord.gateway.PrivilegedIntent

import kotlinx.serialization.Serializable
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.plugins.cors.routing.*
import io.ktor.server.plugins.contentnegotiation.*
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.http.*

import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation as ClientContentNegotiation

@Serializable
data class ChatRequest(val message: String)

@Serializable
data class ChatResponse(val response: String)

@Serializable
data class PythonChatRequest(val message: String)

@Serializable
data class PythonChatResponse(val response: String)

@Serializable
data class CategoryDto(val id: Int? = null, val name: String)

@Serializable
data class ProductDto(
    val id: Int? = null,
    val name: String,
    val price: Double,
    val category_id: Int? = null,
    val category: CategoryDto? = null
)

val httpClient = HttpClient(CIO) {
    install(ClientContentNegotiation) {
        json()
    }
}

val pythonServiceUrl = System.getenv("PYTHON_SERVICE_URL") ?: "http://localhost:8000"
val backendUrl = System.getenv("BACKEND_URL") ?: "http://backend:8080"

suspend fun callPythonGptService(message: String): String {
    return try {
        val response: HttpResponse = httpClient.post("$pythonServiceUrl/chat") {
            contentType(ContentType.Application.Json)
            setBody(PythonChatRequest(message))
        }
        val responseBody = response.body<PythonChatResponse>()
        responseBody.response
    } catch (e: Exception) {
        println("Error calling python-service: ${e.message}")
        "Przepraszam, wystąpił problem z połączeniem z serwisem analizującym tekst."
    }
}

suspend fun fetchCategories(): List<String> {
    return try {
        val response: HttpResponse = httpClient.get("$backendUrl/categories")
        val categories = response.body<List<CategoryDto>>()
        categories.map { it.name }
    } catch (e: Exception) {
        println("Error fetching categories: ${e.message}")
        listOf("elektronika", "ksiazki", "jedzenie")
    }
}

suspend fun fetchProductsByCategory(categoryName: String): List<String> {
    return try {
        val response: HttpResponse = httpClient.get("$backendUrl/products")
        val products = response.body<List<ProductDto>>()
        products
            .filter { it.category?.name?.lowercase() == categoryName.lowercase() }
            .map { it.name }
    } catch (e: Exception) {
        println("Error fetching products: ${e.message}")
        emptyList()
    }
}

suspend fun main() {
    val token = System.getenv("BOT_TOKEN")
    
    embeddedServer(Netty, port = 8081) {
        install(CORS) {
            anyHost()
            allowMethod(HttpMethod.Options)
            allowMethod(HttpMethod.Post)
            allowMethod(HttpMethod.Get)
            allowHeader(HttpHeaders.ContentType)
            allowHeader(HttpHeaders.Authorization)
        }
        install(ContentNegotiation) {
            json()
        }
        routing {
            post("/api/chat") {
                try {
                    //3.5 /////////////////////
                    // Obsługa połączenia interfejsu frontendowego JS
                    val req = call.receive<ChatRequest>()
                    val reply = callPythonGptService(req.message)
                    call.respond(ChatResponse(reply))
                } catch (e: Exception) {
                    call.respond(HttpStatusCode.BadRequest, ChatResponse("Błąd bramki Ktor: ${e.localizedMessage}"))
                }
            }
            get("/health") {
                call.respond(mapOf("status" to "ok"))
            }
        }
    }.start(wait = false)
    
    println("Serwer Ktor uruchomiony na porcie 8081")

    if (token.isNullOrEmpty() || token == "DISCORD_BOT_TOKEN_HERE") {
        println("OSTRZEŻENIE: Brak prawidłowego BOT_TOKEN. Bot Discord nie został uruchomiony. Bramka API Ktor działa dalej.")
        while (true) {
            kotlinx.coroutines.delay(5000)
        }
    } else {
        val kord = Kord(token)
        kord.on<MessageCreateEvent> {
            if (message.author?.isBot == true) return@on

            val tresc = message.content
            val trescLower = tresc.lowercase().trim()

            if (trescLower == "!hello") {
                message.channel.createMessage("world!")
                return@on
            }

            if (trescLower == "!kategorie") {
                val categories = fetchCategories()
                val listaKategorii = categories.joinToString(", ")
                message.channel.createMessage("Kategorie: $listaKategorii")
                return@on
            }

            if (trescLower.startsWith("!produkty ")) {
                val zadanaKategoria = trescLower.removePrefix("!produkty ").trim()
                val produkty = fetchProductsByCategory(zadanaKategoria)

                if (produkty.isNotEmpty()) {
                    val listaProduktow = produkty.joinToString(", ")
                    message.channel.createMessage("Produkty w $zadanaKategoria: $listaProduktow")
                } else {
                    message.channel.createMessage("Nie ma produktów w kategorii $zadanaKategoria lub kategoria nie istnieje. Użyj !kategorie")
                }
                return@on
            }

            //3.5 /////////////////////
            // Przekazanie zapytania z bota do Pythona
            val reply = callPythonGptService(tresc)
            message.channel.createMessage(reply)
        }

        println("Bot Discord uruchomiony...")
        kord.login {
            @OptIn(PrivilegedIntent::class)
            intents += Intent.MessageContent
        }
    }
}