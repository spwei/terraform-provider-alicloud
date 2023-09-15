package cloudapi

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

// BatchAbolishApis invokes the cloudapi.BatchAbolishApis API synchronously
func (client *Client) BatchAbolishApis(request *BatchAbolishApisRequest) (response *BatchAbolishApisResponse, err error) {
	response = CreateBatchAbolishApisResponse()
	err = client.DoAction(request, response)
	return
}

// BatchAbolishApisWithChan invokes the cloudapi.BatchAbolishApis API asynchronously
func (client *Client) BatchAbolishApisWithChan(request *BatchAbolishApisRequest) (<-chan *BatchAbolishApisResponse, <-chan error) {
	responseChan := make(chan *BatchAbolishApisResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BatchAbolishApis(request)
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

// BatchAbolishApisWithCallback invokes the cloudapi.BatchAbolishApis API asynchronously
func (client *Client) BatchAbolishApisWithCallback(request *BatchAbolishApisRequest, callback func(response *BatchAbolishApisResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BatchAbolishApisResponse
		var err error
		defer close(result)
		response, err = client.BatchAbolishApis(request)
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

// BatchAbolishApisRequest is the request struct for api BatchAbolishApis
type BatchAbolishApisRequest struct {
	*requests.RpcRequest
	SecurityToken string                 `position:"Query" name:"SecurityToken"`
	Api           *[]BatchAbolishApisApi `position:"Query" name:"Api"  type:"Repeated"`
}

// BatchAbolishApisApi is a repeated param struct in BatchAbolishApisRequest
type BatchAbolishApisApi struct {
	StageName string `name:"StageName"`
	GroupId   string `name:"GroupId"`
	ApiUid    string `name:"ApiUid"`
	StageId   string `name:"StageId"`
}

// BatchAbolishApisResponse is the response struct for api BatchAbolishApis
type BatchAbolishApisResponse struct {
	*responses.BaseResponse
	OperationId string `json:"OperationId" xml:"OperationId"`
	RequestId   string `json:"RequestId" xml:"RequestId"`
}

// CreateBatchAbolishApisRequest creates a request to invoke BatchAbolishApis API
func CreateBatchAbolishApisRequest() (request *BatchAbolishApisRequest) {
	request = &BatchAbolishApisRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "BatchAbolishApis", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateBatchAbolishApisResponse creates a response to parse from BatchAbolishApis response
func CreateBatchAbolishApisResponse() (response *BatchAbolishApisResponse) {
	response = &BatchAbolishApisResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}