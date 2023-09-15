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

// DescribeDcdnDomainQpsDataByLayer invokes the dcdn.DescribeDcdnDomainQpsDataByLayer API synchronously
func (client *Client) DescribeDcdnDomainQpsDataByLayer(request *DescribeDcdnDomainQpsDataByLayerRequest) (response *DescribeDcdnDomainQpsDataByLayerResponse, err error) {
	response = CreateDescribeDcdnDomainQpsDataByLayerResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDcdnDomainQpsDataByLayerWithChan invokes the dcdn.DescribeDcdnDomainQpsDataByLayer API asynchronously
func (client *Client) DescribeDcdnDomainQpsDataByLayerWithChan(request *DescribeDcdnDomainQpsDataByLayerRequest) (<-chan *DescribeDcdnDomainQpsDataByLayerResponse, <-chan error) {
	responseChan := make(chan *DescribeDcdnDomainQpsDataByLayerResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDcdnDomainQpsDataByLayer(request)
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

// DescribeDcdnDomainQpsDataByLayerWithCallback invokes the dcdn.DescribeDcdnDomainQpsDataByLayer API asynchronously
func (client *Client) DescribeDcdnDomainQpsDataByLayerWithCallback(request *DescribeDcdnDomainQpsDataByLayerRequest, callback func(response *DescribeDcdnDomainQpsDataByLayerResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDcdnDomainQpsDataByLayerResponse
		var err error
		defer close(result)
		response, err = client.DescribeDcdnDomainQpsDataByLayer(request)
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

// DescribeDcdnDomainQpsDataByLayerRequest is the request struct for api DescribeDcdnDomainQpsDataByLayer
type DescribeDcdnDomainQpsDataByLayerRequest struct {
	*requests.RpcRequest
	DomainName     string `position:"Query" name:"DomainName"`
	EndTime        string `position:"Query" name:"EndTime"`
	Interval       string `position:"Query" name:"Interval"`
	LocationNameEn string `position:"Query" name:"LocationNameEn"`
	StartTime      string `position:"Query" name:"StartTime"`
	IspNameEn      string `position:"Query" name:"IspNameEn"`
	Layer          string `position:"Query" name:"Layer"`
}

// DescribeDcdnDomainQpsDataByLayerResponse is the response struct for api DescribeDcdnDomainQpsDataByLayer
type DescribeDcdnDomainQpsDataByLayerResponse struct {
	*responses.BaseResponse
	EndTime         string          `json:"EndTime" xml:"EndTime"`
	StartTime       string          `json:"StartTime" xml:"StartTime"`
	RequestId       string          `json:"RequestId" xml:"RequestId"`
	Layer           string          `json:"Layer" xml:"Layer"`
	DomainName      string          `json:"DomainName" xml:"DomainName"`
	DataInterval    string          `json:"DataInterval" xml:"DataInterval"`
	QpsDataInterval QpsDataInterval `json:"QpsDataInterval" xml:"QpsDataInterval"`
}

// CreateDescribeDcdnDomainQpsDataByLayerRequest creates a request to invoke DescribeDcdnDomainQpsDataByLayer API
func CreateDescribeDcdnDomainQpsDataByLayerRequest() (request *DescribeDcdnDomainQpsDataByLayerRequest) {
	request = &DescribeDcdnDomainQpsDataByLayerRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("dcdn", "2018-01-15", "DescribeDcdnDomainQpsDataByLayer", "", "")
	request.Method = requests.POST
	return
}

// CreateDescribeDcdnDomainQpsDataByLayerResponse creates a response to parse from DescribeDcdnDomainQpsDataByLayer response
func CreateDescribeDcdnDomainQpsDataByLayerResponse() (response *DescribeDcdnDomainQpsDataByLayerResponse) {
	response = &DescribeDcdnDomainQpsDataByLayerResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
