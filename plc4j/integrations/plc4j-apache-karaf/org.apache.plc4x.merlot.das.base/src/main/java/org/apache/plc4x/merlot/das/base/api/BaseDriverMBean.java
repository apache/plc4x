/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.apache.plc4x.merlot.das.base.api;

import java.time.LocalDateTime;

/**
 *
 * @author cgarcia
 */
public interface BaseDriverMBean {
    
    public void init() throws Exception;

    public void destroy() throws Exception;

    public void start();

    public void stop();	
        
    public void setUrl(String url);
    
    public String getUrl();   
    
    public LocalDateTime getStartDateTime();
    
    public LocalDateTime getLastFailDateTime();
    
    public long getRunningTime();
    
    public long getStoppedTime();    
    
    public long getSendMessages();
    
    public long getFailedMessages();          
        
    public long getDelay();
    
    public long getMaxDelay();
    
    public long getMinDelay();   
       
    public long getRequestQueueSize();
    
    public long getRequestQueueItems();   

    public long getResponseQueueSize();
    
    public long getResponseQueueItems();      
    
}
