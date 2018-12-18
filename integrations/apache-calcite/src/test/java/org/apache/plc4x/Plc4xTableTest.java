package org.apache.plc4x;

import org.apache.calcite.linq4j.Enumerator;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.util.concurrent.ArrayBlockingQueue;

class Plc4xTableTest implements WithAssertions {

    @Test
    void testOnBlockingQueue() {
        ArrayBlockingQueue<Object[]> queue = new ArrayBlockingQueue<Object[]>(100);
        Plc4xTable table = new Plc4xTable(queue, null);

        Object[] objects = new Object[0];
        queue.add(objects);

        Enumerator<Object[]> enumerator = table.scan(null).enumerator();

        assertThat(enumerator.moveNext()).isTrue();
        assertThat(enumerator.current()).isEqualTo(objects);
    }
}