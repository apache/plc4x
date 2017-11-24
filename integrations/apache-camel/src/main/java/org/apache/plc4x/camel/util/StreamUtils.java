package org.apache.plc4x.camel.util;

import java.util.Optional;
import java.util.stream.Stream;

public class StreamUtils {

    public static <T> Stream<T> streamOf(@SuppressWarnings("OptionalUsedAsFieldOrParameterType") Optional<T> optional) {
        return Stream.of(optional)
            .filter(Optional::isPresent)
            .map(Optional::get);
    }
}
