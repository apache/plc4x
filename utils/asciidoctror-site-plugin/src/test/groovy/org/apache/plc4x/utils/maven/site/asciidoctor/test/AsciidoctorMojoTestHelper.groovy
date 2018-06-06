package org.apache.plc4x.utils.maven.site.asciidoctor.test

class AsciidoctorMojoTestHelper {
    def getAvailablePort() {
        ServerSocket socket = new ServerSocket(0)
        def port = socket.getLocalPort()
        socket.close()
        return port
    }
}
