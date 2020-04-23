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

#include <plc4c/api/types.h>
#include <plc4c/spi/types_private.h>

char *plc4c_return_code_to_message(return_code err) {
    return "hurz";
}

void plc4c_promise_set_success_callback(plc4c_promise* promise, plc4c_success_callback successCallback) {
    promise->successCallback = successCallback;
}

void plc4c_promise_set_failure_callback(plc4c_promise* promise, plc4c_failure_callback failureCallback) {
    promise->failureCallback = failureCallback;
}

bool plc4c_promise_completed(plc4c_promise* promise) {
    return promise->returnCode != UNFINISHED;
}

bool plc4c_promise_completed_successfully(plc4c_promise* promise) {
    return promise->returnCode == OK;
}

bool plc4c_promise_completed_unsuccessfully(plc4c_promise* promise) {
    return plc4c_promise_completed(promise) && !plc4c_promise_completed_successfully(promise);
}

