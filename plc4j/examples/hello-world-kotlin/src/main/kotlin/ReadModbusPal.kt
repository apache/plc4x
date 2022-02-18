import org.apache.plc4x.java.PlcDriverManager
import org.apache.plc4x.java.api.types.PlcResponseCode
import java.util.concurrent.TimeUnit
import kotlin.system.exitProcess

fun main() {
    PlcDriverManager()
        .getConnection("modbus://localhost:502")
        .use { conn ->
            if (!conn.metadata.canRead()) {
                println("Cannot read!!")
                return
            }

            val readRequest = conn.readRequestBuilder()
                .addItem("value-1", "coil:1")
                .addItem("value-2", "coil:3[4]")
                .addItem("value-3", "holding-register:1")
                .addItem("value-4", "holding-register:3[4]")
                .build()

            val response = readRequest.execute().get(1, TimeUnit.MINUTES)
            response.fieldNames.forEach { fieldName ->
                if (response.getResponseCode(fieldName) !== PlcResponseCode.OK) {
                    println("Error[$fieldName]: $response.getResponseCode(fieldName).name")
                    return
                }
                val numValues = response.getNumberOfValues(fieldName)
                // If it's just one element, output just one single line.
                if (numValues == 1) {
                    println("Value[$fieldName]: response.getObject(fieldName)")
                } else {
                    println("Value[$fieldName]:")
                    for (i in 0 until numValues) {
                        println(" - " + response.getObject(fieldName, i))
                    }
                }
            }
        }

    exitProcess(0)
}