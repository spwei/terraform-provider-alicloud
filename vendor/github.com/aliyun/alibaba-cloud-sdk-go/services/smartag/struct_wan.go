package smartag

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

// Wan is a nested struct in smartag response
type Wan struct {
	Weight       int    `json:"Weight" xml:"Weight"`
	Mask         string `json:"Mask" xml:"Mask"`
	Gateway      string `json:"Gateway" xml:"Gateway"`
	PortName     string `json:"PortName" xml:"PortName"`
	IPType       string `json:"IPType" xml:"IPType"`
	Priority     int    `json:"Priority" xml:"Priority"`
	BandWidth    int    `json:"BandWidth" xml:"BandWidth"`
	TrafficState string `json:"TrafficState" xml:"TrafficState"`
	IP           string `json:"IP" xml:"IP"`
	ISP          string `json:"ISP" xml:"ISP"`
	Username     string `json:"Username" xml:"Username"`
	Vlan         string `json:"Vlan" xml:"Vlan"`
	SourceIps    string `json:"SourceIps" xml:"SourceIps"`
	StartIp      string `json:"StartIp" xml:"StartIp"`
	StopIp       string `json:"StopIp" xml:"StopIp"`
}
