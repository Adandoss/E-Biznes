import java.sql.DriverManager

fun main() {
    println("Hello World!")

    try {
        val connection = DriverManager.getConnection("jdbc:sqlite:test.db")
        println("Success - database connected!!!")
        connection.close()
    } catch (e: Exception) {
        println("Error: \${e.message}")
    }
}