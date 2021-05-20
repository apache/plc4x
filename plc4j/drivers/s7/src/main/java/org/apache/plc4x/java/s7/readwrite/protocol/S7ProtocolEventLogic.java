/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.apache.plc4x.java.s7.readwrite.protocol;

import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;

/**
 *
 * @author cgarcia
 */
public class S7ProtocolEventLogic implements PlcSubscriptionHandle {

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }
    
    
    private class EventDispacher implements Runnable {

        @Override
        public void run() {
            while(true){
                System.out.println("Paso por aqui...");
                try {
                    Thread.sleep(5000);
                } catch (InterruptedException ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
        }
        
    }
    
    
}
