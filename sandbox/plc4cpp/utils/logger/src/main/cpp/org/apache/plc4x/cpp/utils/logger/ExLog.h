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

#ifndef _EX_LOG_H_
#define _EX_LOG_H_

#include <boost/exception/exception.hpp>
#include <exception>
#include <string>

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
                    typedef boost::error_info<struct tag_errmsg, std::string> errmsg_info;

                    class ExLog : public boost::exception, public std::exception
                    {

                    public:
                        ExLog();

                        virtual ~ExLog() throw();

                        ExLog(unsigned int uiUserMessageID, std::exception& ex, const char* pError);

                        ExLog(unsigned int uiUserMessageID, boost::exception& ex, const char* pError);

                        void logInfo(std::string strInfo);

                        inline unsigned int getUserMessageID() const { return _uiUserMessageID; }
                        virtual const char* what() const throw() { return _strWhat.c_str(); }

                    protected:

                        std::string printStackTrace(unsigned int max_frames = 63);

                    protected:

                        std::string _strWhat;
                        std::string _strLastErrorMessage;
                        unsigned int _uiUserMessageID;
                    };
                }
            }
        }
    }
}
#endif
