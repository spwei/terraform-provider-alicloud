package emr

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

// ClusterInstalledService is a nested struct in emr response
type ClusterInstalledService struct {
	ServiceName        string                                         `json:"ServiceName" xml:"ServiceName"`
	ServiceDisplayName string                                         `json:"ServiceDisplayName" xml:"ServiceDisplayName"`
	ServiceVersion     string                                         `json:"ServiceVersion" xml:"ServiceVersion"`
	ServiceEcmVersion  string                                         `json:"ServiceEcmVersion" xml:"ServiceEcmVersion"`
	ServiceStatus      string                                         `json:"serviceStatus" xml:"serviceStatus"`
	OnlyClient         bool                                           `json:"onlyClient" xml:"onlyClient"`
	NotStartedNum      int                                            `json:"notStartedNum" xml:"notStartedNum"`
	NeedRestartNum     int                                            `json:"needRestartNum" xml:"needRestartNum"`
	AbnormalNum        int                                            `json:"abnormalNum" xml:"abnormalNum"`
	Comment            string                                         `json:"comment" xml:"comment"`
	State              string                                         `json:"State" xml:"State"`
	ServiceActionList  ServiceActionListInListClusterInstalledService `json:"ServiceActionList" xml:"ServiceActionList"`
}
