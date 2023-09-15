package ess

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

// ModifyScalingConfiguration invokes the ess.ModifyScalingConfiguration API synchronously
func (client *Client) ModifyScalingConfiguration(request *ModifyScalingConfigurationRequest) (response *ModifyScalingConfigurationResponse, err error) {
	response = CreateModifyScalingConfigurationResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyScalingConfigurationWithChan invokes the ess.ModifyScalingConfiguration API asynchronously
func (client *Client) ModifyScalingConfigurationWithChan(request *ModifyScalingConfigurationRequest) (<-chan *ModifyScalingConfigurationResponse, <-chan error) {
	responseChan := make(chan *ModifyScalingConfigurationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyScalingConfiguration(request)
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

// ModifyScalingConfigurationWithCallback invokes the ess.ModifyScalingConfiguration API asynchronously
func (client *Client) ModifyScalingConfigurationWithCallback(request *ModifyScalingConfigurationRequest, callback func(response *ModifyScalingConfigurationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyScalingConfigurationResponse
		var err error
		defer close(result)
		response, err = client.ModifyScalingConfiguration(request)
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

// ModifyScalingConfigurationRequest is the request struct for api ModifyScalingConfiguration
type ModifyScalingConfigurationRequest struct {
	*requests.RpcRequest
	HpcClusterId                    string                                            `position:"Query" name:"HpcClusterId"`
	KeyPairName                     string                                            `position:"Query" name:"KeyPairName"`
	SpotPriceLimit                  *[]ModifyScalingConfigurationSpotPriceLimit       `position:"Query" name:"SpotPriceLimit"  type:"Repeated"`
	ResourceGroupId                 string                                            `position:"Query" name:"ResourceGroupId"`
	PrivatePoolOptionsMatchCriteria string                                            `position:"Query" name:"PrivatePoolOptions.MatchCriteria"`
	HostName                        string                                            `position:"Query" name:"HostName"`
	Password                        string                                            `position:"Query" name:"Password"`
	InstanceDescription             string                                            `position:"Query" name:"InstanceDescription"`
	SystemDiskAutoSnapshotPolicyId  string                                            `position:"Query" name:"SystemDisk.AutoSnapshotPolicyId"`
	PrivatePoolOptionsId            string                                            `position:"Query" name:"PrivatePoolOptions.Id"`
	Ipv6AddressCount                requests.Integer                                  `position:"Query" name:"Ipv6AddressCount"`
	Cpu                             requests.Integer                                  `position:"Query" name:"Cpu"`
	SystemDiskCategories            *[]string                                         `position:"Query" name:"SystemDiskCategories"  type:"Repeated"`
	OwnerId                         requests.Integer                                  `position:"Query" name:"OwnerId"`
	ScalingConfigurationName        string                                            `position:"Query" name:"ScalingConfigurationName"`
	Tags                            string                                            `position:"Query" name:"Tags"`
	ScalingConfigurationId          string                                            `position:"Query" name:"ScalingConfigurationId"`
	SpotStrategy                    string                                            `position:"Query" name:"SpotStrategy"`
	SystemDiskBurstingEnabled       requests.Boolean                                  `position:"Query" name:"SystemDisk.BurstingEnabled"`
	InstanceName                    string                                            `position:"Query" name:"InstanceName"`
	InternetChargeType              string                                            `position:"Query" name:"InternetChargeType"`
	ZoneId                          string                                            `position:"Query" name:"ZoneId"`
	InstancePatternInfo             *[]ModifyScalingConfigurationInstancePatternInfo  `position:"Query" name:"InstancePatternInfo"  type:"Repeated"`
	Affinity                        string                                            `position:"Query" name:"Affinity"`
	ImageId                         string                                            `position:"Query" name:"ImageId"`
	Memory                          requests.Integer                                  `position:"Query" name:"Memory"`
	SpotInterruptionBehavior        string                                            `position:"Query" name:"SpotInterruptionBehavior"`
	IoOptimized                     string                                            `position:"Query" name:"IoOptimized"`
	InstanceTypes                   *[]string                                         `position:"Query" name:"InstanceTypes"  type:"Repeated"`
	InternetMaxBandwidthOut         requests.Integer                                  `position:"Query" name:"InternetMaxBandwidthOut"`
	SecurityGroupId                 string                                            `position:"Query" name:"SecurityGroupId"`
	SystemDiskKMSKeyId              string                                            `position:"Query" name:"SystemDisk.KMSKeyId"`
	SystemDiskCategory              string                                            `position:"Query" name:"SystemDisk.Category"`
	SystemDiskPerformanceLevel      string                                            `position:"Query" name:"SystemDisk.PerformanceLevel"`
	UserData                        string                                            `position:"Query" name:"UserData"`
	PasswordInherit                 requests.Boolean                                  `position:"Query" name:"PasswordInherit"`
	ImageName                       string                                            `position:"Query" name:"ImageName"`
	Override                        requests.Boolean                                  `position:"Query" name:"Override"`
	SchedulerOptions                map[string]interface{}                            `position:"Query" name:"SchedulerOptions"`
	DeploymentSetId                 string                                            `position:"Query" name:"DeploymentSetId"`
	ResourceOwnerAccount            string                                            `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount                    string                                            `position:"Query" name:"OwnerAccount"`
	Tenancy                         string                                            `position:"Query" name:"Tenancy"`
	SystemDiskDiskName              string                                            `position:"Query" name:"SystemDisk.DiskName"`
	RamRoleName                     string                                            `position:"Query" name:"RamRoleName"`
	SystemDiskEncryptAlgorithm      string                                            `position:"Query" name:"SystemDisk.EncryptAlgorithm"`
	DedicatedHostId                 string                                            `position:"Query" name:"DedicatedHostId"`
	CreditSpecification             string                                            `position:"Query" name:"CreditSpecification"`
	SpotDuration                    requests.Integer                                  `position:"Query" name:"SpotDuration"`
	SecurityGroupIds                *[]string                                         `position:"Query" name:"SecurityGroupIds"  type:"Repeated"`
	DataDisk                        *[]ModifyScalingConfigurationDataDisk             `position:"Query" name:"DataDisk"  type:"Repeated"`
	InstanceTypeOverride            *[]ModifyScalingConfigurationInstanceTypeOverride `position:"Query" name:"InstanceTypeOverride"  type:"Repeated"`
	SystemDiskProvisionedIops       requests.Integer                                  `position:"Query" name:"SystemDisk.ProvisionedIops"`
	LoadBalancerWeight              requests.Integer                                  `position:"Query" name:"LoadBalancerWeight"`
	SystemDiskSize                  requests.Integer                                  `position:"Query" name:"SystemDisk.Size"`
	ImageFamily                     string                                            `position:"Query" name:"ImageFamily"`
	SystemDiskDescription           string                                            `position:"Query" name:"SystemDisk.Description"`
	SystemDiskEncrypted             requests.Boolean                                  `position:"Query" name:"SystemDisk.Encrypted"`
}

// ModifyScalingConfigurationSpotPriceLimit is a repeated param struct in ModifyScalingConfigurationRequest
type ModifyScalingConfigurationSpotPriceLimit struct {
	InstanceType string `name:"InstanceType"`
	PriceLimit   string `name:"PriceLimit"`
}

// ModifyScalingConfigurationInstancePatternInfo is a repeated param struct in ModifyScalingConfigurationRequest
type ModifyScalingConfigurationInstancePatternInfo struct {
	Cores                string    `name:"Cores"`
	InstanceFamilyLevel  string    `name:"InstanceFamilyLevel"`
	Memory               string    `name:"Memory"`
	MaxPrice             string    `name:"MaxPrice"`
	ExcludedInstanceType *[]string `name:"ExcludedInstanceType" type:"Repeated"`
	BurstablePerformance string    `name:"BurstablePerformance"`
	Architecture         *[]string `name:"Architecture" type:"Repeated"`
}

// ModifyScalingConfigurationDataDisk is a repeated param struct in ModifyScalingConfigurationRequest
type ModifyScalingConfigurationDataDisk struct {
	SnapshotId           string    `name:"SnapshotId"`
	PerformanceLevel     string    `name:"PerformanceLevel"`
	AutoSnapshotPolicyId string    `name:"AutoSnapshotPolicyId"`
	Description          string    `name:"Description"`
	BurstingEnabled      string    `name:"BurstingEnabled"`
	DiskName             string    `name:"DiskName"`
	ProvisionedIops      string    `name:"ProvisionedIops"`
	Encrypted            string    `name:"Encrypted"`
	Size                 string    `name:"Size"`
	Categories           *[]string `name:"Categories" type:"Repeated"`
	Category             string    `name:"Category"`
	KMSKeyId             string    `name:"KMSKeyId"`
	Device               string    `name:"Device"`
	DeleteWithInstance   string    `name:"DeleteWithInstance"`
}

// ModifyScalingConfigurationInstanceTypeOverride is a repeated param struct in ModifyScalingConfigurationRequest
type ModifyScalingConfigurationInstanceTypeOverride struct {
	WeightedCapacity string `name:"WeightedCapacity"`
	InstanceType     string `name:"InstanceType"`
}

// ModifyScalingConfigurationResponse is the response struct for api ModifyScalingConfiguration
type ModifyScalingConfigurationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyScalingConfigurationRequest creates a request to invoke ModifyScalingConfiguration API
func CreateModifyScalingConfigurationRequest() (request *ModifyScalingConfigurationRequest) {
	request = &ModifyScalingConfigurationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ess", "2014-08-28", "ModifyScalingConfiguration", "ess", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyScalingConfigurationResponse creates a response to parse from ModifyScalingConfiguration response
func CreateModifyScalingConfigurationResponse() (response *ModifyScalingConfigurationResponse) {
	response = &ModifyScalingConfigurationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
