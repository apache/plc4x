/*
 * Copyright 2015 The Ascidoctor Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.utils.maven.site.asciidoctor.site;

import org.apache.maven.doxia.parser.module.AbstractParserModule;
import org.apache.maven.doxia.parser.module.ParserModule;
import org.codehaus.plexus.component.annotations.Component;

/**
 * This class is the entry point for integration with the Maven Site Plugin
 * integration since Doxia 1.6 (i.e., maven-site-plugin 3.4 and above):
 * it defines source directory and file extensions to be added to
 * <a href="https://maven.apache.org/doxia/references/">Doxia provided modules</a>.
 *
 * @author jdlee
 */
@Component(role = ParserModule.class, hint = AsciidoctorDoxiaParser.ROLE_HINT)
public class AsciidoctorDoxiaParserModule extends AbstractParserModule {

    /**
     * The source directory for AsciiDoc files.
     */
    public static final String SOURCE_DIRECTORY = AsciidoctorDoxiaParser.ROLE_HINT;

    /**
     * The extension for AsciiDoc files.
     */
    // TODO change type to String[] and value to { "adoc", "asciidoc" } once available in Doxia
    public static final String FILE_EXTENSION = "adoc";

    /**
     * Build a new instance of {@link AsciidoctorDoxiaParserModule}.
     */
    public AsciidoctorDoxiaParserModule() {
        super(SOURCE_DIRECTORY, FILE_EXTENSION, AsciidoctorDoxiaParser.ROLE_HINT);
    }
}
