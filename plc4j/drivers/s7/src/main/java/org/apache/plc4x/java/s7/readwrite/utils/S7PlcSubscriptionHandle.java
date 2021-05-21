/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.apache.plc4x.java.s7.readwrite.utils;

import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.readwrite.types.EventType;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;

/**
 *
 * @author cgarcia
 */
public class S7PlcSubscriptionHandle extends DefaultPlcSubscriptionHandle {
    
    private EventType eventtype;

    public S7PlcSubscriptionHandle(EventType eventtype, PlcSubscriber plcSubscriber) {
        super(plcSubscriber);
        this.eventtype = eventtype;
    }
    
    public EventType getEventType() {
        return eventtype;   
    }
    
}
