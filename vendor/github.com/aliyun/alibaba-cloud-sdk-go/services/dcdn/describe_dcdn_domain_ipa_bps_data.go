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

// DescribeDcdnDomainIpaBpsData invokes the dcdn.DescribeDcdnDomainIpaBpsData API synchronously
func (client *Client) DescribeDcdnDomainIpaBpsData(request *DescribeDcdnDomainIpaBpsDataRequest) (response *DescribeDcdnDomainIpaBpsDataResponse, err error) {
	response = CreateDescribeDcdnDomainIpaBpsDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainIpaBpsDataWithChan invokes the dcdn.DescribeDcdnDomainIpaBpsData API asynchronously
func (client *Client) DescribeDcdnDomainIpaBpsDataWithChan(request *DescribeDcdnDomainIpaBpsDataRequest) (<-chan *DescribeDcdnDomainIpaBpsDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainIpaBpsDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainIpaBpsData(request)
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

// DescribeDcdnDomainIpaBpsDataWithCallback invokes the dcdn.DescribeDcdnDomainIpaBpsData API asynchronously
func (client *Client) DescribeDcdnDomainIpaBpsDataWithCallback(request *DescribeDcdnDomainIpaBpsDataRequest, callback func(response *DescribeDcdnDomainIpaBpsDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainIpaBpsDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainIpaBpsData(request)
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

// DescribeDcdnDomainIpaBpsDataRequest is the request struct for api DescribeDcdnDomainIpaBpsData
type DescribeDcdnDomainIpaBpsDataRequest struct {
	*requests.RpcRequest
	FixTimeGap     string `position:"Query" name:"FixTimeGap"`
	TimeMerge      string `position:"Query" name:"TimeMerge"`
	DomainName     string `position:"Query" name:"DomainName"`
	EndTime        string `position:"Query" name:"EndTime"`
	Interval       string `position:"Query" name:"Interval"`
	LocationNameEn string `position:"Query" name:"LocationNameEn"`
	StartTime      string `position:"Query" name:"StartTime"`
	IspNameEn      string `position:"Query" name:"IspNameEn"`
}

// DescribeDcdnDomainIpaBpsDataResponse is the response struct for api DescribeDcdnDomainIpaBpsData
type DescribeDcdnDomainIpaBpsDataResponse struct {
	*responses.BaseResponse
	EndTime            string                                           `json:"EndTime" xml:"EndTime"`
	StartTime          string                                           `json:"StartTime" xml:"StartTime"`
	RequestId          string                                           `json:"RequestId" xml:"RequestId"`
	DomainName         string                                           `json:"DomainName" xml:"DomainName"`
	DataInterval       string                                           `json:"DataInterval" xml:"DataInterval"`
	BpsDataPerInterval BpsDataPerIntervalInDescribeDcdnDomainIpaBpsData `json:"BpsDataPerInterval" xml:"BpsDataPerInterval"`
}

// CreateDescribeDcdnDomainIpaBpsDataRequest creates a request to invoke DescribeDcdnDomainIpaBpsData API
func CreateDescribeDcdnDomainIpaBpsDataRequest() (request *DescribeDcdnDomainIpaBpsDataRequest) {
	request = &DescribeDcdnDomainIpaBpsDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainIpaBpsData", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnDomainIpaBpsDataResponse creates a response to parse from DescribeDcdnDomainIpaBpsData response
func CreateDescribeDcdnDomainIpaBpsDataResponse() (response *DescribeDcdnDomainIpaBpsDataResponse) {
	response = &DescribeDcdnDomainIpaBpsDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
