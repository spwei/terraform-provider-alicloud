package elasticsearch

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

// ListCollectors invokes the elasticsearch.ListCollectors API synchronously
func (client *Client) ListCollectors(request *ListCollectorsRequest) (response *ListCollectorsResponse, err error) {
	response = CreateListCollectorsResponse()
	err = client.DoAction(request, response)
	return
}

// ListCollectorsWithChan invokes the elasticsearch.ListCollectors API asynchronously
func (client *Client) ListCollectorsWithChan(request *ListCollectorsRequest) (<-chan *ListCollectorsResponse, <-chan error) {
	responseChan := make(chan *ListCollectorsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListCollectors(request)
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

// ListCollectorsWithCallback invokes the elasticsearch.ListCollectors API asynchronously
func (client *Client) ListCollectorsWithCallback(request *ListCollectorsRequest, callback func(response *ListCollectorsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListCollectorsResponse
		var err error
		defer close(result)
		response, err = client.ListCollectors(request)
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

// ListCollectorsRequest is the request struct for api ListCollectors
type ListCollectorsRequest struct {
	*requests.RoaRequest
	InstanceId string           `position:"Query" name:"instanceId"`
	Size       requests.Integer `position:"Query" name:"size"`
	Name       string           `position:"Query" name:"name"`
	SourceType string           `position:"Query" name:"sourceType"`
	Page       requests.Integer `position:"Query" name:"page"`
	ResId      string           `position:"Query" name:"resId"`
}

// ListCollectorsResponse is the response struct for api ListCollectors
type ListCollectorsResponse struct {
	*responses.BaseResponse
	RequestId string       `json:"RequestId" xml:"RequestId"`
	Headers   Headers      `json:"Headers" xml:"Headers"`
	Result    []ResultItem `json:"Result" xml:"Result"`
}

// CreateListCollectorsRequest creates a request to invoke ListCollectors API
func CreateListCollectorsRequest() (request *ListCollectorsRequest) {
	request = &ListCollectorsRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("elasticsearch", "2017-06-13", "ListCollectors", "/openapi/collectors", "elasticsearch", "openAPI")
	request.Method = requests.GET
	return
}

// CreateListCollectorsResponse creates a response to parse from ListCollectors response
func CreateListCollectorsResponse() (response *ListCollectorsResponse) {
	response = &ListCollectorsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
