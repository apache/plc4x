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

package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.LinkedList;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class PlcValues {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcValues.class);

    public static PlcValue of(Boolean b) {
        return new PlcBOOL(b);
    }

    public static PlcValue of(boolean b) {
        return new PlcBOOL(b);
    }

    public static PlcValue of(Boolean[] b) {
        if(b != null) {
            if(b.length == 1) {
                return new PlcBOOL(b[0]);
            } else if(b.length > 1) {
                List<PlcBOOL> plcValues = new LinkedList<>();
                for (int i = 0; i < b.length; i++) {
                    plcValues.add(new PlcBOOL(b[i]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(boolean[] b) {
        if(b != null) {
            if(b.length == 1) {
                return new PlcBOOL(b[0]);
            } else if(b.length > 1) {
                List<PlcBOOL> plcValues = new LinkedList<>();
                for (int i = 0; i < b.length; i++) {
                    plcValues.add(new PlcBOOL(b[i]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Byte i) {
        return new PlcSINT(i);
    }

    public static PlcValue of(byte i) {
        return new PlcSINT(i);
    }

    public static PlcValue of(Byte[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcSINT(i[0]);
            } else if(i.length > 1) {
                List<PlcBYTE> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcBYTE(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(byte[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcSINT(i[0]);
            } else if(i.length > 1) {
                List<PlcBYTE> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcBYTE(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Short i) {
        return new PlcINT(i);
    }

    public static PlcValue of(short i) {
        return new PlcINT(i);
    }

    public static PlcValue of(Short[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcINT(i[0]);
            } else if(i.length > 1) {
                List<PlcINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(short[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcINT(i[0]);
            } else if(i.length > 1) {
                List<PlcINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Integer i) {
        return new PlcDINT(i);
    }

    public static PlcValue of(int i) {
        return new PlcDINT(i);
    }

    public static PlcValue of(Integer[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcDINT(i[0]);
            } else if(i.length > 1) {
                List<PlcDINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcDINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(int[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcDINT(i[0]);
            } else if(i.length > 1) {
                List<PlcDINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcDINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Long i) {
        return new PlcLINT(i);
    }

    public static PlcValue of(long i) {
        return new PlcLINT(i);
    }

    public static PlcValue of(Long[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcLINT(i[0]);
            } else if(i.length > 1) {
                List<PlcLINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcLINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(long[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcLINT(i[0]);
            } else if(i.length > 1) {
                List<PlcLINT> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcLINT(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(BigInteger i) {
        return new PlcBigInteger(i);
    }

    public static PlcValue of(BigInteger[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcBigInteger(i[0]);
            } else if(i.length > 1) {
                List<PlcBigInteger> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcBigInteger(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Float i) {
        return new PlcREAL(i);
    }

    public static PlcValue of(float i) {
        return new PlcREAL(i);
    }

    public static PlcValue of(Float[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcREAL(i[0]);
            } else if(i.length > 1) {
                List<PlcREAL> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcREAL(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(float[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcREAL(i[0]);
            } else if(i.length > 1) {
                List<PlcREAL> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcREAL(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(Double i) {
        return new PlcLREAL(i);
    }

    public static PlcValue of(double i) {
        return new PlcLREAL(i);
    }

    public static PlcValue of(Double[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcLREAL(i[0]);
            } else if(i.length > 1) {
                List<PlcLREAL> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcLREAL(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(double[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcLREAL(i[0]);
            } else if(i.length > 1) {
                List<PlcLREAL> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcLREAL(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(BigDecimal i) {
        return new PlcBigDecimal(i);
    }

    public static PlcValue of(BigDecimal[] i) {
        if(i != null) {
            if(i.length == 1) {
                return new PlcBigDecimal(i[0]);
            } else if(i.length > 1) {
                List<PlcBigDecimal> plcValues = new LinkedList<>();
                for (int j = 0; j < i.length; j++) {
                    plcValues.add(new PlcBigDecimal(i[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(String s) {
        return new PlcSTRING(s);
    }

    public static PlcValue of(String[] s) {
        if(s != null) {
            if(s.length == 1) {
                return new PlcSTRING(s[0]);
            } else if(s.length > 1) {
                List<PlcSTRING> plcValues = new LinkedList<>();
                for (int j = 0; j < s.length; j++) {
                    plcValues.add(new PlcSTRING(s[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(LocalTime s) {
        return new PlcTime(s);
    }

    public static PlcValue of(LocalTime[] s) {
        if(s != null) {
            if(s.length == 1) {
                return new PlcTime(s[0]);
            } else if(s.length > 1) {
                List<PlcTime> plcValues = new LinkedList<>();
                for (int j = 0; j < s.length; j++) {
                    plcValues.add(new PlcTime(s[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(LocalDate s) {
        return new PlcDate(s);
    }

    public static PlcValue of(LocalDate[] s) {
        if(s != null) {
            if(s.length == 1) {
                return new PlcDate(s[0]);
            } else if(s.length > 1) {
                List<PlcDate> plcValues = new LinkedList<>();
                for (int j = 0; j < s.length; j++) {
                    plcValues.add(new PlcDate(s[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(LocalDateTime s) {
        return new PlcDateTime(s);
    }

    public static PlcValue of(LocalDateTime[] s) {
        if(s != null) {
            if(s.length == 1) {
                return new PlcDateTime(s[0]);
            } else if(s.length > 1) {
                List<PlcDateTime> plcValues = new LinkedList<>();
                for (int j = 0; j < s.length; j++) {
                    plcValues.add(new PlcDateTime(s[j]));
                }
                return new PlcList(plcValues);
            }
        }
        return null;
    }

    public static PlcValue of(List<PlcValue> list) {
        return new PlcList(list);
    }

    public static PlcValue of(PlcValue... items) {
        return new PlcList(Arrays.asList(items));
    }

    public static PlcValue of(String key, PlcValue value) {
        return new PlcStruct(Collections.singletonMap(key, value));
    }

    public static PlcValue of(Map<String, PlcValue> map) {
        return new PlcStruct(map);
    }

    public static PlcValue of(Object o) {
        if(o == null) {
            return new PlcNull();
        }
        try {
            String simpleName = o.getClass().getSimpleName();
            Class<?> clazz = o.getClass();
            if (o instanceof  List) {
                simpleName = "List";
                clazz = List.class;
            } else if(clazz.isArray()) {
                simpleName = "List";
                clazz = List.class;
                Object[] objectArray = (Object[]) o;
                o = Arrays.asList(objectArray);
            }
            // If it's one of the LocalDate, LocalTime or LocalDateTime, cut off the "Local".
            if(simpleName.startsWith("Local")) {
                simpleName = simpleName.substring(5);
            }
            Constructor<?> constructor = Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + simpleName).getDeclaredConstructor(clazz);
            return ((PlcValue) constructor.newInstance(o));
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException | ClassNotFoundException e) {
            LOGGER.warn("Cannot wrap", e);
            throw new PlcIncompatibleDatatypeException(o.getClass());
        }
    }
}
