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

#include <stdlib.h>
#include <plc4c/write.h>
#include <plc4c/spi/types_private.h>

plc4c_return_code plc4c_write_request_execute(plc4c_write_request *write_request,
                                       plc4c_write_request_execution **write_request_execution) {
    plc4c_system_task *system_task;
    write_request->connection->driver->write_function(&system_task);
    // Add the new task to the task-list.
    plc4c_utils_list_insert_tail_value(write_request->connection->system->task_list, system_task);

    // Inject the default write context into the system task.
    plc4c_write_request_execution *new_write_request_execution = malloc(sizeof(plc4c_write_request_execution));
    new_write_request_execution->write_request = write_request;
    new_write_request_execution->write_response = NULL;
    new_write_request_execution->system_task = system_task;
    new_write_request_execution->system_task->context = new_write_request_execution;

    *write_request_execution = new_write_request_execution;
    return OK;
}

bool plc4c_write_request_finished_successfully(plc4c_write_request_execution *write_request_execution) {
    if(write_request_execution == NULL) {
        return true;
    }
    if(write_request_execution->system_task == NULL) {
        return true;
    }
    return write_request_execution->system_task->completed;
}

bool plc4c_write_request_has_error(plc4c_write_request_execution *write_request_execution) {
    return false;
}

plc4c_write_response *plc4c_write_request_get_response(plc4c_write_request_execution *write_request_execution) {
    if(write_request_execution == NULL) {
        return NULL;
    }
    return write_request_execution->write_response;
}

void plc4c_write_request_destroy(plc4c_write_request *write_request) {
    free(write_request);
}

void plc4c_write_request_execution_destroy(plc4c_write_request_execution *write_request_execution) {
    free(write_request_execution);
}
