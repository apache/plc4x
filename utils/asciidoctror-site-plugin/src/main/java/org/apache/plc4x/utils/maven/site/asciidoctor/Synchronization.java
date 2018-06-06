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

import org.apache.maven.plugins.annotations.Parameter;

public class Synchronization {
    public static final String PREFIX = AsciidoctorMaven.PREFIX + "synchronization.";
    @Parameter(property = PREFIX + "source")
    protected File source;

    @Parameter(property = PREFIX + "target")
    protected File target;

    public Synchronization() {
    }

    public Synchronization(File source, File target) {
        this.source = source;
        this.target = target;
    }

    public File getSource() {
        return source;
    }

    public void setSource(final File source) {
        this.source = source;
    }

    public File getTarget() {
        return target;
    }

    public void setTarget(final File target) {
        this.target = target;
    }
}
