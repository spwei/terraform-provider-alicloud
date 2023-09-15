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

// DescribeDcdnDomainHitRateData invokes the dcdn.DescribeDcdnDomainHitRateData API synchronously
func (client *Client) DescribeDcdnDomainHitRateData(request *DescribeDcdnDomainHitRateDataRequest) (response *DescribeDcdnDomainHitRateDataResponse, err error) {
	response = CreateDescribeDcdnDomainHitRateDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainHitRateDataWithChan invokes the dcdn.DescribeDcdnDomainHitRateData API asynchronously
func (client *Client) DescribeDcdnDomainHitRateDataWithChan(request *DescribeDcdnDomainHitRateDataRequest) (<-chan *DescribeDcdnDomainHitRateDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainHitRateDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainHitRateData(request)
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

// DescribeDcdnDomainHitRateDataWithCallback invokes the dcdn.DescribeDcdnDomainHitRateData API asynchronously
func (client *Client) DescribeDcdnDomainHitRateDataWithCallback(request *DescribeDcdnDomainHitRateDataRequest, callback func(response *DescribeDcdnDomainHitRateDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainHitRateDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainHitRateData(request)
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

// DescribeDcdnDomainHitRateDataRequest is the request struct for api DescribeDcdnDomainHitRateData
type DescribeDcdnDomainHitRateDataRequest struct {
	*requests.RpcRequest
	DomainName string `position:"Query" name:"DomainName"`
	EndTime    string `position:"Query" name:"EndTime"`
	Interval   string `position:"Query" name:"Interval"`
	StartTime  string `position:"Query" name:"StartTime"`
}

// DescribeDcdnDomainHitRateDataResponse is the response struct for api DescribeDcdnDomainHitRateData
type DescribeDcdnDomainHitRateDataResponse struct {
	*responses.BaseResponse
	EndTime            string             `json:"EndTime" xml:"EndTime"`
	StartTime          string             `json:"StartTime" xml:"StartTime"`
	RequestId          string             `json:"RequestId" xml:"RequestId"`
	DomainName         string             `json:"DomainName" xml:"DomainName"`
	DataInterval       string             `json:"DataInterval" xml:"DataInterval"`
	HitRatePerInterval HitRatePerInterval `json:"HitRatePerInterval" xml:"HitRatePerInterval"`
}

// CreateDescribeDcdnDomainHitRateDataRequest creates a request to invoke DescribeDcdnDomainHitRateData API
func CreateDescribeDcdnDomainHitRateDataRequest() (request *DescribeDcdnDomainHitRateDataRequest) {
	request = &DescribeDcdnDomainHitRateDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainHitRateData", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnDomainHitRateDataResponse creates a response to parse from DescribeDcdnDomainHitRateData response
func CreateDescribeDcdnDomainHitRateDataResponse() (response *DescribeDcdnDomainHitRateDataResponse) {
	response = &DescribeDcdnDomainHitRateDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
