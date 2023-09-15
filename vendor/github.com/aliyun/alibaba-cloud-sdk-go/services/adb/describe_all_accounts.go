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

// DescribeAllAccounts invokes the adb.DescribeAllAccounts API synchronously
func (client *Client) DescribeAllAccounts(request *DescribeAllAccountsRequest) (response *DescribeAllAccountsResponse, err error) {
	response = CreateDescribeAllAccountsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeAllAccountsWithChan invokes the adb.DescribeAllAccounts API asynchronously
func (client *Client) DescribeAllAccountsWithChan(request *DescribeAllAccountsRequest) (<-chan *DescribeAllAccountsResponse, <-chan error) {
	responseChan := make(chan *DescribeAllAccountsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeAllAccounts(request)
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

// DescribeAllAccountsWithCallback invokes the adb.DescribeAllAccounts API asynchronously
func (client *Client) DescribeAllAccountsWithCallback(request *DescribeAllAccountsRequest, callback func(response *DescribeAllAccountsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeAllAccountsResponse
		var err error
		defer close(result)
		response, err = client.DescribeAllAccounts(request)
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

// DescribeAllAccountsRequest is the request struct for api DescribeAllAccounts
type DescribeAllAccountsRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeAllAccountsResponse is the response struct for api DescribeAllAccounts
type DescribeAllAccountsResponse struct {
	*responses.BaseResponse
	RequestId   string        `json:"RequestId" xml:"RequestId"`
	AccountList []AccountInfo `json:"AccountList" xml:"AccountList"`
}

// CreateDescribeAllAccountsRequest creates a request to invoke DescribeAllAccounts API
func CreateDescribeAllAccountsRequest() (request *DescribeAllAccountsRequest) {
	request = &DescribeAllAccountsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("adb", "2019-03-15", "DescribeAllAccounts", "ads", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeAllAccountsResponse creates a response to parse from DescribeAllAccounts response
func CreateDescribeAllAccountsResponse() (response *DescribeAllAccountsResponse) {
	response = &DescribeAllAccountsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
