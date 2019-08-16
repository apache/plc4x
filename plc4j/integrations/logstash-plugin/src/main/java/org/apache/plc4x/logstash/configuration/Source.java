package org.apache.plc4x.logstash.configuration;

public class Source {
    private final String connectionUri;
    private final String sourceAlias;

    @Override
    public String toString() {
        return "Source{" +
            "connectionUri='" + connectionUri + '\'' +
            ", sourceAlias='" + sourceAlias + '\'' +
            '}';
    }

    public String getConnectionUri() {
        return connectionUri;
    }

    public String getSourceAlias() {
        return sourceAlias;
    }

    public Source(String connectionUri, String sourceAlias) {
        this.connectionUri = connectionUri;
        this.sourceAlias = sourceAlias;
    }
}
