package ddoscoo

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

// Data is a nested struct in ddoscoo response
type Data struct {
	StartTime int64  `json:"StartTime" xml:"StartTime"`
	Region    string `json:"Region" xml:"Region"`
	Domain    string `json:"Domain" xml:"Domain"`
	EventType string `json:"EventType" xml:"EventType"`
	Bps       int64  `json:"Bps" xml:"Bps"`
	Port      string `json:"Port" xml:"Port"`
	Ip        string `json:"Ip" xml:"Ip"`
	Pps       int64  `json:"Pps" xml:"Pps"`
	MaxQps    int64  `json:"MaxQps" xml:"MaxQps"`
	Count     int64  `json:"Count" xml:"Count"`
	EndTime   int64  `json:"EndTime" xml:"EndTime"`
	Attack    int64  `json:"Attack" xml:"Attack"`
}
