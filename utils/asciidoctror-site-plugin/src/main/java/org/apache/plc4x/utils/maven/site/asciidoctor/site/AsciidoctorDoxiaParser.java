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

import org.apache.maven.doxia.module.xhtml.XhtmlParser;
import org.apache.maven.doxia.parser.ParseException;
import org.apache.maven.doxia.parser.Parser;
import org.apache.maven.doxia.sink.Sink;
import org.apache.maven.doxia.site.decoration.inheritance.PathUtils;
import org.apache.maven.doxia.siterenderer.RenderingContext;
import org.apache.maven.doxia.siterenderer.sink.SiteRendererSink;
import org.apache.maven.project.MavenProject;
import org.asciidoctor.Asciidoctor;
import org.asciidoctor.AttributesBuilder;
import org.asciidoctor.OptionsBuilder;
import org.asciidoctor.SafeMode;
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorHelper;
import org.codehaus.plexus.component.annotations.Component;
import org.codehaus.plexus.component.annotations.Requirement;
import org.codehaus.plexus.util.IOUtil;
import org.codehaus.plexus.util.xml.Xpp3Dom;

import java.io.File;
import java.io.IOException;
import java.io.Reader;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * This class is used by <a href="https://maven.apache.org/doxia/overview.html">the Doxia framework</a>
 * to handle the actual parsing of the AsciiDoc input files into HTML to be consumed/wrapped
 * by the Maven site generation process
 * (see <a href="https://maven.apache.org/plugins/maven-site-plugin/">maven-site-plugin</a>).
 *
 * @author jdlee
 * @author mojavelinux
 */
@Component(role = Parser.class, hint = AsciidoctorDoxiaParser.ROLE_HINT)
public class AsciidoctorDoxiaParser extends XhtmlParser {

    @Requirement
    protected MavenProject project;

    /**
     * The role hint for the {@link AsciidoctorDoxiaParser} Plexus component.
     */
    public static final String ROLE_HINT = "asciidoc";

    protected final Asciidoctor asciidoctor = Asciidoctor.Factory.create();

    public AsciidoctorDoxiaParser() {
        super();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void parse(Reader reader, Sink sink) {
        String source;
        try {
            if ((source = IOUtil.toString(reader)) == null) {
                source = "";
            }
        }
        catch (IOException ex) {
            getLog().error("Could not read AsciiDoc source: " + ex.getLocalizedMessage());
            return;
        }

        File moduleBaseDirectory = null;
        if(sink instanceof SiteRendererSink) {
            SiteRendererSink siteRendererSink = (SiteRendererSink) sink;
            RenderingContext renderingContext = siteRendererSink.getRenderingContext();
            String asciidocBaseDir = renderingContext.getBasedir().getAbsolutePath();
            String asciidocRelativeDir = renderingContext.getBasedirRelativePath();
            moduleBaseDirectory = new File(asciidocBaseDir.substring(0, asciidocBaseDir.length() - asciidocRelativeDir.length()));
        }

        Xpp3Dom siteConfig = getSiteConfig(project);

        File siteDirectory = resolveSiteDirectory(project, siteConfig);

        OptionsBuilder options = processAsciiDocConfig(
                siteConfig,
                initOptions(project, siteDirectory),
                initAttributes(project, siteDirectory),
            moduleBaseDirectory);
        // QUESTION should we keep OptionsBuilder & AttributesBuilder separate for call to convertAsciiDoc?
        sink.rawText(convertAsciiDoc(source, options));
    }

    protected Xpp3Dom getSiteConfig(MavenProject project) {
        return project.getGoalConfiguration("org.apache.maven.plugins", "maven-site-plugin", "site", "site");
    }

    protected File resolveSiteDirectory(MavenProject project, Xpp3Dom siteConfig) {
        File siteDirectory = new File(project.getBasedir(), "src/site");
        if (siteConfig != null) {
            Xpp3Dom siteDirectoryNode = siteConfig.getChild("siteDirectory");
            if (siteDirectoryNode != null) {
                siteDirectory = new File(siteDirectoryNode.getValue());
            }
        }
        return siteDirectory;
    }

    protected OptionsBuilder initOptions(MavenProject project, File siteDirectory) {
        return OptionsBuilder.options()
                .backend("xhtml")
                .safe(SafeMode.UNSAFE)
                .baseDir(new File(siteDirectory, ROLE_HINT));
    }

    protected AttributesBuilder initAttributes(MavenProject project, File siteDirectory) {
        return AttributesBuilder.attributes()
            .attribute("idprefix", "@")
            .attribute("showtitle", "@");
    }

    protected OptionsBuilder processAsciiDocConfig(Xpp3Dom siteConfig, OptionsBuilder options, AttributesBuilder attributes,
                                                   File moduleBaseDirectory) {
        if (siteConfig == null) {
            return options.attributes(attributes);
        }

        Xpp3Dom asciidocConfig = siteConfig.getChild("asciidoc");
        if (asciidocConfig == null) {
            return options.attributes(attributes);
        }

        if (this.project.getProperties() != null) {
            for ( Map.Entry<Object, Object> entry : this.project.getProperties().entrySet() ) {
                attributes.attribute(((String) entry.getKey()).replaceAll("\\.", "-"), entry.getValue());
            }
        }

        for (Xpp3Dom asciidocOpt : asciidocConfig.getChildren()) {
            String optName = asciidocOpt.getName();
            if ("attributes".equals(optName)) {
                for (Xpp3Dom asciidocAttr : asciidocOpt.getChildren()) {
                    String attrValue = asciidocAttr.getValue();
                    if((attrValue != null) && attrValue.contains("@{project.basedir}")) {
                        attrValue = attrValue.replace("@{project.basedir}", moduleBaseDirectory.getAbsolutePath());
                    }
                    AsciidoctorHelper.addAttribute(asciidocAttr.getName(), attrValue, attributes);
                }
            }
            else if ("requires".equals(optName)) {
                Xpp3Dom[] requires = asciidocOpt.getChildren("require");
                // supports variant:
                // <requires>
                //     <require>time</require>
                // </requires>
                if (requires.length > 0) {
                    for (Xpp3Dom require : requires) {
                        requireLibrary(require.getValue());
                    }
                }
                else {
                    // supports variant:
                    // <requires>time, base64</requires>
                    for (String require : asciidocOpt.getValue().split(",")) {
                        requireLibrary(require);
                    }
                }
            }
            else if ("templateDir".equals(optName) || "template_dir".equals(optName)) {
                options.templateDir(resolveProjectDir(project, asciidocOpt.getValue()));
            }
            else if ("templateDirs".equals(optName) || "template_dirs".equals(optName)) {
                List<File> templateDirs = new ArrayList<File>();
                for (Xpp3Dom dir : asciidocOpt.getChildren("dir")) {
                    templateDirs.add(resolveProjectDir(project, dir.getValue()));
                }
                if (!templateDirs.isEmpty()) {
                    options.templateDirs(templateDirs.toArray(new File[0]));
                }
            }
            else if ("baseDir".equals(optName)) {
                options.baseDir(resolveProjectDir(project, asciidocOpt.getValue()));
            }
            else {
                options.option(optName.replaceAll("(?<!_)([A-Z])", "_$1").toLowerCase(), asciidocOpt.getValue());
            }
        }
        return options.attributes(attributes);
    }

    protected String convertAsciiDoc(String source, OptionsBuilder options) {
        return asciidoctor.convert(source, options);
    }

    protected File resolveProjectDir(MavenProject project, String path) {
        File filePath = new File(path);
        if (!filePath.isAbsolute()) {
            filePath = new File(project.getBasedir(), filePath.toString());
        }
        return filePath;
    }

    private void requireLibrary(String require) {
        if (!(require = require.trim()).isEmpty()) {
            try {
                asciidoctor.requireLibrary(require);
            } catch (Exception ex) {
                getLog().error(ex.getLocalizedMessage());
            }
        }
    }
}
