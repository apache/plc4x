package org.apache.plc4x;

import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.*;

class Plc4xSchemaFactoryTest implements WithAssertions {

    @Test
    void create() {
        Plc4xSchemaFactory factory = new Plc4xSchemaFactory();
        assertThatThrownBy(() -> factory.create(null, "", Collections.emptyMap()))
            .isInstanceOf(NullPointerException.class)
            .hasMessageContaining("Please specify operand 'config'...");
    }
}