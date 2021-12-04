/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.events;

import java.util.Map;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;

/**
 * Like JMS but different.
 * Maintain the same pattern as JMS MapMessage. 
 * S7Event is a Map whose key is a String and its value is an 
 * Object of primitive values.
 * Classes derived from the S7Event must respect the following 
 * conversion table.
 * 
 *  |        | boolean byte short char int long float double String byte[]
 *  |--------+-------------------------------------------------------------
 *  |boolean |    X                                            X
 *  |byte    |          X     X         X   X                  X
 *  |short   |                X         X   X                  X
 *  |char    |                     X                           X
 *  |int     |                          X   X                  X
 *  |long    |                              X                  X
 *  |float   |                                    X     X      X
 *  |double  |                                          X      X
 *  |String  |    X     X     X         X   X     X     X      X
 *  |byte[]  |                                                        X
 *  |--------+-------------------------------------------------------------
 */
public interface S7Event extends PlcSubscriptionEvent{
     Map<String, Object> getMap();
}
