package org.apache.plc4x.utils.maven.site.asciidoctor.io;

import org.apache.maven.model.Resource;
import org.codehaus.plexus.util.Scanner;
import org.sonatype.plexus.build.incremental.BuildContext;

import java.io.File;
import java.util.*;

/**
 * Recursively traverses directories returning the list of AsciiDoc files that match the applied filters.
 * If no filters are set, AsciiDoc documents with extensions .adoc, .ad, .asc and .asciidoc are returned.
 */
public class AsciidoctorFileScanner {

    public static String[] DEFAULT_FILE_EXTENSIONS = {"**/*.adoc", "**/*.ad", "**/*.asc","**/*.asciidoc"};

    public static String[] IGNORED_FOLDERS_AND_FILES = {
            // Files and directories beginning with underscore are ignored
            "**/_*.*",
            "**/_*",
            "**/_*/**/*.*",
            // docinfo snippets should not be copied
            "**/docinfo.html",
            "**/docinfo-header.html",
            "**/docinfo-footer.html",
            "**/*-docinfo.html",
            "**/*-docinfo-header.html",
            "**/*-docinfo-footer.html",
            "**/docinfo.xml",
            "**/docinfo-header.xml",
            "**/docinfo-footer.xml",
            "**/*-docinfo.xml",
            "**/*-docinfo-header.xml",
            "**/*-docinfo-footer.xml"};

    private BuildContext buildContext;

    public AsciidoctorFileScanner(BuildContext buildContext) {
        this.buildContext = buildContext;
    }

    /**
     * Scans a resource directory (and sub-subdirectories) returning all AsciiDoc documents found.
     *
     * @param resource {@link Resource} to scan (the directory property is mandatory)
     * @return List of found documents matching the resource properties
     */
    public List<File> scan(Resource resource) {
        Scanner scanner = buildContext.newScanner(new File(resource.getDirectory()), true);
        setupScanner(scanner, resource);
        scanner.scan();
        List<File> files = new ArrayList<File>();
        for (String file : scanner.getIncludedFiles()) {
            files.add(new File(resource.getDirectory(), file));
        }
        return files;
    }

    /**
     * Scans a list of resources returning all AsciiDoc documents found.
     *
     * @param resources List of {@link Resource} to scan (the directory property is mandatory)
     * @return List of found documents matching the resources properties
     */
    public List<File> scan(List<Resource> resources) {
        List<File> files = new ArrayList<File>();
        for (Resource resource: resources) {
            files.addAll(scan(resource));
        }
        return files;
    }

    /**
     * Initializes the Scanner with the default values.
     * <br>
     * By default:
     * <ul>
     *     <li>includes adds extension .adoc, .ad, .asc and .asciidoc
     *     <li>excludes adds filters to avoid hidden files and directoris beginning with undersore
     * </ul>
     *
     * NOTE: Patterns both in inclusions and exclusions are automatically excluded.
     */
    private void setupScanner(Scanner scanner, Resource resource) {

        if (resource.getIncludes() == null || resource.getIncludes().isEmpty()) {
            scanner.setIncludes(DEFAULT_FILE_EXTENSIONS);
        } else {
            scanner.setIncludes(resource.getIncludes().toArray(new String[] {}));
        }

        if (resource.getExcludes() == null || resource.getExcludes().isEmpty()) {
            scanner.setExcludes(IGNORED_FOLDERS_AND_FILES);
        } else {
            scanner.setExcludes(mergeAndConvert(resource.getExcludes(), IGNORED_FOLDERS_AND_FILES));
        }
        // adds exclusions like SVN or GIT files
        scanner.addDefaultExcludes();
    }

    /**
     * Returns a String[] with the values of both input parameters.
     * Duplicated values are inserted only once.
     *
     * @param list List of string
     * @param array Array of String
     * @return Array of String with all values
     */
    private String[] mergeAndConvert(List<String> list, String[] array) {
        Set<String> set = new HashSet<String>(Arrays.asList(array));
        set.addAll(list);
        return set.toArray(new String[set.size()]);
    }

}
