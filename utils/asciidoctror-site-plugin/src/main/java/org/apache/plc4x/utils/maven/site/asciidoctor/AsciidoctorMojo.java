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

import org.apache.commons.io.FileUtils;
import org.apache.maven.execution.MavenSession;
import org.apache.maven.model.Resource;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugin.MojoFailureException;
import org.apache.maven.plugins.annotations.Component;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;
import org.apache.maven.shared.filtering.MavenFilteringException;
import org.apache.maven.shared.filtering.MavenResourcesExecution;
import org.apache.maven.shared.filtering.MavenResourcesFiltering;
import org.asciidoctor.*;
import org.asciidoctor.internal.JRubyRuntimeContext;
import org.apache.plc4x.utils.maven.site.asciidoctor.extensions.AsciidoctorJExtensionRegistry;
import org.apache.plc4x.utils.maven.site.asciidoctor.extensions.ExtensionConfiguration;
import org.apache.plc4x.utils.maven.site.asciidoctor.extensions.ExtensionRegistry;
import org.apache.plc4x.utils.maven.site.asciidoctor.io.AsciidoctorFileScanner;
import org.jruby.Ruby;
import org.sonatype.plexus.build.incremental.BuildContext;

import java.io.File;
import java.io.IOException;
import java.util.*;


/**
 * Basic maven plugin goal to render AsciiDoc files using Asciidoctor, a ruby port.
 */
@Mojo(name = "process-asciidoc")
public class AsciidoctorMojo extends AbstractMojo {
    // copied from org.asciidoctor.AsciiDocDirectoryWalker.ASCIIDOC_REG_EXP_EXTENSION
    // should probably be configured in AsciidoctorMojo through @Parameter 'extension'
    protected static final String ASCIIDOC_REG_EXP_EXTENSION = ".*\\.a((sc(iidoc)?)|d(oc)?)$";

    @Parameter(defaultValue = "${project.build.sourceEncoding}")
    protected String encoding;

    @Parameter(property = AsciidoctorMaven.PREFIX + "sourceDirectory", defaultValue = "${basedir}/src/main/asciidoc", required = true)
    protected File sourceDirectory;

    @Parameter(property = AsciidoctorMaven.PREFIX + "outputDirectory", defaultValue = "${project.build.directory}/generated-docs", required = true)
    protected File outputDirectory;

    @Parameter(property = AsciidoctorMaven.PREFIX + "outputFile", required = false)
    protected File outputFile;

    @Parameter(property = AsciidoctorMaven.PREFIX + "preserveDirectories", defaultValue = "false", required = false)
    protected boolean preserveDirectories = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "relativeBaseDir", defaultValue = "false", required = false)
    protected boolean relativeBaseDir = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "projectDirectory", defaultValue = "${basedir}", required = false, readonly = false)
    protected File projectDirectory;

    @Parameter(property = AsciidoctorMaven.PREFIX + "rootDir", defaultValue = "${basedir}", required = false, readonly = false)
    protected File rootDir;

    @Parameter(property = AsciidoctorMaven.PREFIX + "baseDir", required = false)
    protected File baseDir;

    @Parameter(property = AsciidoctorMaven.PREFIX + "skip", required = false)
    protected boolean skip = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "gemPath", defaultValue = "", required = false)
    protected String gemPath = "";

    @Parameter(property = AsciidoctorMaven.PREFIX + "requires")
    protected List<String> requires = new ArrayList<String>();

    @Parameter(required = false)
    protected Map<String, Object> attributes = new HashMap<String, Object>();

    @Parameter(property = AsciidoctorMaven.PREFIX + Options.ATTRIBUTES, required = false)
    protected String attributesChain = "";

    @Parameter(property = AsciidoctorMaven.PREFIX + Options.BACKEND, defaultValue = "docbook", required = true)
    protected String backend = "";

    @Parameter(property = AsciidoctorMaven.PREFIX + Options.DOCTYPE, required = false)
    protected String doctype;

    @Parameter(property = AsciidoctorMaven.PREFIX + Options.ERUBY, required = false)
    protected String eruby = "";

    @Parameter(property = AsciidoctorMaven.PREFIX + "headerFooter", required = false)
    protected boolean headerFooter = true;

    @Parameter(property = AsciidoctorMaven.PREFIX + "templateDir", required = false)
    protected File templateDir;

    @Parameter(property = AsciidoctorMaven.PREFIX + "templateEngine", required = false)
    protected String templateEngine;

    @Parameter(property = AsciidoctorMaven.PREFIX + "templateCache")
    protected boolean templateCache = true;

    @Parameter(property = AsciidoctorMaven.PREFIX + "imagesDir", required = false)
    protected String imagesDir = "images@"; // '@' Allows override by :imagesdir: document attribute

    @Parameter(property = AsciidoctorMaven.PREFIX + "sourceHighlighter", required = false)
    protected String sourceHighlighter;

    @Parameter(property = AsciidoctorMaven.PREFIX + Attributes.TITLE, required = false)
    protected String title = "";

    @Parameter(property = AsciidoctorMaven.PREFIX + "sourceDocumentName", required = false)
    protected String sourceDocumentName;

    @Parameter(property = AsciidoctorMaven.PREFIX + "sourceDocumentExtensions")
    protected List<String> sourceDocumentExtensions = new ArrayList<String>();

    @Parameter(property = AsciidoctorMaven.PREFIX + "sourcemap")
    protected boolean sourcemap = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "catalogAssets")
    protected boolean catalogAssets = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "synchronizations", required = false)
    protected List<Synchronization> synchronizations = new ArrayList<Synchronization>();

    @Parameter(property = AsciidoctorMaven.PREFIX + "extensions")
    protected List<ExtensionConfiguration> extensions = new ArrayList<ExtensionConfiguration>();

    @Parameter(property = AsciidoctorMaven.PREFIX + "embedAssets")
    protected boolean embedAssets = false;

    @Parameter(property = AsciidoctorMaven.PREFIX + "attributeMissing")
    protected String attributeMissing = "skip";

    @Parameter(property = AsciidoctorMaven.PREFIX + "attributeUndefined")
    protected String attributeUndefined = "drop-line";

    // List of resources to copy to the output directory (e.g., images, css). By default everything is copied
    @Parameter(property = AsciidoctorMaven.PREFIX + "sources")
    protected List<Resource> resources;

    @Parameter(defaultValue = "${project}", readonly = true, required = true)
    protected MavenProject project;

    @Parameter(defaultValue = "${session}", readonly = true, required = true)
    protected MavenSession session;

    @Component
    protected MavenResourcesFiltering outputResourcesFiltering;

    @Component
    protected BuildContext buildContext;

    @Override
    public void execute() throws MojoExecutionException, MojoFailureException {
        if (skip) {
            getLog().info("AsciiDoc processing is skipped.");
            return;
        }

        if (sourceDirectory == null) {
            throw new MojoExecutionException("Required parameter 'asciidoctor.sourceDirectory' not set.");
        }

        if (!sourceDirectory.exists()) {
            getLog().info("sourceDirectory " + sourceDirectory.getPath() + " does not exist. Skip processing");
            return;
        }

        ensureOutputExists();

        // Validate resources to avoid errors later on
        if (resources != null) {
            for (Resource resource : resources) {
                if (resource.getDirectory() == null || resource.getDirectory().isEmpty()) {
                    throw new MojoExecutionException("Found empty resource directory");
                }
            }
        }

        final Asciidoctor asciidoctor = getAsciidoctorInstance(gemPath);

        asciidoctor.requireLibraries(requires);

        final OptionsBuilder optionsBuilder = OptionsBuilder.options();
        setOptionsOnBuilder(optionsBuilder);

        final AttributesBuilder attributesBuilder = AttributesBuilder.attributes();
        setAttributesOnBuilder(attributesBuilder);

        optionsBuilder.attributes(attributesBuilder);

        ExtensionRegistry extensionRegistry = new AsciidoctorJExtensionRegistry(asciidoctor);
        for (ExtensionConfiguration extension : extensions) {
            try {
                extensionRegistry.register(extension.getClassName(), extension.getBlockName());
            } catch (Exception e) {
                throw new MojoExecutionException(e.getMessage(), e);
            }
        }

        // Copy output resources
        prepareResources();
        copyResources();

        // Prepare sources
        final List<File> sourceFiles = sourceDocumentName == null ?
                scanSourceFiles() : Arrays.asList(new File(sourceDirectory, sourceDocumentName));

        final Set<File> dirs = new HashSet<File>();
        for (final File source : sourceFiles) {
            final File destinationPath = setDestinationPaths(optionsBuilder, source);
            if (!dirs.add(destinationPath))
                getLog().warn("Duplicated destination found: overwriting file: " + destinationPath.getAbsolutePath());

            renderFile(asciidoctor, optionsBuilder.asMap(), source);
        }

        if (synchronizations != null && !synchronizations.isEmpty()) {
            synchronize();
        }
    }

    /**
     * Initializes resources attribute excluding AsciiDoc documents and hidden directories/files (those prefixed with
     * underscore).
     * By default everything in the sources directories is copied.
     */
    private void prepareResources() {
        if (resources == null || resources.isEmpty()) {
            resources = new ArrayList<Resource>();
            // we don't want to copy files considered sources
            Resource resource = new Resource();
            resource.setDirectory(sourceDirectory.getAbsolutePath());
            resource.setExcludes(new ArrayList<String>());
            // exclude sourceDocumentName if defined
            if (sourceDocumentName != null && sourceDocumentName.isEmpty()) {
                resource.getExcludes().add(sourceDocumentName);
            }
            // exclude filename extensions if defined
            resources.add(resource);
        }

        // All resources must exclude AsciiDoc documents and folders beginning with underscore
        for (Resource resource : resources) {
            if (resource.getExcludes() == null || resource.getExcludes().isEmpty()) {
                resource.setExcludes(new ArrayList<String>());
            }
            List<String> excludes = new ArrayList<String>();
            for (String value : AsciidoctorFileScanner.IGNORED_FOLDERS_AND_FILES) {
                excludes.add(value);
            }
            for (String value : AsciidoctorFileScanner.DEFAULT_FILE_EXTENSIONS) {
                excludes.add(value);
            }
            for (String docExtension : sourceDocumentExtensions) {
                resource.getExcludes().add("**/*." + docExtension);
            }
            // in case someone wants to include some of the default excluded files (.e.g., AsciiDoc docs)
            excludes.removeAll(resource.getIncludes());
            resource.getExcludes().addAll(excludes);
        }
    }

    /**
     * Copies the resources defined in the 'resources' attribute
     */
    private void copyResources() throws MojoExecutionException {
        try {
            // Right now it's not used at all, but could be used to apply resource filters/replacements
            MavenResourcesExecution resourcesExecution =
                    new MavenResourcesExecution(resources, outputDirectory, project, encoding,
                            Collections.<String>emptyList(), Collections.<String>emptyList(), session);
            resourcesExecution.setIncludeEmptyDirs(true);
            resourcesExecution.setAddDefaultExcludes(true);
            outputResourcesFiltering.filterResources(resourcesExecution);
        } catch (MavenFilteringException e) {
            throw new MojoExecutionException("Could not copy resources", e);
        }
    }

    /**
     * Updates optionsBuilder object's baseDir and destination(s) accordingly to the options.
     *
     * @param optionsBuilder AsciidoctorJ options to be updated.
     * @param sourceFile     AsciiDoc source file to process.
     * @return the final destination file path.
     */
    private File setDestinationPaths(OptionsBuilder optionsBuilder, final File sourceFile) throws MojoExecutionException {
        try {
            if (baseDir != null) {
                optionsBuilder.baseDir(baseDir);
            } else {
                // when preserveDirectories == false, parent and sourceDirectory are the same
                if (relativeBaseDir) {
                    optionsBuilder.baseDir(sourceFile.getParentFile());
                } else {
                    optionsBuilder.baseDir(sourceDirectory);
                }
            }
            if (preserveDirectories) {
                final String candidatePath = sourceFile.getParentFile().getCanonicalPath().substring(sourceDirectory.getCanonicalPath().length());
                final File relativePath = new File(outputDirectory.getCanonicalPath() + candidatePath);
                optionsBuilder.toDir(relativePath).destinationDir(relativePath);
            } else {
                optionsBuilder.toDir(outputDirectory).destinationDir(outputDirectory);
            }
            if (outputFile != null) {
                //allow overriding the output file name
                optionsBuilder.toFile(outputFile);
            }
            // return destination file path
            if (outputFile != null) {
                return outputFile.isAbsolute() ?
                        outputFile : new File((String) optionsBuilder.asMap().get(Options.DESTINATION_DIR), outputFile.getPath());
            } else
                return new File((String) optionsBuilder.asMap().get(Options.DESTINATION_DIR), sourceFile.getName());
        } catch (IOException e) {
            throw new MojoExecutionException("Unable to locate output directory", e);
        }
    }

    protected Asciidoctor getAsciidoctorInstance(String gemPath) throws MojoExecutionException {
        Asciidoctor asciidoctor = null;
        if (gemPath == null) {
            asciidoctor = Asciidoctor.Factory.create();
        } else {
            // Replace Windows path separator to avoid paths with mixed \ and /.
            // This happens for instance when setting: <gemPath>${project.build.directory}/gems-provided</gemPath>
            // because the project's path is converted to string.
            String normalizedGemPath = (File.separatorChar == '\\') ? gemPath.replaceAll("\\\\", "/") : gemPath;
            asciidoctor = Asciidoctor.Factory.create(normalizedGemPath);
        }

        Ruby rubyInstance = null;
        try {
            rubyInstance = (Ruby) JRubyRuntimeContext.class.getMethod("get")
                    .invoke(null);
        } catch (NoSuchMethodException e) {
            if (rubyInstance == null) {
                try {
                    rubyInstance = (Ruby) JRubyRuntimeContext.class.getMethod(
                            "get", Asciidoctor.class).invoke(null, asciidoctor);
                } catch (Exception e1) {
                    throw new MojoExecutionException(
                            "Failed to invoke get(AsciiDoctor) for JRubyRuntimeContext",
                            e1);
                }

            }
        } catch (Exception e) {
            throw new MojoExecutionException(
                    "Failed to invoke get for JRubyRuntimeContext", e);
        }

        String gemHome = rubyInstance.evalScriptlet("ENV['GEM_HOME']").toString();
        String gemHomeExpected = (gemPath == null || "".equals(gemPath)) ? "" : gemPath.split(File.pathSeparator)[0];

        if (!"".equals(gemHome) && !gemHomeExpected.equals(gemHome)) {
            getLog().warn("Using inherited external environment to resolve gems (" + gemHome + "), i.e. build is platform dependent!");
        }

        return asciidoctor;
    }

    private List<File> scanSourceFiles() {
        final List<File> asciidoctorFiles;
        if (sourceDocumentExtensions == null || sourceDocumentExtensions.isEmpty()) {
            final DirectoryWalker directoryWalker = new AsciiDocDirectoryWalker(sourceDirectory.getAbsolutePath());
            asciidoctorFiles = directoryWalker.scan();
        } else {
            final DirectoryWalker directoryWalker = new CustomExtensionDirectoryWalker(sourceDirectory.getAbsolutePath(), sourceDocumentExtensions);
            asciidoctorFiles = directoryWalker.scan();
        }
        String absoluteSourceDirectory = sourceDirectory.getAbsolutePath();
        for (Iterator<File> iter = asciidoctorFiles.iterator(); iter.hasNext(); ) {
            File f = iter.next();
            do {
                // stop when we hit the source directory root
                if (absoluteSourceDirectory.equals(f.getAbsolutePath())) {
                    break;
                }
                // skip if the filename or directory begins with _
                if (f.getName().startsWith("_")) {
                    iter.remove();
                    break;
                }
            } while ((f = f.getParentFile()) != null);
        }
        return asciidoctorFiles;
    }

    private void synchronize() {
        for (final Synchronization synchronization : synchronizations) {
            synchronize(synchronization);
        }
    }

    protected void renderFile(Asciidoctor asciidoctor, Map<String, Object> options, File f) {
        asciidoctor.renderFile(f, options);
        logRenderedFile(f);
    }

    protected void logRenderedFile(File f) {
        getLog().info("Rendered " + f.getAbsolutePath());
    }

    protected void synchronize(final Synchronization synchronization) {
        if (synchronization.getSource().isDirectory()) {
            try {
                FileUtils.copyDirectory(synchronization.getSource(), synchronization.getTarget());
            } catch (IOException e) {
                getLog().error(String.format("Can't synchronize %s -> %s", synchronization.getSource(), synchronization.getTarget()));
            }
        } else {
            try {
                FileUtils.copyFile(synchronization.getSource(), synchronization.getTarget());
            } catch (IOException e) {
                getLog().error(String.format("Can't synchronize %s -> %s", synchronization.getSource(), synchronization.getTarget()));
            }
        }
    }

    protected void ensureOutputExists() {
        if (!outputDirectory.exists()) {
            if (!outputDirectory.mkdirs()) {
                getLog().error("Can't create " + outputDirectory.getPath());
            }
        }
    }

    /**
     * Updates and OptionsBuilder instance with the options defined in the configuration.
     *
     * @param optionsBuilder AsciidoctorJ options to be updated.
     */
    protected void setOptionsOnBuilder(OptionsBuilder optionsBuilder) {
        optionsBuilder
                .backend(backend)
                .safe(SafeMode.UNSAFE)
                .headerFooter(headerFooter)
                .eruby(eruby)
                .mkDirs(true);

        // Following options are only set when the value is different than the default
        if (sourcemap)
            optionsBuilder.option("sourcemap", true);

        if (catalogAssets)
            optionsBuilder.option("catalog_assets", true);

        if (!templateCache)
            optionsBuilder.option("template_cache", false);

        if (doctype != null)
            optionsBuilder.docType(doctype);

        if (templateEngine != null)
            optionsBuilder.templateEngine(templateEngine);

        if (templateDir != null)
            optionsBuilder.templateDir(templateDir);
    }

    protected void setAttributesOnBuilder(AttributesBuilder attributesBuilder) throws MojoExecutionException {
        if (sourceHighlighter != null) {
            attributesBuilder.sourceHighlighter(sourceHighlighter);
        }

        if (embedAssets) {
            attributesBuilder.linkCss(false);
            attributesBuilder.dataUri(true);
        }

        if (imagesDir != null) {
            attributesBuilder.imagesDir(imagesDir);
        }

        if ("skip".equals(attributeMissing) || "drop".equals(attributeMissing) || "drop-line".equals(attributeMissing)) {
            attributesBuilder.attributeMissing(attributeMissing);
        } else {
            throw new MojoExecutionException(attributeMissing + " is not valid. Must be one of 'skip', 'drop' or 'drop-line'");
        }

        if ("drop".equals(attributeUndefined) || "drop-line".equals(attributeUndefined)) {
            attributesBuilder.attributeUndefined(attributeUndefined);
        } else {
            throw new MojoExecutionException(attributeUndefined + " is not valid. Must be one of 'drop' or 'drop-line'");
        }

        AsciidoctorHelper.addAttributes(attributes, attributesBuilder);

        if (!attributesChain.isEmpty()) {
            getLog().info("Attributes: " + attributesChain);
            attributesBuilder.arguments(attributesChain);
        }

    }

    public File getSourceDirectory() {
        return sourceDirectory;
    }

    public void setSourceDirectory(File sourceDirectory) {
        this.sourceDirectory = sourceDirectory;
    }

    public File getOutputDirectory() {
        return outputDirectory;
    }

    public void setOutputDirectory(File outputDirectory) {
        this.outputDirectory = outputDirectory;
    }

    public File getOutputFile() {
        return outputFile;
    }

    public void setOutputFile(File outputFile) {
        this.outputFile = outputFile;
    }

    public String getBackend() {
        return backend;
    }

    public void setBackend(String backend) {
        this.backend = backend;
    }

    public boolean isSkip() {
        return skip;
    }

    public void setSkip(boolean skip) {
        this.skip = skip;
    }

    public String getDoctype() {
        return doctype;
    }

    public void setDoctype(String doctype) {
        this.doctype = doctype;
    }

    public Map<String, Object> getAttributes() {
        return attributes;
    }

    public void setAttributes(Map<String, Object> attributes) {
        this.attributes = attributes;
    }

    public boolean isHeaderFooter() {
        return headerFooter;
    }

    public void setHeaderFooter(boolean headerFooter) {
        this.headerFooter = headerFooter;
    }

    public File getTemplateDir() {
        return templateDir;
    }

    public void setTemplateDir(File templateDir) {
        this.templateDir = templateDir;
    }

    public String getTemplateEngine() {
        return templateEngine;
    }

    public void setTemplateEngine(String templateEngine) {
        this.templateEngine = templateEngine;
    }

    public String getImagesDir() {
        return imagesDir;
    }

    public void setImagesDir(String imagesDir) {
        this.imagesDir = imagesDir;
    }

    public String getSourceHighlighter() {
        return sourceHighlighter;
    }

    public void setSourceHighlighter(String sourceHighlighter) {
        this.sourceHighlighter = sourceHighlighter;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public List<String> getSourceDocumentExtensions() {
        return sourceDocumentExtensions;
    }

    public void setSourceDocumentExtensions(final List<String> sourceDocumentExtensions) {
        this.sourceDocumentExtensions = sourceDocumentExtensions;
    }

    public String getEruby() {
        return eruby;
    }

    public void setEruby(String eruby) {
        this.eruby = eruby;
    }

    public String getSourceDocumentName() {
        return sourceDocumentName;
    }

    public void setSourceDocumentName(String sourceDocumentName) {
        this.sourceDocumentName = sourceDocumentName;
    }

    public List<Synchronization> getSynchronizations() {
        return synchronizations;
    }

    public void setSynchronizations(List<Synchronization> synchronizations) {
        this.synchronizations = synchronizations;
    }

    public boolean isEmbedAssets() {
        return embedAssets;
    }

    public void setEmbedAssets(boolean embedAssets) {
        this.embedAssets = embedAssets;
    }

    public String getAttributeMissing() {
        return attributeMissing;
    }

    public void setAttributeMissing(String attributeMissing) {
        this.attributeMissing = attributeMissing;
    }

    public String getAttributeUndefined() {
        return attributeUndefined;
    }

    public void setAttributeUndefined(String attributeUndefined) {
        this.attributeUndefined = attributeUndefined;
    }

    public File getProjectDirectory() {
        return projectDirectory;
    }

    public void setProjectDirectory(File projectDirectory) {
        this.projectDirectory = projectDirectory;
    }

    public File getRootDir() {
        return rootDir;
    }

    public void setRootDir(File rootDir) {
        this.rootDir = rootDir;
    }

    public String getGemPath() {
        return gemPath;
    }

    public void setGemPath(String gemPath) {
        this.gemPath = gemPath;
    }

    public File getBaseDir() {
        return baseDir;
    }

    public void setBaseDir(File baseDir) {
        this.baseDir = baseDir;
    }

    public void setPreserveDirectories(boolean preserveDirectories) {
        this.preserveDirectories = preserveDirectories;
    }

    public void setRelativeBaseDir(boolean relativeBaseDir) {
        this.relativeBaseDir = relativeBaseDir;
    }

    public List<ExtensionConfiguration> getExtensions() {
        return extensions;
    }

    public void setExtensions(List<ExtensionConfiguration> extensions) {
        this.extensions = extensions;
    }

    public List<Resource> getResources() {
        return resources;
    }

    public void setResources(List<Resource> resources) {
        this.resources = resources;
    }

    private static class CustomExtensionDirectoryWalker extends AbstractDirectoryWalker {
        private final List<String> fileExtensions;

        public CustomExtensionDirectoryWalker(final String absolutePath, final List<String> fileExtensions) {
            super(absolutePath);
            this.fileExtensions = fileExtensions;
        }

        @Override
        protected boolean isAcceptedFile(final File filename) {
            final String name = filename.getName();
            for (final String extension : fileExtensions) {
                if (name.endsWith(extension)) {
                    return true;
                }
            }
            return false;
        }
    }

}
