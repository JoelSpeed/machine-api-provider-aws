package fake

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/openshift/machine-api-provider-aws/pkg/actuators/machine"
	"github.com/openshift/machine-api-provider-aws/pkg/client"
	"k8s.io/client-go/kubernetes"
)

type awsClient struct {
}

func (c *awsClient) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return &ec2.DescribeImagesOutput{
		Images: []*ec2.Image{
			{
				ImageId: aws.String("ami-a9acbbd6"),
			},
		},
	}, nil
}

func (c *awsClient) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return machine.StubDescribeVPCs()
}

func (c *awsClient) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{
				SubnetId: aws.String("subnet-28fddb3c45cae61b5"),
			},
		},
	}, nil
}

func (c *awsClient) DescribeAvailabilityZones(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return &ec2.DescribeAvailabilityZonesOutput{}, nil
}

func (c *awsClient) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{
				GroupId: aws.String("sg-05acc3c38a35ce63b"),
			},
		},
	}, nil
}

func (c *awsClient) DescribePlacementGroups(*ec2.DescribePlacementGroupsInput) (*ec2.DescribePlacementGroupsOutput, error) {
	return &ec2.DescribePlacementGroupsOutput{}, nil
}

func (c *awsClient) DescribeDHCPOptions(input *ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
	return machine.StubDescribeDHCPOptions()
}

func (c *awsClient) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return &ec2.Reservation{
		Instances: []*ec2.Instance{
			{
				ImageId:    aws.String("ami-a9acbbd6"),
				InstanceId: aws.String("i-02fcb933c5da7085c"),
				State: &ec2.InstanceState{
					Code: aws.Int64(16),
				},
			},
		},
	}, nil
}

func (c *awsClient) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						ImageId:    aws.String("ami-a9acbbd6"),
						InstanceId: aws.String("i-02fcb933c5da7085c"),
						State: &ec2.InstanceState{
							Name: aws.String(ec2.InstanceStateNameRunning),
							Code: aws.Int64(16),
						},
						LaunchTime: aws.Time(time.Now()),
					},
				},
			},
		},
	}, nil
}

func (c *awsClient) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	// Feel free to extend the returned values
	return &ec2.TerminateInstancesOutput{}, nil
}

func (c *awsClient) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	// Feel free to extend the returned values
	return &ec2.DescribeVolumesOutput{}, nil
}

func (c *awsClient) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return &ec2.CreateTagsOutput{}, nil
}

func (c *awsClient) CreatePlacementGroup(input *ec2.CreatePlacementGroupInput) (*ec2.CreatePlacementGroupOutput, error) {
	return &ec2.CreatePlacementGroupOutput{}, nil
}

func (c *awsClient) DeletePlacementGroup(input *ec2.DeletePlacementGroupInput) (*ec2.DeletePlacementGroupOutput, error) {
	return &ec2.DeletePlacementGroupOutput{}, nil
}

func (c *awsClient) RegisterInstancesWithLoadBalancer(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
	// Feel free to extend the returned values
	return &elb.RegisterInstancesWithLoadBalancerOutput{}, nil
}

func (c *awsClient) ELBv2DescribeLoadBalancers(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	// Feel free to extend the returned values
	return &elbv2.DescribeLoadBalancersOutput{}, nil
}

func (c *awsClient) ELBv2DescribeTargetGroups(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	// Feel free to extend the returned values
	return &elbv2.DescribeTargetGroupsOutput{}, nil
}

func (c *awsClient) ELBv2DescribeTargetHealth(*elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	return &elbv2.DescribeTargetHealthOutput{}, nil
}

func (c *awsClient) ELBv2RegisterTargets(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
	// Feel free to extend the returned values
	return &elbv2.RegisterTargetsOutput{}, nil
}

func (c *awsClient) ELBv2DeregisterTargets(*elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error) {
	// Feel free to extend the returned values
	return &elbv2.DeregisterTargetsOutput{}, nil
}

// NewClient creates our client wrapper object for the actual AWS clients we use.
// For authentication the underlying clients will use either the cluster AWS credentials
// secret if defined (i.e. in the root cluster),
// otherwise the IAM profile of the master where the actuator will run. (target clusters)
func NewClient(kubeClient kubernetes.Interface, secretName, namespace, region string) (client.Client, error) {
	return &awsClient{}, nil
}
