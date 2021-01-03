package org.apache.plc4x.java.utils.connectionpool;

public interface CachedDriverManagerMBean {

    String getStateString();

    int getNumberOfConnects();

    int getNumberOfBorrows();

    int getNumberOfWachtdogs();

    int getNumberOfRejections();

    void triggerReconnect();

    int getQueueSize();

}
