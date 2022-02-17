import org.apache.plc4x.java.PlcDriverManager
import org.apache.plc4x.java.api.types.PlcResponseCode
import kotlin.system.exitProcess

fun main() {
    PlcDriverManager().getConnection("modbus://localhost:502").use { conn ->
        if (!conn.metadata.canRead()) {
            println("Cannot read!!")
            return
        }

        val builder = conn.readRequestBuilder()
        builder.addItem("value-1", "coil:1")
        builder.addItem("value-2", "coil:3[4]")
        builder.addItem("value-3", "holding-register:1")
        builder.addItem("value-4", "holding-register:3[4]")
        val readRequest = builder.build()

        val response = readRequest.execute().get()
        for (fieldName in response.fieldNames) {
            if (response.getResponseCode(fieldName) === PlcResponseCode.OK) {
                val numValues = response.getNumberOfValues(fieldName)
                // If it's just one element, output just one single line.
                if (numValues == 1) {
                    println("Value[" + fieldName + "]: " + response.getObject(fieldName))
                } else {
                    println("Value[$fieldName]:")
                    for (i in 0 until numValues) {
                        println(" - " + response.getObject(fieldName, i))
                    }
                }
            } else {
                println("Error[" + fieldName + "]: " + response.getResponseCode(fieldName).name)
            }
        }
    }
    exitProcess(0)
}