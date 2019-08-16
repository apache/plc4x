package org.apache.plc4x.logstash.configuration;

import java.util.List;

public class Job {
    private final Integer rate;
    private final List<String> queries;
    private final List<String> sources;

    public Job(Integer rate, List<String> queries, List<String> sources) {
        this.rate = rate;
        this.queries = queries;
        this.sources = sources;
    }

    @Override
    public String toString() {
        return "Job{" +
            "rate=" + rate +
            ", queries=" + queries +
            ", sources=" + sources +
            '}';
    }

    public Integer getRate() {
        return rate;
    }

    public List<String> getQueries() {
        return queries;
    }

    public List<String> getSources() {
        return sources;
    }
}
