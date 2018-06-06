import java.io.*


File outputDir = new File(basedir, "target/docs")

String[] expectedFiles = ['sample.html', 'invoker.properties', 'StringUtils.java']

// output files should be copied
for (String expectedFile : expectedFiles) {
    File file = new File(outputDir, expectedFile)
    println("Checking for existence of " + file)
    if (!file.isFile()) {
        throw new Exception("Missing file " + file)
    }
}

String fileText = new File(outputDir, 'StringUtils.java').text
// properties in 'StringUtils.java' should be filtered (replaced)
['java.version', 'command.property', 'pom.property'].each {
    if (fileText.contains(it)) {
        throw new Exception("Propert " + file)
    }
}

return true