/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.merlot.scheduler.impl;

import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerWhiteBoard;
import java.util.Date;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.ServiceReference;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class SchedulerWhiteBoardImpl implements SchedulerWhiteBoard {
    
    private BundleContext bc = null;
    
    private Scheduler scheduler = null;
    
    private final Logger LOGGER = LoggerFactory.getLogger(this.getClass());

    @Override
    public void init() {

    }

    @Override
    public void destroy() {

    }

    @Override
    public void start() {
        
    }

    @Override
    public void stop() {

    }

    @Override
    public void bind(ServiceReference ref) {
        final String name = getServiceIdentifier(ref);
        final Object job = bc.getService(ref);
        Boolean concurrent = true;
       
        if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_CONCURRENT) != null){
            if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_CONCURRENT) instanceof Boolean) {
                concurrent = (Boolean) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_CONCURRENT);
            } else {
                concurrent = Boolean.valueOf((String) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_CONCURRENT));  
            }            
        }
        final String expression = (String) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_EXPRESSION);
        try {
            if (expression != null) {
                this.scheduler.schedule(job, this.scheduler.EXPR(expression)
                        .name(name)
                        .canRunConcurrently(concurrent));
            } else {
                Integer times = -1;
                {
                    final Object v = ref.getProperty(Scheduler.PROPERTY_SCHEDULER_TIMES);
                    if (null != v) {
                        if (v instanceof Integer) {
                            times = (Integer) v;
                        } else if (v instanceof Long) {
                            times = ((Long) v).intValue();
                        } else if (v instanceof Number) {
                            times = ((Number) v).intValue();
                        } else {
                            times = new Integer(v.toString());
                        }
                    }
                }

                Long period = null;
                if (ref.getProperty("MERLOT_SCHEDULER_PERIOD") != null) {
                    if (ref.getProperty("MERLOT_SCHEDULER_PERIOD") instanceof Long) {
                        period = (Long) ref.getProperty("MERLOT_SCHEDULER_PERIOD");
                    } else {
                        period = new Long((String) ref.getProperty("MERLOT_SCHEDULER_PERIOD"));
                    }
                    if (period < 1) {
                        LOGGER.debug("Ignoring service {} : scheduler period is less than 1.", job);
                    } else if (times < -1) {
                        LOGGER.debug("Ignoring service {} : scheduler times is defined but is less than -1.", job);
                    } else {
                        boolean immediate = false;
                        if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_IMMEDIATE) != null) {
                            if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_IMMEDIATE) instanceof Boolean) {
                                immediate = (Boolean) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_IMMEDIATE);
                            } else {
                                immediate = Boolean.valueOf((String) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_IMMEDIATE));
                            }
                        }
                        final Date date = new Date();
                        if (!immediate) {     
                            
                            date.setTime(System.currentTimeMillis() + period);
                        }
                        System.out.println("Paso por aqui >>> ....: " + period);
                        this.scheduler.schedule(job, this.scheduler.NOW(times, period)
                                .name(name)
                                .canRunConcurrently((concurrent != null ? concurrent : true)));
                    }
                } else {
                    LOGGER.info("Ignoring service {} : no scheduling property found.", job);
                }
            }
        } catch (Exception e) {
            LOGGER.warn("Error scheduling job", e);
        }        
    }

    @Override
    public void unbind(ServiceReference ref) {
        final String name = getServiceIdentifier(ref);
        this.scheduler.unschedule(name);
    }

    @Override
    public void setScheduler(Scheduler scheduler) {
        this.scheduler = scheduler;
    }

    @Override
    public Scheduler getScheduler() {
        return this.scheduler;
    }    
    

    private String getServiceIdentifier(final ServiceReference ref) {
        String name = (String) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_NAME);
        if ( name == null ) {
            name = (String) ref.getProperty(Constants.SERVICE_PID);
            if ( name == null ) {
                name = "Registered Service";
            }
        }
        // now append service id to create a unique identifier
        name = name + "." + ref.getProperty(Constants.SERVICE_ID);
        return name;
    }

    @Override
    public void setBundleContext(BundleContext bc) {
        this.bc = bc;
    }



    
}
