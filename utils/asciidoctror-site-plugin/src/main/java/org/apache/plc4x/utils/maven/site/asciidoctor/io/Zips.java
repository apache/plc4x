/*
 * Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.utils.maven.site.asciidoctor.io;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.zip.ZipEntry;
import java.util.zip.ZipOutputStream;

import org.apache.commons.io.IOUtils;

public final class Zips {
    public static void zip(final File dir, final File zipName) throws IOException, IllegalArgumentException {
        final String[] entries = dir.list();
        final ZipOutputStream out = new ZipOutputStream(new FileOutputStream(zipName));

        String prefix = dir.getAbsolutePath();
        if (!prefix.endsWith(File.separator)) {
            prefix += File.separator;
        }

        for (final String entry : entries) {
            File f = new File(dir, entry);
            zip(out, f, prefix, zipName.getName().substring(0, zipName.getName().length() - 4));
        }
        IOUtils.closeQuietly(out);
    }

    private static void zip(final ZipOutputStream out, final File f, final String prefix, final String root) throws IOException {
        if (f.isDirectory()) {
            final File[] files = f.listFiles();
            if (files != null) {
                for (final File child : files) {
                    zip(out, child, prefix, root);
                }
            }
        } else {
            final FileInputStream in = new FileInputStream(f);
            final ZipEntry entry = new ZipEntry(root + "/" + f.getPath().replace(prefix, ""));
            out.putNextEntry(entry);
            IOUtils.copy(in, out);
            IOUtils.closeQuietly(in);
        }
    }

    private Zips() {
        // no-op
    }
}
