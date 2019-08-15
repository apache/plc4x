package org.apache.plc4x.logstash;

import co.elastic.logstash.api.Configuration;
import org.apache.commons.lang3.StringUtils;
import org.junit.Assert;
import org.junit.Test;
import org.logstash.plugins.ConfigurationImpl;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

public class JavaInputExampleTest {

    @Test
    public void testJavaInputExample() {
        String prefix = "This is message";
        long eventCount = 5;
        Map<String, Object> configValues = new HashMap<>();
        configValues.put(Plc4xInput.PREFIX_CONFIG.name(), prefix);
        configValues.put(Plc4xInput.EVENT_COUNT_CONFIG.name(), eventCount);
        Configuration config = new ConfigurationImpl(configValues);
        Plc4xInput input = new Plc4xInput("test-id", config, null);
        TestConsumer testConsumer = new TestConsumer();
        input.start(testConsumer);

        List<Map<String, Object>> events = testConsumer.getEvents();
        Assert.assertEquals(eventCount, events.size());
        for (int k = 1; k <= events.size(); k++) {
            Assert.assertEquals(prefix + " " + StringUtils.center(k + " of " + eventCount, 20),
                    events.get(k - 1).get("message"));
        }
    }

    private static class TestConsumer implements Consumer<Map<String, Object>> {

        private List<Map<String, Object>> events = new ArrayList<>();

        @Override
        public void accept(Map<String, Object> event) {
            synchronized (this) {
                events.add(event);
            }
        }

        public List<Map<String, Object>> getEvents() {
            return events;
        }
    }

}
