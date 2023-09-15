package bssopenapi

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

// InquiryPriceRefundInstance invokes the bssopenapi.InquiryPriceRefundInstance API synchronously
func (client *Client) InquiryPriceRefundInstance(request *InquiryPriceRefundInstanceRequest) (response *InquiryPriceRefundInstanceResponse, err error) {
	response = CreateInquiryPriceRefundInstanceResponse()
	err = client.DoAction(request, response)
	return
}

// InquiryPriceRefundInstanceWithChan invokes the bssopenapi.InquiryPriceRefundInstance API asynchronously
func (client *Client) InquiryPriceRefundInstanceWithChan(request *InquiryPriceRefundInstanceRequest) (<-chan *InquiryPriceRefundInstanceResponse, <-chan error) {
	responseChan := make(chan *InquiryPriceRefundInstanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.InquiryPriceRefundInstance(request)
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

// InquiryPriceRefundInstanceWithCallback invokes the bssopenapi.InquiryPriceRefundInstance API asynchronously
func (client *Client) InquiryPriceRefundInstanceWithCallback(request *InquiryPriceRefundInstanceRequest, callback func(response *InquiryPriceRefundInstanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *InquiryPriceRefundInstanceResponse
		var err error
		defer close(result)
		response, err = client.InquiryPriceRefundInstance(request)
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

// InquiryPriceRefundInstanceRequest is the request struct for api InquiryPriceRefundInstance
type InquiryPriceRefundInstanceRequest struct {
	*requests.RpcRequest
	ProductCode string `position:"Query" name:"ProductCode"`
	ClientToken string `position:"Query" name:"ClientToken"`
	ProductType string `position:"Query" name:"ProductType"`
	InstanceId  string `position:"Query" name:"InstanceId"`
}

// InquiryPriceRefundInstanceResponse is the response struct for api InquiryPriceRefundInstance
type InquiryPriceRefundInstanceResponse struct {
	*responses.BaseResponse
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Success   bool   `json:"Success" xml:"Success"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateInquiryPriceRefundInstanceRequest creates a request to invoke InquiryPriceRefundInstance API
func CreateInquiryPriceRefundInstanceRequest() (request *InquiryPriceRefundInstanceRequest) {
	request = &InquiryPriceRefundInstanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "InquiryPriceRefundInstance", "bssopenapi", "openAPI")
	request.Method = requests.POST
	return
}

// CreateInquiryPriceRefundInstanceResponse creates a response to parse from InquiryPriceRefundInstance response
func CreateInquiryPriceRefundInstanceResponse() (response *InquiryPriceRefundInstanceResponse) {
	response = &InquiryPriceRefundInstanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}