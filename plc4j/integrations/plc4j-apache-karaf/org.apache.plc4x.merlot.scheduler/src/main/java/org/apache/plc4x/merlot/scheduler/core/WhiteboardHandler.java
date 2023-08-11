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
package org.apache.plc4x.merlot.scheduler.core;

import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import java.util.Date;
import org.osgi.framework.BundleContext;
import org.osgi.framework.Constants;
import org.osgi.framework.InvalidSyntaxException;
import org.osgi.framework.ServiceReference;
import org.osgi.util.tracker.ServiceTracker;
import org.osgi.util.tracker.ServiceTrackerCustomizer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


/**
 * The quartz based implementation of the scheduler.
 *
 */
public class WhiteboardHandler {

    /** Default logger. */
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    private Scheduler scheduler;

    private ServiceTracker<?,?> serviceTracker;

    public WhiteboardHandler(final BundleContext context, Scheduler scheduler) throws InvalidSyntaxException {
        this.scheduler = scheduler;
        this.serviceTracker = new ServiceTracker<>(context,
                context.createFilter("(|(" + Constants.OBJECTCLASS + "=" + Runnable.class.getName() + ")" +
                        "(" + Constants.OBJECTCLASS + "=" + Job.class.getName() + "))"),
                new ServiceTrackerCustomizer<Object,Object>() {

                    public synchronized void  removedService(final ServiceReference reference, final Object service) {
                        context.ungetService(reference);
                        unregister(reference, service);
                    }

                    public synchronized void modifiedService(final ServiceReference reference, final Object service) {
                        unregister(reference, service);
                        register(reference, service);
                    }

                    public synchronized Object addingService(final ServiceReference reference) {
                        final Object obj = context.getService(reference);
                        if ( obj != null ) {
                            register(reference, obj);
                        }
                        return obj;
                    }
                });
        this.serviceTracker.open();
    }

    /**
     * Deactivate this component.
     */
    public void deactivate() {
        this.serviceTracker.close();
    }


    /**
     * Create unique identifier
     */
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

    /**
     * Register a job or task
     */
    private void register(final ServiceReference ref, final Object job) {
        final String name = getServiceIdentifier(ref);
        Boolean concurrent = true;
        if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_CONCURRENT) != null) {
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
                            times = Integer.valueOf(v.toString());
                        }
                    }
                }

                Long period = null;
                if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_PERIOD) != null) {
                    if (ref.getProperty(Scheduler.PROPERTY_SCHEDULER_PERIOD) instanceof Long) {
                        period = (Long) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_PERIOD);
                    } else {
                        period = Long.valueOf((String) ref.getProperty(Scheduler.PROPERTY_SCHEDULER_PERIOD));
                    }
                    if (period < 1) {
                        this.logger.debug("Ignoring service {} : scheduler period is less than 1.", ref);
                    } else if (times < -1) {
                        this.logger.debug("Ignoring service {} : scheduler times is defined but is less than -1.", ref);
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
                            date.setTime(System.currentTimeMillis() + period * 1000);
                        }
                        this.scheduler.schedule(job, this.scheduler.AT(date, times, period)
                                .name(name)
                                .canRunConcurrently((concurrent != null ? concurrent : true)));
                    }
                } else {
                    this.logger.debug("Ignoring service {} : no scheduling property found.", ref);
                }
            }
        } catch (Exception e) {
            logger.warn("Error scheduling job", e);
        }
    }

    /**
     * Unregister a service.
     */
    private void unregister(final ServiceReference reference, final Object service) {
        final String name = getServiceIdentifier(reference);
        this.scheduler.unschedule(name);
    }
}
