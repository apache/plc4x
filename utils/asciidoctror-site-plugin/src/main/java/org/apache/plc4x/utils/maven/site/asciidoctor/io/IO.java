package org.apache.plc4x.utils.maven.site.asciidoctor.io;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

public final class IO {
    public static String slurp(final InputStream from) throws IOException {
        final ByteArrayOutputStream to = new ByteArrayOutputStream();
        final byte[] buffer = new byte[1024];
        int length;
        while ((length = from.read(buffer)) != -1) {
            to.write(buffer, 0, length);
        }
        to.flush();
        return new String(to.toByteArray());
    }

    private IO() {
        // no-op
    }
}
