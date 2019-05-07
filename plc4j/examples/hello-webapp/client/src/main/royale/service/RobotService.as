/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package service
{
import org.apache.royale.events.Event;
import org.apache.royale.events.EventDispatcher;
import org.apache.royale.net.HTTPConstants;
import org.apache.royale.net.HTTPService;

[Event(name="success", type="org.apache.royale.events.Event")]
public class RobotService extends EventDispatcher {

    private var remoteService:HTTPService;
    private var _url:String = null;

    /**
     * constructor
     */
    public function RobotService() {
        remoteService = new HTTPService();
        remoteService.addEventListener(HTTPConstants.COMPLETE, completeHandler);
        _url = "api/robot/move";
    }

    private function completeHandler(event:Event):void {
        dispatchEvent(new Event("success"));
    }

    public function moveForwardLeft():void {
        remoteService.url = _url + "?direction=forward-left";
        remoteService.send();
    }

    public function moveForward():void {
        remoteService.url = _url + "?direction=forward";
        remoteService.send();
    }

    public function moveForwardRight():void {
        remoteService.url = _url + "?direction=forward-right";
        remoteService.send();
    }

    public function turnLeft():void {
        remoteService.url = _url + "?direction=left";
        remoteService.send();
    }

    public function stop():void {
        remoteService.url = _url + "?direction=stop";
        remoteService.send();
    }

    public function turnRight():void {
        remoteService.url = _url + "?direction=right";
        remoteService.send();
    }

    public function moveBackwardLeft():void {
        remoteService.url = _url + "?direction=backward-left";
        remoteService.send();
    }

    public function moveBackward():void {
        remoteService.url = _url + "?direction=backward";
        remoteService.send();
    }

    public function moveBackwardRight():void {
        remoteService.url = _url + "?direction=backward-right";
        remoteService.send();
    }

}
}