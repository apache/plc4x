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

package org.apache.plc4x.utils.maven.site.asciidoctor;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.util.Map;

import org.apache.commons.io.IOUtils;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugin.MojoFailureException;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.asciidoctor.Asciidoctor;
import org.apache.plc4x.utils.maven.site.asciidoctor.http.AsciidoctorHttpServer;
import org.apache.plc4x.utils.maven.site.asciidoctor.io.IO;

@Mojo(name = "http")
public class AsciidoctorHttpMojo extends AsciidoctorRefreshMojo {
    public static final String PREFIX = AsciidoctorMaven.PREFIX + "http.";

    @Parameter(property = PREFIX + "home", required = false, defaultValue = "index")
    protected String home;

    @Parameter(property = PREFIX + "reload-interval", required = false, defaultValue = "0")
    protected int autoReloadInterval;

    @Override
    protected void doWork() throws MojoFailureException, MojoExecutionException {
        final AsciidoctorHttpServer server = new AsciidoctorHttpServer(getLog(), port, outputDirectory, home);
        server.start();

        super.doWork();

        server.stop();
    }

    @Override
    protected void renderFile(final Asciidoctor asciidoctorInstance, final Map<String, Object> options, final File f) {
        asciidoctorInstance.renderFile(f, options);

        if (autoReloadInterval > 0 && backend.toLowerCase().startsWith("html")) {
            final String filename = f.getName();
            final File out = new File(outputDirectory, filename.substring(0, filename.lastIndexOf(".")) + ".html");
            if (out.exists()) {

                String content = null;

                { // read
                    FileInputStream fis = null;
                    try {
                        fis = new FileInputStream(out); // java asciidoctor render() doesn't work ATM so read the rendered file instead of doing it in memory
                        content = IO.slurp(fis);
                    } catch (final Exception e) {
                        getLog().error(e);
                    } finally {
                        IOUtils.closeQuietly(fis);
                    }
                }

                if (content != null) { // render + write
                    FileOutputStream fos = null;
                    try {
                        fos = new FileOutputStream(out);
                        fos.write(addRefreshing(content).getBytes());
                    } catch (final Exception e) {
                        getLog().error(e);
                    } finally {
                        IOUtils.closeQuietly(fos);
                    }
                }
            }
        } else {
            asciidoctorInstance.renderFile(f, options);
        }

        logRenderedFile(f);
    }

    private String addRefreshing(final String html) {
        return html.replace("</body>", "<script>setTimeout(\"location.reload(true);\"," + autoReloadInterval + ");</script>\n</body>");
    }

    public String getHome() {
        return home;
    }

    public void setHome(final String home) {
        this.home = home;
    }
}
