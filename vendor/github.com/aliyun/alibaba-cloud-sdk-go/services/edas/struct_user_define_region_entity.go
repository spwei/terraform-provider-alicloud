package edas

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// UserDefineRegionEntity is a nested struct in edas response
type UserDefineRegionEntity struct {
	RegionId      string `json:"RegionId" xml:"RegionId"`
	RegistryType  string `json:"RegistryType" xml:"RegistryType"`
	BelongRegion  string `json:"BelongRegion" xml:"BelongRegion"`
	DebugEnable   bool   `json:"DebugEnable" xml:"DebugEnable"`
	UserId        string `json:"UserId" xml:"UserId"`
	Id            int64  `json:"Id" xml:"Id"`
	RegionName    string `json:"RegionName" xml:"RegionName"`
	Description   string `json:"Description" xml:"Description"`
	MseInstanceId string `json:"MseInstanceId" xml:"MseInstanceId"`
}
