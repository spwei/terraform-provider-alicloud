package drds

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeShardTaskList invokes the drds.DescribeShardTaskList API synchronously
func (client *Client) DescribeShardTaskList(request *DescribeShardTaskListRequest) (response *DescribeShardTaskListResponse, err error) {
	response = CreateDescribeShardTaskListResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeShardTaskListWithChan invokes the drds.DescribeShardTaskList API asynchronously
func (client *Client) DescribeShardTaskListWithChan(request *DescribeShardTaskListRequest) (<-chan *DescribeShardTaskListResponse, <-chan error) {
	responseChan := make(chan *DescribeShardTaskListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeShardTaskList(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeShardTaskListWithCallback invokes the drds.DescribeShardTaskList API asynchronously
func (client *Client) DescribeShardTaskListWithCallback(request *DescribeShardTaskListRequest, callback func(response *DescribeShardTaskListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeShardTaskListResponse
		var err error
		defer close(result)
		response, err = client.DescribeShardTaskList(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeShardTaskListRequest is the request struct for api DescribeShardTaskList
type DescribeShardTaskListRequest struct {
	*requests.RpcRequest
	TaskType       string           `position:"Query" name:"TaskType"`
	Query          string           `position:"Query" name:"Query"`
	CurrentPage    requests.Integer `position:"Query" name:"CurrentPage"`
	DrdsInstanceId string           `position:"Query" name:"DrdsInstanceId"`
	DbName         string           `position:"Query" name:"DbName"`
	PageSize       requests.Integer `position:"Query" name:"PageSize"`
}

// DescribeShardTaskListResponse is the response struct for api DescribeShardTaskList
type DescribeShardTaskListResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	Success    bool       `json:"Success" xml:"Success"`
	PageNumber int        `json:"PageNumber" xml:"PageNumber"`
	PageSize   int        `json:"PageSize" xml:"PageSize"`
	Total      int        `json:"Total" xml:"Total"`
	List       []ListItem `json:"List" xml:"List"`
}

// CreateDescribeShardTaskListRequest creates a request to invoke DescribeShardTaskList API
func CreateDescribeShardTaskListRequest() (request *DescribeShardTaskListRequest) {
	request = &DescribeShardTaskListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Drds", "2019-01-23", "DescribeShardTaskList", "drds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeShardTaskListResponse creates a response to parse from DescribeShardTaskList response
func CreateDescribeShardTaskListResponse() (response *DescribeShardTaskListResponse) {
	response = &DescribeShardTaskListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
