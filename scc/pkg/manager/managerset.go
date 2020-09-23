/*
 * Copyright 2020 Intel Corporation, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package manager

type Managerset struct {
    Overlay *OverlayObjectManager
    Proposal *ProposalObjectManager
    Hub *HubObjectManager
    HubConn *HubConnObjectManager
    HubDevice *HubDeviceObjectManager
    Device *DeviceObjectManager
    DeviceConn *DeviceConnObjectManager
    IPRange *IPRangeObjectManager
}

var mgrset = Managerset{}

func GetManagerset() *Managerset {
    return &mgrset
}
