/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.camel;

import org.apache.camel.*;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriEndpoint;
import org.apache.camel.spi.UriParam;
import org.apache.camel.spi.UriPath;
import org.apache.camel.support.DefaultEndpoint;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import java.util.List;
import java.util.Map;
import java.util.Objects;

@UriEndpoint(scheme = "plc4x", title = "PLC4X", syntax = "plc4x:driver", label = "plc4x")
public class Plc4XEndpoint extends DefaultEndpoint {

    /**
     * The name of the PLC4X driver
     */
    @UriPath(label = "common")
    @Metadata(required = true)
    @SuppressWarnings("unused")
    private String driver;

    /**
     * The address for the PLC4X driver
     */
    @UriParam
    @Metadata(required = false)
    @SuppressWarnings("unused")
    private List<TagData> tags;

    /**
     * TODO: document me
     */
    @UriParam
    @Metadata(required = false)
    @SuppressWarnings("unused")
    private Class dataType ;

    private final PlcDriverManager plcDriverManager;
    private  PlcConnection connection;
    private String uri;

    public Plc4XEndpoint(String endpointUri, Component component) {
        super(endpointUri, component);
        plcDriverManager= new PlcDriverManager();
        uri = endpointUri;
        //Here I established the connection in the endpoint, as it is created once during the context
        // to avoid disconnecting and reconnecting for every request
        try {
            String plc4xURI = uri.replaceFirst("plc4x:/?/?", "");
            connection = plcDriverManager.getConnection(plc4xURI);
        } catch (PlcConnectionException e) {
            e.printStackTrace();
        }
    }

    public PlcConnection getConnection() {
        return connection;
    }


    @Override
    public Producer createProducer() throws Exception {
        //Checking if connection is still up and reconnecting if not
        if(!connection.isConnected()){
            try{
                connection= plcDriverManager.getConnection(uri.replaceFirst("plc4x:/?/?", ""));
            }
            catch (Exception e){
                e.printStackTrace();
            }
        }
        return new Plc4XProducer(this);
    }

    @Override
    public Consumer createConsumer(Processor processor) throws Exception {
        //Checking if connection is still up and reconnecting if not
        if(!connection.isConnected()){
            try{
                connection= plcDriverManager.getConnection(uri.replaceFirst("plc4x:/?/?", ""));
            }
            catch (Exception e){
                e.printStackTrace();
            }
        }
        return new Plc4XConsumer(this, processor);
    }

    @Override
    public PollingConsumer createPollingConsumer() throws Exception {
        //Checking if connection is still up and reconnecting if not
        if(!connection.isConnected()){
            try{
                connection= plcDriverManager.getConnection(uri.replaceFirst("plc4x:/?/?", ""));
            }
            catch (Exception e){
                e.printStackTrace();
            }
        }
        return new Plc4XPollingConsumer(this);
    }

    @Override
    public boolean isSingleton() {
        return true;
    }

    public PlcDriverManager getPlcDriverManager() {
        return plcDriverManager;
    }

    public String getDriver() {
        return driver;
    }

    public void setDriver(String driver) {
        this.driver = driver;
    }

    public List<TagData> getTags() {
        return tags;
    }

    public void setTags(List<TagData> tags) {
        this.tags = tags;
    }

    public Class getDataType() {
        return dataType;
    }

    public void setDataType(Class dataType) {
        this.dataType = dataType;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Plc4XEndpoint)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        Plc4XEndpoint that = (Plc4XEndpoint) o;
        return Objects.equals(getDriver(), that.getDriver()) &&
            Objects.equals(getTags(), that.getTags()) &&
            Objects.equals(getDataType(), that.getDataType()) &&
            Objects.equals(getPlcDriverManager(), that.getPlcDriverManager());
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), getDriver(), getTags(), getDataType(), getPlcDriverManager());
    }

    @Override
    public void doStop(){
        //Shutting down the connection when leaving the Context
        try{
            if(connection!=null){
                if(connection.isConnected()){
                    connection.close();
                }
            }
        }
        catch (Exception e){
            e.printStackTrace();
        }
    }

}
