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

#include "ErrorCategory.h"

namespace org
{
    namespace apache
    {
        namespace plc4x
        {
            namespace cpp
            {
                namespace utils
                {

                    const char* Error_Category::_paErrorMessage[] =
                    {
                        "success",
                        "err0",
                        "err1",
                        "err2",
                        "err3",
                        "err4",
                        "err5",
                        "err6"
                    };


                    std::string Error_Category::message(int iEv) const
                    {
                        std::string strMessage = "";
                        if (IsValidMessageIndex(iEv))
                        {
                            strMessage = _paErrorMessage[iEv];
                        }
                        else
                        {
                            strMessage = DefaultMessage(iEv);
                        }

                        return strMessage;
                    }

                    std::string Error_Category::DefaultMessage(int iEv) const
                    {
                        boost::system::error_code ecCode(iEv, boost::system::system_category());

                        return ecCode.message();
                    }
                }
            }
        }
    }
}
