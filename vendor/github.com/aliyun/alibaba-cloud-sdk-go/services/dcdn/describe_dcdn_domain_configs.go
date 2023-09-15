package dcdn

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

// DescribeDcdnDomainConfigs invokes the dcdn.DescribeDcdnDomainConfigs API synchronously
func (client *Client) DescribeDcdnDomainConfigs(request *DescribeDcdnDomainConfigsRequest) (response *DescribeDcdnDomainConfigsResponse, err error) {
	response = CreateDescribeDcdnDomainConfigsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainConfigsWithChan invokes the dcdn.DescribeDcdnDomainConfigs API asynchronously
func (client *Client) DescribeDcdnDomainConfigsWithChan(request *DescribeDcdnDomainConfigsRequest) (<-chan *DescribeDcdnDomainConfigsResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainConfigsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainConfigs(request)
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

// DescribeDcdnDomainConfigsWithCallback invokes the dcdn.DescribeDcdnDomainConfigs API asynchronously
func (client *Client) DescribeDcdnDomainConfigsWithCallback(request *DescribeDcdnDomainConfigsRequest, callback func(response *DescribeDcdnDomainConfigsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainConfigsResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainConfigs(request)
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

// DescribeDcdnDomainConfigsRequest is the request struct for api DescribeDcdnDomainConfigs
type DescribeDcdnDomainConfigsRequest struct {
	*requests.RpcRequest
	FunctionNames string           `position:"Query" name:"FunctionNames"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	ConfigId      string           `position:"Query" name:"ConfigId"`
}

// DescribeDcdnDomainConfigsResponse is the response struct for api DescribeDcdnDomainConfigs
type DescribeDcdnDomainConfigsResponse struct {
	*responses.BaseResponse
	RequestId     string                                   `json:"RequestId" xml:"RequestId"`
	DomainConfigs DomainConfigsInDescribeDcdnDomainConfigs `json:"DomainConfigs" xml:"DomainConfigs"`
}

// CreateDescribeDcdnDomainConfigsRequest creates a request to invoke DescribeDcdnDomainConfigs API
func CreateDescribeDcdnDomainConfigsRequest() (request *DescribeDcdnDomainConfigsRequest) {
	request = &DescribeDcdnDomainConfigsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainConfigs", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnDomainConfigsResponse creates a response to parse from DescribeDcdnDomainConfigs response
func CreateDescribeDcdnDomainConfigsResponse() (response *DescribeDcdnDomainConfigsResponse) {
	response = &DescribeDcdnDomainConfigsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
