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

[enum uint 8 BACnetDaysOfWeek
    ['0'  MONDAY                        ]
    ['1'  TUESDAY                       ]
    ['2'  WEDNESDAY                     ]
    ['3'  THURSDAY                      ]
    ['4'  FRIDAY                        ]
    ['5'  SATURDAY                      ]
    ['6'  SUNDAY                        ]
]

[enum uint 8 BACnetEventTransitionBits
    ['0'  TO_OFFNORMAL                  ]
    ['1'  TO_FAULT                      ]
    ['2'  TO_NORMAL                     ]
]

[enum uint 8 BACnetLimitEnable
    ['0'  LOW_LIMIT_ENABLE              ]
    ['1'  HIGH_LIMIT_ENABLE             ]
]

[enum uint 8 BACnetLogStatus
    ['0'  LOG_DISABLED                  ]
    ['1'  BUFFER_PURGED                 ]
    ['2'  LOG_INTERRUPTED               ]
]

[enum uint 8 BACnetObjectTypesSupported
    ['0'    ANALOG_INPUT                ]
    ['1'    ANALOG_OUTPUT               ]
    ['2'    ANALOG_VALUE                ]
    ['3'    BINARY_INPUT                ]
    ['4'    BINARY_OUTPUT               ]
    ['5'    BINARY_VALUE                ]
    ['6'    CALENDAR                    ]
    ['7'    COMMAND                     ]
    ['8'    DEVICE                      ]
    ['9'    EVENT_ENROLLMENT            ]
    ['10'   FILE                        ]
    ['11'   GROUP                       ]
    ['12'   LOOP                        ]
    ['13'   MULTI_STATE_INPUT           ]
    ['14'   MULTI_STATE_OUTPUT          ]
    ['15'   NOTIFICATION_CLASS          ]
    ['16'   PROGRAM                     ]
    ['17'   SCHEDULE                    ]
    ['18'   AVERAGING                   ]
    ['19'   MULTI_STATE_VALUE           ]
    ['20'   TREND_LOG                   ]
    ['21'   LIFE_SAFETY_POINT           ]
    ['22'   LIFE_SAFETY_ZONE            ]
    ['23'   ACCUMULATOR                 ]
    ['24'   PULSE_CONVERTER             ]
    ['25'   EVENT_LOG                   ]
    ['26'   GLOBAL_GROUP                ]
    ['27'   TREND_LOG_MULTIPLE          ]
    ['28'   LOAD_CONTROL                ]
    ['29'   STRUCTURED_VIEW             ]
    ['30'   ACCESS_DOOR                 ]
    ['31'   TIMER                       ]
    ['32'   ACCESS_CREDENTIAL           ]
    ['33'   ACCESS_POINT                ]
    ['34'   ACCESS_RIGHTS               ]
    ['35'   ACCESS_USER                 ]
    ['36'   ACCESS_ZONE                 ]
    ['37'   CREDENTIAL_DATA_INPUT       ]
    ['38'   NETWORK_SECURITY            ]
    ['39'   BITSTRING_VALUE             ]
    ['40'   CHARACTERSTRING_VALUE       ]
    ['41'   DATEPATTERN_VALUE           ]
    ['42'   DATE_VALUE                  ]
    ['43'   DATETIMEPATTERN_VALUE       ]
    ['44'   DATETIME_VALUE              ]
    ['45'   INTEGER_VALUE               ]
    ['46'   LARGE_ANALOG_VALUE          ]
    ['47'   OCTETSTRING_VALUE           ]
    ['48'   POSITIVE_INTEGER_VALUE      ]
    ['49'   TIMEPATTERN_VALUE           ]
    ['50'   TIME_VALUE                  ]
    ['51'   NOTIFICATION_FORWARDER      ]
    ['52'   ALERT_ENROLLMENT            ]
    ['53'   CHANNEL                     ]
    ['54'   LIGHTING_OUTPUT             ]
    ['55'   BINARY_LIGHTING_OUTPUT      ]
    ['56'   NETWORK_PORT                ]
    ['57'   ELEVATOR_GROUP              ]
    ['58'   ESCALATOR                   ]
    ['59'   LIFT                        ]
]

[enum uint 8 BACnetResultFlags
    ['0'  FIRST_ITEM                    ]
    ['1'  LAST_ITEM                     ]
    ['2'  MORE_ITEMS                    ]
]

[enum uint 8 BACnetServicesSupported
// Alarm and Event Services
    ['0'  ACKNOWLEDGE_ALARM                     ]
    ['1'  CONFIRMED_COV_NOTIFICATION            ]
    //['42'  CONFIRMED_COV_NOTIFICATION_MULTIPLE]
    ['2'  CONFIRMED_EVENT_NOTIFICATION          ]
    ['3'  GET_ALARM_SUMMARY                     ]
    ['4'  GET_ENROLLMENT_SUMMARY                ]
    //['39'  GET_EVENT_INFORMATION              ]
    //['37'  LIFE_SAFETY_OPERATION              ]
    ['5'  SUBSCRIBE_COV                         ]
    //['38'  SUBSCRIBE_COV_PROPERTY             ]
    //['41'  SUBSCRIBE_COV_PROPERTY_MULTIPLE    ]
// File Access Services
    ['6'   ATOMIC_READ_FILE                     ]
    ['7'   ATOMIC_WRITE_FILE                    ]
// Object Access Services
    ['8'   ADD_LIST_ELEMENT                     ]
    ['9'   REMOVE_LIST_ELEMENT                  ]
    ['10'  CREATE_OBJECT                        ]
    ['11'  DELETE_OBJECT                        ]
    ['12'  READ_PROPERTY                        ]
    ['14'  READ_PROPERTY_MULTIPLE               ]
    //['35'  READ_RANGE                         ]
    //['40'  WRITE_GROUP                        ]
    ['15'  WRITE_PROPERTY                       ]
    ['16'  WRITE_PROPERTY_MULTIPLE              ]
// Remote Device Management Services
    ['17'  DEVICE_COMMUNICATION_CONTROL         ]
    ['18'  CONFIRMED_PRIVATE_TRANSFER           ]
    ['19'  CONFIRMED_TEXT_MESSAGE               ]
    ['20'  REINITIALIZE_DEVICE                  ]
// Virtual Terminal Services
    ['21'  VT_OPEN                              ]
    ['22'  VT_CLOSE                             ]
    ['23'  VT_DATA                              ]
// Removed Services
    // formerly: ['13'  READ_PROPERTY_CONDITIONAL ] removed in version 1 revision 12
    // formerly: ['24'  AUTHENTICATE ]removed in version 1 revision 11
    // formerly: ['25'  REQUEST_KEY  ] removed in version 1 revision 11
// Unconfirmed Services
    ['26'  I_AM                                 ]
    ['27'  I_HAVE                               ]
    ['28'  UNCONFIRMED_COV_NOTIFICATION         ]
    // ['43'  UNCONFIRMED_COV_NOTIFICATION_MULTIPLE]
    ['29'  UNCONFIRMED_EVENT_NOTIFICATION       ]
    ['30'  UNCONFIRMED_PRIVATE_TRANSFER         ]
    ['31'  UNCONFIRMED_TEXT_MESSAGE             ]
    ['32'  TIME_SYNCHRONIZATION                 ]
    // ['36'  UTC_TIME_SYNCHRONIZATION          ]
    ['33'  WHO_HAS                              ]
    ['34'  WHO_IS                               ]
// Services added after 1995
    ['35'  READ_RANGE                           ] // Object Access Service
    ['36'  UTC_TIME_SYNCHRONIZATION             ] // Remote Device Management Service
    ['37'  LIFE_SAFETY_OPERATION                ] // Alarm and Event Service
    ['38'  SUBSCRIBE_COV_PROPERTY               ] // Alarm and Event Service
    ['39'  GET_EVENT_INFORMATION                ] // Alarm and Event Service
    ['40'  WRITE_GROUP                          ] // Object Access Services
// Services added after 2012
    ['41'  SUBSCRIBE_COV_PROPERTY_MULTIPLE      ] // Alarm and Event Service
    ['42'  CONFIRMED_COV_NOTIFICATION_MULTIPLE  ] // Alarm and Event Service
    ['43'  UNCONFIRMED_COV_NOTIFICATION_MULTIPLE]// Alarm and Event Service
]

[enum uint 8 BACnetStatusFlags
    ['0'  IN_ALARM              ]
    ['1'  FAULT                 ]
    ['2'  OVERRIDDEN            ]
    ['3'  OUT_OF_SERVICE        ]
]
