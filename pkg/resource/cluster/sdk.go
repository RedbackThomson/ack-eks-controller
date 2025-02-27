// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package cluster

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/eks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/eks-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.EKS{}
	_ = &svcapitypes.Cluster{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeClusterOutput
	resp, err = rm.sdkapi.DescribeClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeCluster", err)
	if err != nil {
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
		}
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "ResourceNotFoundException" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Cluster.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Cluster.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Cluster.CertificateAuthority != nil {
		f1 := &svcapitypes.Certificate{}
		if resp.Cluster.CertificateAuthority.Data != nil {
			f1.Data = resp.Cluster.CertificateAuthority.Data
		}
		ko.Status.CertificateAuthority = f1
	} else {
		ko.Status.CertificateAuthority = nil
	}
	if resp.Cluster.ClientRequestToken != nil {
		ko.Spec.ClientRequestToken = resp.Cluster.ClientRequestToken
	} else {
		ko.Spec.ClientRequestToken = nil
	}
	if resp.Cluster.ConnectorConfig != nil {
		f3 := &svcapitypes.ConnectorConfigResponse{}
		if resp.Cluster.ConnectorConfig.ActivationCode != nil {
			f3.ActivationCode = resp.Cluster.ConnectorConfig.ActivationCode
		}
		if resp.Cluster.ConnectorConfig.ActivationExpiry != nil {
			f3.ActivationExpiry = &metav1.Time{*resp.Cluster.ConnectorConfig.ActivationExpiry}
		}
		if resp.Cluster.ConnectorConfig.ActivationId != nil {
			f3.ActivationID = resp.Cluster.ConnectorConfig.ActivationId
		}
		if resp.Cluster.ConnectorConfig.Provider != nil {
			f3.Provider = resp.Cluster.ConnectorConfig.Provider
		}
		if resp.Cluster.ConnectorConfig.RoleArn != nil {
			f3.RoleARN = resp.Cluster.ConnectorConfig.RoleArn
		}
		ko.Status.ConnectorConfig = f3
	} else {
		ko.Status.ConnectorConfig = nil
	}
	if resp.Cluster.CreatedAt != nil {
		ko.Status.CreatedAt = &metav1.Time{*resp.Cluster.CreatedAt}
	} else {
		ko.Status.CreatedAt = nil
	}
	if resp.Cluster.EncryptionConfig != nil {
		f5 := []*svcapitypes.EncryptionConfig{}
		for _, f5iter := range resp.Cluster.EncryptionConfig {
			f5elem := &svcapitypes.EncryptionConfig{}
			if f5iter.Provider != nil {
				f5elemf0 := &svcapitypes.Provider{}
				if f5iter.Provider.KeyArn != nil {
					f5elemf0.KeyARN = f5iter.Provider.KeyArn
				}
				f5elem.Provider = f5elemf0
			}
			if f5iter.Resources != nil {
				f5elemf1 := []*string{}
				for _, f5elemf1iter := range f5iter.Resources {
					var f5elemf1elem string
					f5elemf1elem = *f5elemf1iter
					f5elemf1 = append(f5elemf1, &f5elemf1elem)
				}
				f5elem.Resources = f5elemf1
			}
			f5 = append(f5, f5elem)
		}
		ko.Spec.EncryptionConfig = f5
	} else {
		ko.Spec.EncryptionConfig = nil
	}
	if resp.Cluster.Endpoint != nil {
		ko.Status.Endpoint = resp.Cluster.Endpoint
	} else {
		ko.Status.Endpoint = nil
	}
	if resp.Cluster.Health != nil {
		f7 := &svcapitypes.ClusterHealth{}
		if resp.Cluster.Health.Issues != nil {
			f7f0 := []*svcapitypes.ClusterIssue{}
			for _, f7f0iter := range resp.Cluster.Health.Issues {
				f7f0elem := &svcapitypes.ClusterIssue{}
				if f7f0iter.Code != nil {
					f7f0elem.Code = f7f0iter.Code
				}
				if f7f0iter.Message != nil {
					f7f0elem.Message = f7f0iter.Message
				}
				if f7f0iter.ResourceIds != nil {
					f7f0elemf2 := []*string{}
					for _, f7f0elemf2iter := range f7f0iter.ResourceIds {
						var f7f0elemf2elem string
						f7f0elemf2elem = *f7f0elemf2iter
						f7f0elemf2 = append(f7f0elemf2, &f7f0elemf2elem)
					}
					f7f0elem.ResourceIDs = f7f0elemf2
				}
				f7f0 = append(f7f0, f7f0elem)
			}
			f7.Issues = f7f0
		}
		ko.Status.Health = f7
	} else {
		ko.Status.Health = nil
	}
	if resp.Cluster.Id != nil {
		ko.Status.ID = resp.Cluster.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.Cluster.Identity != nil {
		f9 := &svcapitypes.Identity{}
		if resp.Cluster.Identity.Oidc != nil {
			f9f0 := &svcapitypes.OIDC{}
			if resp.Cluster.Identity.Oidc.Issuer != nil {
				f9f0.Issuer = resp.Cluster.Identity.Oidc.Issuer
			}
			f9.OIDC = f9f0
		}
		ko.Status.Identity = f9
	} else {
		ko.Status.Identity = nil
	}
	if resp.Cluster.KubernetesNetworkConfig != nil {
		f10 := &svcapitypes.KubernetesNetworkConfigRequest{}
		if resp.Cluster.KubernetesNetworkConfig.IpFamily != nil {
			f10.IPFamily = resp.Cluster.KubernetesNetworkConfig.IpFamily
		}
		if resp.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr != nil {
			f10.ServiceIPv4CIDR = resp.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr
		}
		ko.Spec.KubernetesNetworkConfig = f10
	} else {
		ko.Spec.KubernetesNetworkConfig = nil
	}
	if resp.Cluster.Logging != nil {
		f11 := &svcapitypes.Logging{}
		if resp.Cluster.Logging.ClusterLogging != nil {
			f11f0 := []*svcapitypes.LogSetup{}
			for _, f11f0iter := range resp.Cluster.Logging.ClusterLogging {
				f11f0elem := &svcapitypes.LogSetup{}
				if f11f0iter.Enabled != nil {
					f11f0elem.Enabled = f11f0iter.Enabled
				}
				if f11f0iter.Types != nil {
					f11f0elemf1 := []*string{}
					for _, f11f0elemf1iter := range f11f0iter.Types {
						var f11f0elemf1elem string
						f11f0elemf1elem = *f11f0elemf1iter
						f11f0elemf1 = append(f11f0elemf1, &f11f0elemf1elem)
					}
					f11f0elem.Types = f11f0elemf1
				}
				f11f0 = append(f11f0, f11f0elem)
			}
			f11.ClusterLogging = f11f0
		}
		ko.Spec.Logging = f11
	} else {
		ko.Spec.Logging = nil
	}
	if resp.Cluster.Name != nil {
		ko.Spec.Name = resp.Cluster.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Cluster.OutpostConfig != nil {
		f13 := &svcapitypes.OutpostConfigRequest{}
		if resp.Cluster.OutpostConfig.ControlPlaneInstanceType != nil {
			f13.ControlPlaneInstanceType = resp.Cluster.OutpostConfig.ControlPlaneInstanceType
		}
		if resp.Cluster.OutpostConfig.ControlPlanePlacement != nil {
			f13f1 := &svcapitypes.ControlPlanePlacementRequest{}
			if resp.Cluster.OutpostConfig.ControlPlanePlacement.GroupName != nil {
				f13f1.GroupName = resp.Cluster.OutpostConfig.ControlPlanePlacement.GroupName
			}
			f13.ControlPlanePlacement = f13f1
		}
		if resp.Cluster.OutpostConfig.OutpostArns != nil {
			f13f2 := []*string{}
			for _, f13f2iter := range resp.Cluster.OutpostConfig.OutpostArns {
				var f13f2elem string
				f13f2elem = *f13f2iter
				f13f2 = append(f13f2, &f13f2elem)
			}
			f13.OutpostARNs = f13f2
		}
		ko.Spec.OutpostConfig = f13
	} else {
		ko.Spec.OutpostConfig = nil
	}
	if resp.Cluster.PlatformVersion != nil {
		ko.Status.PlatformVersion = resp.Cluster.PlatformVersion
	} else {
		ko.Status.PlatformVersion = nil
	}
	if resp.Cluster.ResourcesVpcConfig != nil {
		f15 := &svcapitypes.VPCConfigRequest{}
		if resp.Cluster.ResourcesVpcConfig.EndpointPrivateAccess != nil {
			f15.EndpointPrivateAccess = resp.Cluster.ResourcesVpcConfig.EndpointPrivateAccess
		}
		if resp.Cluster.ResourcesVpcConfig.EndpointPublicAccess != nil {
			f15.EndpointPublicAccess = resp.Cluster.ResourcesVpcConfig.EndpointPublicAccess
		}
		if resp.Cluster.ResourcesVpcConfig.PublicAccessCidrs != nil {
			f15f3 := []*string{}
			for _, f15f3iter := range resp.Cluster.ResourcesVpcConfig.PublicAccessCidrs {
				var f15f3elem string
				f15f3elem = *f15f3iter
				f15f3 = append(f15f3, &f15f3elem)
			}
			f15.PublicAccessCIDRs = f15f3
		}
		if resp.Cluster.ResourcesVpcConfig.SecurityGroupIds != nil {
			f15f4 := []*string{}
			for _, f15f4iter := range resp.Cluster.ResourcesVpcConfig.SecurityGroupIds {
				var f15f4elem string
				f15f4elem = *f15f4iter
				f15f4 = append(f15f4, &f15f4elem)
			}
			f15.SecurityGroupIDs = f15f4
		}
		if resp.Cluster.ResourcesVpcConfig.SubnetIds != nil {
			f15f5 := []*string{}
			for _, f15f5iter := range resp.Cluster.ResourcesVpcConfig.SubnetIds {
				var f15f5elem string
				f15f5elem = *f15f5iter
				f15f5 = append(f15f5, &f15f5elem)
			}
			f15.SubnetIDs = f15f5
		}
		ko.Spec.ResourcesVPCConfig = f15
	} else {
		ko.Spec.ResourcesVPCConfig = nil
	}
	if resp.Cluster.RoleArn != nil {
		ko.Spec.RoleARN = resp.Cluster.RoleArn
	} else {
		ko.Spec.RoleARN = nil
	}
	if resp.Cluster.Status != nil {
		ko.Status.Status = resp.Cluster.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.Cluster.Tags != nil {
		f18 := map[string]*string{}
		for f18key, f18valiter := range resp.Cluster.Tags {
			var f18val string
			f18val = *f18valiter
			f18[f18key] = &f18val
		}
		ko.Spec.Tags = f18
	} else {
		ko.Spec.Tags = nil
	}
	if resp.Cluster.Version != nil {
		ko.Spec.Version = resp.Cluster.Version
	} else {
		ko.Spec.Version = nil
	}

	rm.setStatusDefaults(ko)
	if !clusterActive(&resource{ko}) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource. No need to return a requeue error here.
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
	} else {
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionTrue, nil, nil)
	}

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeClusterInput, error) {
	res := &svcsdk.DescribeClusterInput{}

	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateClusterOutput
	_ = resp
	resp, err = rm.sdkapi.CreateClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateCluster", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Cluster.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Cluster.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Cluster.CertificateAuthority != nil {
		f1 := &svcapitypes.Certificate{}
		if resp.Cluster.CertificateAuthority.Data != nil {
			f1.Data = resp.Cluster.CertificateAuthority.Data
		}
		ko.Status.CertificateAuthority = f1
	} else {
		ko.Status.CertificateAuthority = nil
	}
	if resp.Cluster.ClientRequestToken != nil {
		ko.Spec.ClientRequestToken = resp.Cluster.ClientRequestToken
	} else {
		ko.Spec.ClientRequestToken = nil
	}
	if resp.Cluster.ConnectorConfig != nil {
		f3 := &svcapitypes.ConnectorConfigResponse{}
		if resp.Cluster.ConnectorConfig.ActivationCode != nil {
			f3.ActivationCode = resp.Cluster.ConnectorConfig.ActivationCode
		}
		if resp.Cluster.ConnectorConfig.ActivationExpiry != nil {
			f3.ActivationExpiry = &metav1.Time{*resp.Cluster.ConnectorConfig.ActivationExpiry}
		}
		if resp.Cluster.ConnectorConfig.ActivationId != nil {
			f3.ActivationID = resp.Cluster.ConnectorConfig.ActivationId
		}
		if resp.Cluster.ConnectorConfig.Provider != nil {
			f3.Provider = resp.Cluster.ConnectorConfig.Provider
		}
		if resp.Cluster.ConnectorConfig.RoleArn != nil {
			f3.RoleARN = resp.Cluster.ConnectorConfig.RoleArn
		}
		ko.Status.ConnectorConfig = f3
	} else {
		ko.Status.ConnectorConfig = nil
	}
	if resp.Cluster.CreatedAt != nil {
		ko.Status.CreatedAt = &metav1.Time{*resp.Cluster.CreatedAt}
	} else {
		ko.Status.CreatedAt = nil
	}
	if resp.Cluster.EncryptionConfig != nil {
		f5 := []*svcapitypes.EncryptionConfig{}
		for _, f5iter := range resp.Cluster.EncryptionConfig {
			f5elem := &svcapitypes.EncryptionConfig{}
			if f5iter.Provider != nil {
				f5elemf0 := &svcapitypes.Provider{}
				if f5iter.Provider.KeyArn != nil {
					f5elemf0.KeyARN = f5iter.Provider.KeyArn
				}
				f5elem.Provider = f5elemf0
			}
			if f5iter.Resources != nil {
				f5elemf1 := []*string{}
				for _, f5elemf1iter := range f5iter.Resources {
					var f5elemf1elem string
					f5elemf1elem = *f5elemf1iter
					f5elemf1 = append(f5elemf1, &f5elemf1elem)
				}
				f5elem.Resources = f5elemf1
			}
			f5 = append(f5, f5elem)
		}
		ko.Spec.EncryptionConfig = f5
	} else {
		ko.Spec.EncryptionConfig = nil
	}
	if resp.Cluster.Endpoint != nil {
		ko.Status.Endpoint = resp.Cluster.Endpoint
	} else {
		ko.Status.Endpoint = nil
	}
	if resp.Cluster.Health != nil {
		f7 := &svcapitypes.ClusterHealth{}
		if resp.Cluster.Health.Issues != nil {
			f7f0 := []*svcapitypes.ClusterIssue{}
			for _, f7f0iter := range resp.Cluster.Health.Issues {
				f7f0elem := &svcapitypes.ClusterIssue{}
				if f7f0iter.Code != nil {
					f7f0elem.Code = f7f0iter.Code
				}
				if f7f0iter.Message != nil {
					f7f0elem.Message = f7f0iter.Message
				}
				if f7f0iter.ResourceIds != nil {
					f7f0elemf2 := []*string{}
					for _, f7f0elemf2iter := range f7f0iter.ResourceIds {
						var f7f0elemf2elem string
						f7f0elemf2elem = *f7f0elemf2iter
						f7f0elemf2 = append(f7f0elemf2, &f7f0elemf2elem)
					}
					f7f0elem.ResourceIDs = f7f0elemf2
				}
				f7f0 = append(f7f0, f7f0elem)
			}
			f7.Issues = f7f0
		}
		ko.Status.Health = f7
	} else {
		ko.Status.Health = nil
	}
	if resp.Cluster.Id != nil {
		ko.Status.ID = resp.Cluster.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.Cluster.Identity != nil {
		f9 := &svcapitypes.Identity{}
		if resp.Cluster.Identity.Oidc != nil {
			f9f0 := &svcapitypes.OIDC{}
			if resp.Cluster.Identity.Oidc.Issuer != nil {
				f9f0.Issuer = resp.Cluster.Identity.Oidc.Issuer
			}
			f9.OIDC = f9f0
		}
		ko.Status.Identity = f9
	} else {
		ko.Status.Identity = nil
	}
	if resp.Cluster.KubernetesNetworkConfig != nil {
		f10 := &svcapitypes.KubernetesNetworkConfigRequest{}
		if resp.Cluster.KubernetesNetworkConfig.IpFamily != nil {
			f10.IPFamily = resp.Cluster.KubernetesNetworkConfig.IpFamily
		}
		if resp.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr != nil {
			f10.ServiceIPv4CIDR = resp.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr
		}
		ko.Spec.KubernetesNetworkConfig = f10
	} else {
		ko.Spec.KubernetesNetworkConfig = nil
	}
	if resp.Cluster.Logging != nil {
		f11 := &svcapitypes.Logging{}
		if resp.Cluster.Logging.ClusterLogging != nil {
			f11f0 := []*svcapitypes.LogSetup{}
			for _, f11f0iter := range resp.Cluster.Logging.ClusterLogging {
				f11f0elem := &svcapitypes.LogSetup{}
				if f11f0iter.Enabled != nil {
					f11f0elem.Enabled = f11f0iter.Enabled
				}
				if f11f0iter.Types != nil {
					f11f0elemf1 := []*string{}
					for _, f11f0elemf1iter := range f11f0iter.Types {
						var f11f0elemf1elem string
						f11f0elemf1elem = *f11f0elemf1iter
						f11f0elemf1 = append(f11f0elemf1, &f11f0elemf1elem)
					}
					f11f0elem.Types = f11f0elemf1
				}
				f11f0 = append(f11f0, f11f0elem)
			}
			f11.ClusterLogging = f11f0
		}
		ko.Spec.Logging = f11
	} else {
		ko.Spec.Logging = nil
	}
	if resp.Cluster.Name != nil {
		ko.Spec.Name = resp.Cluster.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Cluster.OutpostConfig != nil {
		f13 := &svcapitypes.OutpostConfigRequest{}
		if resp.Cluster.OutpostConfig.ControlPlaneInstanceType != nil {
			f13.ControlPlaneInstanceType = resp.Cluster.OutpostConfig.ControlPlaneInstanceType
		}
		if resp.Cluster.OutpostConfig.ControlPlanePlacement != nil {
			f13f1 := &svcapitypes.ControlPlanePlacementRequest{}
			if resp.Cluster.OutpostConfig.ControlPlanePlacement.GroupName != nil {
				f13f1.GroupName = resp.Cluster.OutpostConfig.ControlPlanePlacement.GroupName
			}
			f13.ControlPlanePlacement = f13f1
		}
		if resp.Cluster.OutpostConfig.OutpostArns != nil {
			f13f2 := []*string{}
			for _, f13f2iter := range resp.Cluster.OutpostConfig.OutpostArns {
				var f13f2elem string
				f13f2elem = *f13f2iter
				f13f2 = append(f13f2, &f13f2elem)
			}
			f13.OutpostARNs = f13f2
		}
		ko.Spec.OutpostConfig = f13
	} else {
		ko.Spec.OutpostConfig = nil
	}
	if resp.Cluster.PlatformVersion != nil {
		ko.Status.PlatformVersion = resp.Cluster.PlatformVersion
	} else {
		ko.Status.PlatformVersion = nil
	}
	if resp.Cluster.ResourcesVpcConfig != nil {
		f15 := &svcapitypes.VPCConfigRequest{}
		if resp.Cluster.ResourcesVpcConfig.EndpointPrivateAccess != nil {
			f15.EndpointPrivateAccess = resp.Cluster.ResourcesVpcConfig.EndpointPrivateAccess
		}
		if resp.Cluster.ResourcesVpcConfig.EndpointPublicAccess != nil {
			f15.EndpointPublicAccess = resp.Cluster.ResourcesVpcConfig.EndpointPublicAccess
		}
		if resp.Cluster.ResourcesVpcConfig.PublicAccessCidrs != nil {
			f15f3 := []*string{}
			for _, f15f3iter := range resp.Cluster.ResourcesVpcConfig.PublicAccessCidrs {
				var f15f3elem string
				f15f3elem = *f15f3iter
				f15f3 = append(f15f3, &f15f3elem)
			}
			f15.PublicAccessCIDRs = f15f3
		}
		if resp.Cluster.ResourcesVpcConfig.SecurityGroupIds != nil {
			f15f4 := []*string{}
			for _, f15f4iter := range resp.Cluster.ResourcesVpcConfig.SecurityGroupIds {
				var f15f4elem string
				f15f4elem = *f15f4iter
				f15f4 = append(f15f4, &f15f4elem)
			}
			f15.SecurityGroupIDs = f15f4
		}
		if resp.Cluster.ResourcesVpcConfig.SubnetIds != nil {
			f15f5 := []*string{}
			for _, f15f5iter := range resp.Cluster.ResourcesVpcConfig.SubnetIds {
				var f15f5elem string
				f15f5elem = *f15f5iter
				f15f5 = append(f15f5, &f15f5elem)
			}
			f15.SubnetIDs = f15f5
		}
		ko.Spec.ResourcesVPCConfig = f15
	} else {
		ko.Spec.ResourcesVPCConfig = nil
	}
	if resp.Cluster.RoleArn != nil {
		ko.Spec.RoleARN = resp.Cluster.RoleArn
	} else {
		ko.Spec.RoleARN = nil
	}
	if resp.Cluster.Status != nil {
		ko.Status.Status = resp.Cluster.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.Cluster.Tags != nil {
		f18 := map[string]*string{}
		for f18key, f18valiter := range resp.Cluster.Tags {
			var f18val string
			f18val = *f18valiter
			f18[f18key] = &f18val
		}
		ko.Spec.Tags = f18
	} else {
		ko.Spec.Tags = nil
	}
	if resp.Cluster.Version != nil {
		ko.Spec.Version = resp.Cluster.Version
	} else {
		ko.Spec.Version = nil
	}

	rm.setStatusDefaults(ko)
	// We expect the cluster to be in 'CREATING' status since we just issued
	// the call to create it, but I suppose it doesn't hurt to check here.
	if clusterCreating(&resource{ko}) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource. No need to return a requeue error here.
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
		return &resource{ko}, nil
	}

	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateClusterInput, error) {
	res := &svcsdk.CreateClusterInput{}

	if r.ko.Spec.ClientRequestToken != nil {
		res.SetClientRequestToken(*r.ko.Spec.ClientRequestToken)
	}
	if r.ko.Spec.EncryptionConfig != nil {
		f1 := []*svcsdk.EncryptionConfig{}
		for _, f1iter := range r.ko.Spec.EncryptionConfig {
			f1elem := &svcsdk.EncryptionConfig{}
			if f1iter.Provider != nil {
				f1elemf0 := &svcsdk.Provider{}
				if f1iter.Provider.KeyARN != nil {
					f1elemf0.SetKeyArn(*f1iter.Provider.KeyARN)
				}
				f1elem.SetProvider(f1elemf0)
			}
			if f1iter.Resources != nil {
				f1elemf1 := []*string{}
				for _, f1elemf1iter := range f1iter.Resources {
					var f1elemf1elem string
					f1elemf1elem = *f1elemf1iter
					f1elemf1 = append(f1elemf1, &f1elemf1elem)
				}
				f1elem.SetResources(f1elemf1)
			}
			f1 = append(f1, f1elem)
		}
		res.SetEncryptionConfig(f1)
	}
	if r.ko.Spec.KubernetesNetworkConfig != nil {
		f2 := &svcsdk.KubernetesNetworkConfigRequest{}
		if r.ko.Spec.KubernetesNetworkConfig.IPFamily != nil {
			f2.SetIpFamily(*r.ko.Spec.KubernetesNetworkConfig.IPFamily)
		}
		if r.ko.Spec.KubernetesNetworkConfig.ServiceIPv4CIDR != nil {
			f2.SetServiceIpv4Cidr(*r.ko.Spec.KubernetesNetworkConfig.ServiceIPv4CIDR)
		}
		res.SetKubernetesNetworkConfig(f2)
	}
	if r.ko.Spec.Logging != nil {
		f3 := &svcsdk.Logging{}
		if r.ko.Spec.Logging.ClusterLogging != nil {
			f3f0 := []*svcsdk.LogSetup{}
			for _, f3f0iter := range r.ko.Spec.Logging.ClusterLogging {
				f3f0elem := &svcsdk.LogSetup{}
				if f3f0iter.Enabled != nil {
					f3f0elem.SetEnabled(*f3f0iter.Enabled)
				}
				if f3f0iter.Types != nil {
					f3f0elemf1 := []*string{}
					for _, f3f0elemf1iter := range f3f0iter.Types {
						var f3f0elemf1elem string
						f3f0elemf1elem = *f3f0elemf1iter
						f3f0elemf1 = append(f3f0elemf1, &f3f0elemf1elem)
					}
					f3f0elem.SetTypes(f3f0elemf1)
				}
				f3f0 = append(f3f0, f3f0elem)
			}
			f3.SetClusterLogging(f3f0)
		}
		res.SetLogging(f3)
	}
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.OutpostConfig != nil {
		f5 := &svcsdk.OutpostConfigRequest{}
		if r.ko.Spec.OutpostConfig.ControlPlaneInstanceType != nil {
			f5.SetControlPlaneInstanceType(*r.ko.Spec.OutpostConfig.ControlPlaneInstanceType)
		}
		if r.ko.Spec.OutpostConfig.ControlPlanePlacement != nil {
			f5f1 := &svcsdk.ControlPlanePlacementRequest{}
			if r.ko.Spec.OutpostConfig.ControlPlanePlacement.GroupName != nil {
				f5f1.SetGroupName(*r.ko.Spec.OutpostConfig.ControlPlanePlacement.GroupName)
			}
			f5.SetControlPlanePlacement(f5f1)
		}
		if r.ko.Spec.OutpostConfig.OutpostARNs != nil {
			f5f2 := []*string{}
			for _, f5f2iter := range r.ko.Spec.OutpostConfig.OutpostARNs {
				var f5f2elem string
				f5f2elem = *f5f2iter
				f5f2 = append(f5f2, &f5f2elem)
			}
			f5.SetOutpostArns(f5f2)
		}
		res.SetOutpostConfig(f5)
	}
	if r.ko.Spec.ResourcesVPCConfig != nil {
		f6 := &svcsdk.VpcConfigRequest{}
		if r.ko.Spec.ResourcesVPCConfig.EndpointPrivateAccess != nil {
			f6.SetEndpointPrivateAccess(*r.ko.Spec.ResourcesVPCConfig.EndpointPrivateAccess)
		}
		if r.ko.Spec.ResourcesVPCConfig.EndpointPublicAccess != nil {
			f6.SetEndpointPublicAccess(*r.ko.Spec.ResourcesVPCConfig.EndpointPublicAccess)
		}
		if r.ko.Spec.ResourcesVPCConfig.PublicAccessCIDRs != nil {
			f6f2 := []*string{}
			for _, f6f2iter := range r.ko.Spec.ResourcesVPCConfig.PublicAccessCIDRs {
				var f6f2elem string
				f6f2elem = *f6f2iter
				f6f2 = append(f6f2, &f6f2elem)
			}
			f6.SetPublicAccessCidrs(f6f2)
		}
		if r.ko.Spec.ResourcesVPCConfig.SecurityGroupIDs != nil {
			f6f3 := []*string{}
			for _, f6f3iter := range r.ko.Spec.ResourcesVPCConfig.SecurityGroupIDs {
				var f6f3elem string
				f6f3elem = *f6f3iter
				f6f3 = append(f6f3, &f6f3elem)
			}
			f6.SetSecurityGroupIds(f6f3)
		}
		if r.ko.Spec.ResourcesVPCConfig.SubnetIDs != nil {
			f6f4 := []*string{}
			for _, f6f4iter := range r.ko.Spec.ResourcesVPCConfig.SubnetIDs {
				var f6f4elem string
				f6f4elem = *f6f4iter
				f6f4 = append(f6f4, &f6f4elem)
			}
			f6.SetSubnetIds(f6f4)
		}
		res.SetResourcesVpcConfig(f6)
	}
	if r.ko.Spec.RoleARN != nil {
		res.SetRoleArn(*r.ko.Spec.RoleARN)
	}
	if r.ko.Spec.Tags != nil {
		f8 := map[string]*string{}
		for f8key, f8valiter := range r.ko.Spec.Tags {
			var f8val string
			f8val = *f8valiter
			f8[f8key] = &f8val
		}
		res.SetTags(f8)
	}
	if r.ko.Spec.Version != nil {
		res.SetVersion(*r.ko.Spec.Version)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	return rm.customUpdate(ctx, desired, latest, delta)
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	if clusterDeleting(r) {
		return r, requeueWaitWhileDeleting
	}
	inUse, err := rm.clusterInUse(ctx, r)
	if err != nil {
		return nil, err
	} else if inUse {
		return r, requeueWaitWhileInUse
	}

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteClusterOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteClusterWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteCluster", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteClusterInput, error) {
	res := &svcsdk.DeleteClusterInput{}

	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Cluster,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "ResourceLimitExceeded",
		"ResourceNotFound",
		"ResourceInUse",
		"OptInRequired",
		"InvalidParameterCombination",
		"InvalidParameterValue",
		"InvalidParameterException",
		"InvalidQueryParameter",
		"MalformedQueryString",
		"MissingAction",
		"MissingParameter",
		"ValidationError":
		return true
	default:
		return false
	}
}

// newLogging returns a Logging object
// with each the field set by the resource's corresponding spec field.
func (rm *resourceManager) newLogging(
	r *resource,
) *svcsdk.Logging {
	res := &svcsdk.Logging{}

	if r.ko.Spec.Logging.ClusterLogging != nil {
		resf0 := []*svcsdk.LogSetup{}
		for _, resf0iter := range r.ko.Spec.Logging.ClusterLogging {
			resf0elem := &svcsdk.LogSetup{}
			if resf0iter.Enabled != nil {
				resf0elem.SetEnabled(*resf0iter.Enabled)
			}
			if resf0iter.Types != nil {
				resf0elemf1 := []*string{}
				for _, resf0elemf1iter := range resf0iter.Types {
					var resf0elemf1elem string
					resf0elemf1elem = *resf0elemf1iter
					resf0elemf1 = append(resf0elemf1, &resf0elemf1elem)
				}
				resf0elem.SetTypes(resf0elemf1)
			}
			resf0 = append(resf0, resf0elem)
		}
		res.SetClusterLogging(resf0)
	}

	return res
}

// newVpcConfigRequest returns a VpcConfigRequest object
// with each the field set by the resource's corresponding spec field.
func (rm *resourceManager) newVpcConfigRequest(
	r *resource,
) *svcsdk.VpcConfigRequest {
	res := &svcsdk.VpcConfigRequest{}

	if r.ko.Spec.ResourcesVPCConfig.EndpointPrivateAccess != nil {
		res.SetEndpointPrivateAccess(*r.ko.Spec.ResourcesVPCConfig.EndpointPrivateAccess)
	}
	if r.ko.Spec.ResourcesVPCConfig.EndpointPublicAccess != nil {
		res.SetEndpointPublicAccess(*r.ko.Spec.ResourcesVPCConfig.EndpointPublicAccess)
	}
	if r.ko.Spec.ResourcesVPCConfig.PublicAccessCIDRs != nil {
		resf2 := []*string{}
		for _, resf2iter := range r.ko.Spec.ResourcesVPCConfig.PublicAccessCIDRs {
			var resf2elem string
			resf2elem = *resf2iter
			resf2 = append(resf2, &resf2elem)
		}
		res.SetPublicAccessCidrs(resf2)
	}
	if r.ko.Spec.ResourcesVPCConfig.SecurityGroupIDs != nil {
		resf3 := []*string{}
		for _, resf3iter := range r.ko.Spec.ResourcesVPCConfig.SecurityGroupIDs {
			var resf3elem string
			resf3elem = *resf3iter
			resf3 = append(resf3, &resf3elem)
		}
		res.SetSecurityGroupIds(resf3)
	}
	if r.ko.Spec.ResourcesVPCConfig.SubnetIDs != nil {
		resf4 := []*string{}
		for _, resf4iter := range r.ko.Spec.ResourcesVPCConfig.SubnetIDs {
			var resf4elem string
			resf4elem = *resf4iter
			resf4 = append(resf4, &resf4elem)
		}
		res.SetSubnetIds(resf4)
	}

	return res
}
