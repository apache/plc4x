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

#ifndef _S7_CONTROLLER_TYPE
#define _S7_CONTROLLER_TYPE

namespace org
{
    namespace apache
    {
        namespace plc4x
        {
            namespace cpp
            {
                namespace s7
                {
                    namespace types
                    {
                        enum S7ControllerType {

                            ANY,
                            S7_300,
                            S7_400,
                            S7_1200,
                            S7_1500,
                            LOGO
                        };

                        struct S7ControllerTypeMap : public std::map<std::string, S7ControllerType>
                        {
                            S7ControllerTypeMap()
                            {
                                this->operator[]("any") = S7ControllerType::ANY;
                                this->operator[]("s7_300") = S7ControllerType::S7_300;
                                this->operator[]("s7_400") = S7ControllerType::S7_400;
                                this->operator[]("s7_1200") = S7ControllerType::S7_1200;
                                this->operator[]("s7_1500") = S7ControllerType::S7_1500;
                                this->operator[]("logo") = S7ControllerType::LOGO;
                            };
                            ~S7ControllerTypeMap() {}
                        };
                    }
                }
            }
        }
    }
}

#endif