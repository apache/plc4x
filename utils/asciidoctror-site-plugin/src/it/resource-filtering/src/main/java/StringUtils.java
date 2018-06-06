package example;

/**
 * Using ${java.vendor} ${java.version}
 *
 * Command property: ${command.property}
 * POM property: ${pom.property}
 */
public class StringUtils {
    // tag::contains[]
    public boolean contains(String haystack, String needle) {
        return haystack.contains(needle);
    }
    // end::contains[]
}
