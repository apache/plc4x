/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.apache.plc4x.java.s7.readwrite.field;

import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

/**
 *
 * @author cgarcia
 */
public class S7SslField  implements PlcField  {
    
   //SZL_ID=0xYYYY;INDEX=0xZZZZ
    private static final Pattern SSL_ADDRESS_PATTERN =
        Pattern.compile("^SSL_ID=(?<sslId>16#[0-9a-fA-F]{4});INDEX=(?<index>16#[0-9a-fA-F]{4})");  
    
    private static final String SSL_ID = "sslId";    
    private static final String INDEX = "index";    
    
    private final int szlId;
    private final int index;

    public S7SslField(int sslId, int index) {
        this.szlId = sslId;
        this.index = index;
    }

    public int getSslId() {
        return szlId;
    }

    public int getIndex() {
        return index;
    }
          
    public static boolean matches(String fieldString) {
        return SSL_ADDRESS_PATTERN.matcher(fieldString).matches();      
    }   
    
    public static S7SslField of(String fieldString) {
        Matcher matcher = SSL_ADDRESS_PATTERN.matcher(fieldString); 
        if (matcher.matches()){
            String strSxlId = matcher.group(SSL_ID);
            String strIndex = matcher.group(INDEX);
            strSxlId = strSxlId.replaceAll("16#", "");
            strIndex = strIndex.replaceAll("16#", "");
            return new S7SslField(Integer.parseInt(strSxlId, 16),Integer.parseInt(strIndex, 16));
        }  
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString); 
    }       
}
