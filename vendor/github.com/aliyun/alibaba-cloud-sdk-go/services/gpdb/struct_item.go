package gpdb

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

// Item is a nested struct in gpdb response
type Item struct {
	Name                  string                            `json:"Name" xml:"Name"`
	User                  string                            `json:"User" xml:"User"`
	BlockedBySQLStmt      string                            `json:"BlockedBySQLStmt" xml:"BlockedBySQLStmt"`
	SQLPlan               string                            `json:"SQLPlan" xml:"SQLPlan"`
	SourcePort            int                               `json:"SourcePort" xml:"SourcePort"`
	NotGrantLocks         string                            `json:"NotGrantLocks" xml:"NotGrantLocks"`
	DownloadUrl           string                            `json:"DownloadUrl" xml:"DownloadUrl"`
	OperationClass        string                            `json:"OperationClass" xml:"OperationClass"`
	BlockedByPID          string                            `json:"BlockedByPID" xml:"BlockedByPID"`
	SQLStmt               string                            `json:"SQLStmt" xml:"SQLStmt"`
	QueryId               string                            `json:"QueryId" xml:"QueryId"`
	StartTime             int64                             `json:"StartTime" xml:"StartTime"`
	SQLTruncatedThreshold int                               `json:"SQLTruncatedThreshold" xml:"SQLTruncatedThreshold"`
	Duration              int                               `json:"Duration" xml:"Duration"`
	DBRole                string                            `json:"DBRole" xml:"DBRole"`
	ExecuteState          string                            `json:"ExecuteState" xml:"ExecuteState"`
	DBName                string                            `json:"DBName" xml:"DBName"`
	Cost                  int                               `json:"Cost" xml:"Cost"`
	SourceIP              string                            `json:"SourceIP" xml:"SourceIP"`
	DownloadId            int64                             `json:"DownloadId" xml:"DownloadId"`
	FileName              string                            `json:"FileName" xml:"FileName"`
	OperationExecuteTime  string                            `json:"OperationExecuteTime" xml:"OperationExecuteTime"`
	WaitingTime           int64                             `json:"WaitingTime" xml:"WaitingTime"`
	BlockedByApplication  string                            `json:"BlockedByApplication" xml:"BlockedByApplication"`
	ExecuteCost           float64                           `json:"ExecuteCost" xml:"ExecuteCost"`
	ExceptionMsg          string                            `json:"ExceptionMsg" xml:"ExceptionMsg"`
	SessionID             string                            `json:"SessionID" xml:"SessionID"`
	BlockedByUser         string                            `json:"BlockedByUser" xml:"BlockedByUser"`
	SQLText               string                            `json:"SQLText" xml:"SQLText"`
	QueryID               string                            `json:"QueryID" xml:"QueryID"`
	GrantLocks            string                            `json:"GrantLocks" xml:"GrantLocks"`
	ReturnRowCounts       int64                             `json:"ReturnRowCounts" xml:"ReturnRowCounts"`
	ScanRowCounts         int64                             `json:"ScanRowCounts" xml:"ScanRowCounts"`
	AccountName           string                            `json:"AccountName" xml:"AccountName"`
	PID                   string                            `json:"PID" xml:"PID"`
	SQLTruncated          bool                              `json:"SQLTruncated" xml:"SQLTruncated"`
	Status                string                            `json:"Status" xml:"Status"`
	Database              string                            `json:"Database" xml:"Database"`
	OperationType         string                            `json:"OperationType" xml:"OperationType"`
	Application           string                            `json:"Application" xml:"Application"`
	Series                []SeriesItemInDescribeSQLLogCount `json:"Series" xml:"Series"`
}