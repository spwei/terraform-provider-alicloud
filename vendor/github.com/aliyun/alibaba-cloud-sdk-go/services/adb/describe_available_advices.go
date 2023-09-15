package adb

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

// DescribeAvailableAdvices invokes the adb.DescribeAvailableAdvices API synchronously
func (client *Client) DescribeAvailableAdvices(request *DescribeAvailableAdvicesRequest) (response *DescribeAvailableAdvicesResponse, err error) {
	response = CreateDescribeAvailableAdvicesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAvailableAdvicesWithChan invokes the adb.DescribeAvailableAdvices API asynchronously
func (client *Client) DescribeAvailableAdvicesWithChan(request *DescribeAvailableAdvicesRequest) (<-chan *DescribeAvailableAdvicesResponse, <-chan error) {
	responseChan := make(chan *DescribeAvailableAdvicesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAvailableAdvices(request)
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

// DescribeAvailableAdvicesWithCallback invokes the adb.DescribeAvailableAdvices API asynchronously
func (client *Client) DescribeAvailableAdvicesWithCallback(request *DescribeAvailableAdvicesRequest, callback func(response *DescribeAvailableAdvicesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAvailableAdvicesResponse
		var err error
		defer close(result)
		response, err = client.DescribeAvailableAdvices(request)
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

// DescribeAvailableAdvicesRequest is the request struct for api DescribeAvailableAdvices
type DescribeAvailableAdvicesRequest struct {
	*requests.RpcRequest
	DBClusterId string           `position:"Query" name:"DBClusterId"`
	PageNumber  requests.Integer `position:"Query" name:"PageNumber"`
	AdviceDate  requests.Integer `position:"Query" name:"AdviceDate"`
	PageSize    requests.Integer `position:"Query" name:"PageSize"`
	Lang        string           `position:"Query" name:"Lang"`
}

// DescribeAvailableAdvicesResponse is the response struct for api DescribeAvailableAdvices
type DescribeAvailableAdvicesResponse struct {
	*responses.BaseResponse
	PageNumber int64       `json:"PageNumber" xml:"PageNumber"`
	PageSize   int64       `json:"PageSize" xml:"PageSize"`
	RequestId  string      `json:"RequestId" xml:"RequestId"`
	TotalCount int64       `json:"TotalCount" xml:"TotalCount"`
	Items      []ItemsItem `json:"Items" xml:"Items"`
}

// CreateDescribeAvailableAdvicesRequest creates a request to invoke DescribeAvailableAdvices API
func CreateDescribeAvailableAdvicesRequest() (request *DescribeAvailableAdvicesRequest) {
	request = &DescribeAvailableAdvicesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("adb", "2019-03-15", "DescribeAvailableAdvices", "ads", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAvailableAdvicesResponse creates a response to parse from DescribeAvailableAdvices response
func CreateDescribeAvailableAdvicesResponse() (response *DescribeAvailableAdvicesResponse) {
	response = &DescribeAvailableAdvicesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
