/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.apache.plc4x.merlot.das.base.impl;

import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import org.apache.plc4x.merlot.das.base.api.BaseDriverMBean;
import java.time.LocalDateTime;
import javax.management.NotCompliantMBeanException;
import javax.management.StandardMBean;

/**
 *
 * @author cgarcia
 */
public class BaseDriverMBeanImpl extends StandardMBean implements BaseDriverMBean {

    private final BaseDriver driver;
    
    public BaseDriverMBeanImpl(BaseDriver driver) throws NotCompliantMBeanException {
        super(BaseDriverMBean.class);
        this.driver = driver;
    };

    @Override
    public void init() throws Exception {
        driver.init();
    }

    @Override
    public void destroy() throws Exception {
        driver.destroy();
    }

    @Override
    public void start() {
        driver.start();
    }

    @Override
    public void stop() {
        driver.stop();
    }

    @Override
    public void setUrl(String url) {
        driver.setUrl(url);
    }

    @Override
    public String getUrl() {
        return driver.getUrl();
    }

    @Override
    public LocalDateTime getStartDateTime() {
        return driver.getStartDateTime();
    }

    @Override
    public LocalDateTime getLastFailDateTime() {
        return driver.getLastFailDateTime();
    }

    @Override
    public long getRunningTime() {
        return driver.getRunningTime();
    }

    @Override
    public long getStoppedTime() {
        return driver.getStoppedTime();
    }

    @Override
    public long getSendMessages() {
        return driver.getSendMessages();
    }

    @Override
    public long getFailedMessages() {
        return driver.getFailedMessages();
    }

    @Override
    public long getDelay() {
        return driver.getDelay();
    }

    @Override
    public long getMaxDelay() {
        return driver.getMaxDelay();
    }

    @Override
    public long getMinDelay() {
        return driver.getMinDelay();
    }  


    @Override
    public long getRequestQueueSize() {
        return driver.getRequestQueueSize();
    }

    @Override
    public long getRequestQueueItems() {
        return driver.getRequestQueueItems();
    }

    @Override
    public long getResponseQueueSize() {
        return driver.getResponseQueueSize();
    }

    @Override
    public long getResponseQueueItems() {
        return driver.getResponseQueueItems();
    }
    

    
    
}
