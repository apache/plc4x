package org.apache.plc4x.utils.maven.site.asciidoctor.test.io

import java.util.concurrent.CountDownLatch

class PrefilledInputStream extends ByteArrayInputStream {
    final CountDownLatch latch

    PrefilledInputStream(final byte[] buf, final CountDownLatch latch) {
        super(buf)
        this.latch = latch
    }

    @Override
    synchronized int read(final byte[] b, final int off, final int len) {
        latch.await()
        return super.read(b, off, len)
    }
}
