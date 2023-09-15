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

// SetDomainWebSocketStatus invokes the cloudapi.SetDomainWebSocketStatus API synchronously
func (client *Client) SetDomainWebSocketStatus(request *SetDomainWebSocketStatusRequest) (response *SetDomainWebSocketStatusResponse, err error) {
	response = CreateSetDomainWebSocketStatusResponse()
	err = client.DoAction(request, response)
	return
}

// SetDomainWebSocketStatusWithChan invokes the cloudapi.SetDomainWebSocketStatus API asynchronously
func (client *Client) SetDomainWebSocketStatusWithChan(request *SetDomainWebSocketStatusRequest) (<-chan *SetDomainWebSocketStatusResponse, <-chan error) {
	responseChan := make(chan *SetDomainWebSocketStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetDomainWebSocketStatus(request)
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

// SetDomainWebSocketStatusWithCallback invokes the cloudapi.SetDomainWebSocketStatus API asynchronously
func (client *Client) SetDomainWebSocketStatusWithCallback(request *SetDomainWebSocketStatusRequest, callback func(response *SetDomainWebSocketStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetDomainWebSocketStatusResponse
		var err error
		defer close(result)
		response, err = client.SetDomainWebSocketStatus(request)
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

// SetDomainWebSocketStatusRequest is the request struct for api SetDomainWebSocketStatus
type SetDomainWebSocketStatusRequest struct {
	*requests.RpcRequest
	WSSEnable     string `position:"Query" name:"WSSEnable"`
	GroupId       string `position:"Query" name:"GroupId"`
	DomainName    string `position:"Query" name:"DomainName"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
	ActionValue   string `position:"Query" name:"ActionValue"`
}

// SetDomainWebSocketStatusResponse is the response struct for api SetDomainWebSocketStatus
type SetDomainWebSocketStatusResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetDomainWebSocketStatusRequest creates a request to invoke SetDomainWebSocketStatus API
func CreateSetDomainWebSocketStatusRequest() (request *SetDomainWebSocketStatusRequest) {
	request = &SetDomainWebSocketStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "SetDomainWebSocketStatus", "apigateway", "openAPI")
	request.Method = requests.POST
	return
}

// CreateSetDomainWebSocketStatusResponse creates a response to parse from SetDomainWebSocketStatus response
func CreateSetDomainWebSocketStatusResponse() (response *SetDomainWebSocketStatusResponse) {
	response = &SetDomainWebSocketStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
