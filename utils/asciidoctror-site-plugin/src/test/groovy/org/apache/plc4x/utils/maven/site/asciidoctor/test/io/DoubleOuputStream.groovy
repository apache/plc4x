package org.apache.plc4x.utils.maven.site.asciidoctor.test.io

class DoubleOuputStream extends ByteArrayOutputStream {
    final OutputStream other

    DoubleOuputStream(final OutputStream os) {
        other = os
    }

    @Override
    synchronized void write(final byte[] b, final int off, final int len) {
        other.write(b, off, len)
        super.write(b, off, len)
    }
}
