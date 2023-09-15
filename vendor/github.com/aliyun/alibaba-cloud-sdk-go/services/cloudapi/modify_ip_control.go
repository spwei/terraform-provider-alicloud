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

// ModifyIpControl invokes the cloudapi.ModifyIpControl API synchronously
func (client *Client) ModifyIpControl(request *ModifyIpControlRequest) (response *ModifyIpControlResponse, err error) {
	response = CreateModifyIpControlResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyIpControlWithChan invokes the cloudapi.ModifyIpControl API asynchronously
func (client *Client) ModifyIpControlWithChan(request *ModifyIpControlRequest) (<-chan *ModifyIpControlResponse, <-chan error) {
	responseChan := make(chan *ModifyIpControlResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyIpControl(request)
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

// ModifyIpControlWithCallback invokes the cloudapi.ModifyIpControl API asynchronously
func (client *Client) ModifyIpControlWithCallback(request *ModifyIpControlRequest, callback func(response *ModifyIpControlResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyIpControlResponse
		var err error
		defer close(result)
		response, err = client.ModifyIpControl(request)
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

// ModifyIpControlRequest is the request struct for api ModifyIpControl
type ModifyIpControlRequest struct {
	*requests.RpcRequest
	IpControlName string `position:"Query" name:"IpControlName"`
	Description   string `position:"Query" name:"Description"`
	IpControlId   string `position:"Query" name:"IpControlId"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
}

// ModifyIpControlResponse is the response struct for api ModifyIpControl
type ModifyIpControlResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyIpControlRequest creates a request to invoke ModifyIpControl API
func CreateModifyIpControlRequest() (request *ModifyIpControlRequest) {
	request = &ModifyIpControlRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "ModifyIpControl", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyIpControlResponse creates a response to parse from ModifyIpControl response
func CreateModifyIpControlResponse() (response *ModifyIpControlResponse) {
	response = &ModifyIpControlResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
