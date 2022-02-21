package awsplacementgroup

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	awsclient "github.com/openshift/machine-api-provider-aws/pkg/client"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// upstreamMachineClusterIDLabel is the label that a machine must have to identify the cluster to which it belongs
const upstreamMachineClusterIDLabel = "sigs.k8s.io/cluster-api-cluster"
const requeueAfter = 10 * time.Second

//
const (
	// AWSPlacementGroupCreationSucceeded indicates placement group creation success.
	AWSPlacementGroupCreationSucceededConditionReason string = "AWSPlacementGroupCreationSucceeded"
	// AWSPlacementGroupCreationFailed indicates placement group creation failure.
	AWSPlacementGroupCreationFailedConditionReason string = "AWSPlacementGroupCreationFailed"
	// AWSPlacementGroupDeletionFailed indicates placement group creation failure.
	AWSPlacementGroupDeletionFailedConditionReason string = "AWSPlacementGroupDeletionFailed"
	// AWSPlacementGroupConfigurationMismatch indicates placement group configuration doesn't match the real-world resource configuration.
	AWSPlacementGroupConfigurationMismatchConditionReason string = "AWSPlacementGroupConfigurationMismatch"
	// AWSPlacementGroupConfigurationInSyncConditionReason indicates placement group configuration is in sync with configuration.
	AWSPlacementGroupConfigurationInSyncConditionReason string = "AWSPlacementGroupConfigurationInSync"
)

// Reconciler reconciles AWSPlacementGroup.
type Reconciler struct {
	Client              client.Client
	Log                 logr.Logger
	AWSClient           awsclient.Client
	ConfigManagedClient client.Client
	Infrastructure      *configv1.Infrastructure

	recorder record.EventRecorder
	scheme   *runtime.Scheme
}

// SetupWithManager creates a new controller for a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	_, err := ctrl.NewControllerManagedBy(mgr).
		For(&machinev1.AWSPlacementGroup{}).
		WithOptions(options).
		Build(r)

	if err != nil {
		return fmt.Errorf("failed setting up with a controller manager: %w", err)
	}

	r.recorder = mgr.GetEventRecorderFor("awsplacementgroup-controller")
	r.scheme = mgr.GetScheme()
	return nil
}

func (r *Reconciler) reconcile(ctx context.Context, logger logr.Logger, awsPlacementGroup *machinev1.AWSPlacementGroup) (ctrl.Result, error) {
	// check AWS for the configuration of the placement group and reflect this in the status of the object
	err := reflectObservedConfiguration(r.AWSClient, awsPlacementGroup)
	if err != nil {
		return ctrl.Result{}, err
	}

	// only proceed if the AWSPlacementGroup is Managed
	if awsPlacementGroup.Spec.ManagementSpec.ManagementState == machinev1.UnmanagedManagementState {
		// this AWSPlacementGroup is now unmanaged so clean up any machine finalizer if there is any
		// as we don't want to delete placement group resources on AWS, if we are in Unmanaged mode
		if controllerutil.ContainsFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer) {
			controllerutil.RemoveFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer)
			klog.Infof("%s: removing finalizer from AWSPlacementGroup", awsPlacementGroup.Name)
			return ctrl.Result{}, nil
		}

		klog.Infof("%s: ignoring Unmanaged AWSPlacementGroup", awsPlacementGroup.Name)
		// TODO: remove requeueafter?
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	if awsPlacementGroup.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer) {
			controllerutil.AddFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer)
			klog.Infof("%s: adding finalizer to AWSPlacementGroup", awsPlacementGroup.Name)
			return ctrl.Result{RequeueAfter: requeueAfter}, nil
		}
	}

	// if object DeletionTimestamp is not zero, it means it is being deleted
	// so clean up relevant resources
	if !awsPlacementGroup.DeletionTimestamp.IsZero() {
		// no-op if finalizer has been removed.
		if !controllerutil.ContainsFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer) {
			klog.Infof("%v: reconciling AWSPlacementGroup causes a no-op as there is no finalizer", awsPlacementGroup.Name)
			return ctrl.Result{RequeueAfter: requeueAfter}, nil
		}

		klog.Infof("%s: reconciling AWSPlacementGroup triggers delete", awsPlacementGroup.Name)

		if err = deletePlacementGroup(r.AWSClient, awsPlacementGroup, r.Infrastructure); err != nil {
			condition := metav1.Condition{
				Type:   "Deleting",
				Status: metav1.ConditionTrue,
				Reason: AWSPlacementGroupDeletionFailedConditionReason,
			}
			condition.Message = err.Error()
			awsPlacementGroup.Status.Conditions = setCondition(awsPlacementGroup.Status.Conditions, condition)
			return ctrl.Result{RequeueAfter: requeueAfter}, fmt.Errorf("%s: failed to delete placement group: %v", awsPlacementGroup.Name, err)
		}

		// remove finalizer on successful deletion
		controllerutil.RemoveFinalizer(awsPlacementGroup, machinev1beta1.MachineFinalizer)
		klog.Infof("%s: removing finalizer on successful placement group deletion", awsPlacementGroup.Name)

		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	// conditionally create the placement group
	if err := checkOrCreatePlacementGroup(r.AWSClient, awsPlacementGroup, r.Infrastructure); err != nil {
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// Reconcile implements controller runtime Reconciler interface.
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("name", req.Name, "namespace", req.Namespace)
	logger.V(3).Info("Reconciling")

	awsPlacementGroup := &machinev1.AWSPlacementGroup{}
	if err := r.Client.Get(ctx, req.NamespacedName, awsPlacementGroup); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// get the infrastructure object
	infra := &configv1.Infrastructure{}
	infraName := client.ObjectKey{Name: awsclient.GlobalInfrastuctureName}
	if err := r.Client.Get(ctx, infraName, infra); err != nil {
		return ctrl.Result{}, err
	}

	// check if the CredentialsSecret is defined
	credentialsSecretName := ""
	if awsPlacementGroup.Spec.CredentialsSecret != nil {
		credentialsSecretName = awsPlacementGroup.Spec.CredentialsSecret.Name
	}

	// get an awsClient
	awsClient, err := awsclient.NewValidatedClient(
		r.Client, credentialsSecretName, awsPlacementGroup.Namespace,
		infra.Status.PlatformStatus.AWS.Region, r.ConfigManagedClient)
	if err != nil {
		return ctrl.Result{}, err
	}

	r.Infrastructure = infra
	r.AWSClient = awsClient

	originalAWSPlacementGroupToPatch := client.MergeFrom(awsPlacementGroup.DeepCopy())

	result, err := r.reconcile(ctx, logger, awsPlacementGroup)
	if err != nil {
		logger.Error(err, "failed to reconcile AWSPlacementGroup")
		r.recorder.Eventf(awsPlacementGroup, corev1.EventTypeWarning, "ReconcileError", "%v", err)
		// we don't return here so we want to attempt to patch the AWSPlacementGroup regardless of an error.
	}

	originalAWSPlacementGroupStatus := awsPlacementGroup.Status.DeepCopy()

	if err := r.Client.Patch(ctx, awsPlacementGroup, originalAWSPlacementGroupToPatch); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch AWSPlacementGroup: %v", err)
	}

	awsPlacementGroup.Status = *originalAWSPlacementGroupStatus

	if err := r.Client.Status().Patch(ctx, awsPlacementGroup, originalAWSPlacementGroupToPatch); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch AWSPlacementGroup Status: %v", err)
	}

	return result, err
}

// mergeInfrastructureAndAWSPlacementGroupSpecTags merge list of tags from AWSPlacementGroup provider spec and Infrastructure object platform spec.
// Machine tags have precedence over Infrastructure
func mergeInfrastructureAndAWSPlacementGroupSpecTags(awsPlacementGroupSpecTags []machinev1beta1.TagSpecification, infra *configv1.Infrastructure) []machinev1beta1.TagSpecification {
	if infra == nil || infra.Status.PlatformStatus == nil || infra.Status.PlatformStatus.AWS == nil || infra.Status.PlatformStatus.AWS.ResourceTags == nil {
		return awsPlacementGroupSpecTags
	}

	mergedList := []machinev1beta1.TagSpecification{}
	mergedList = append(mergedList, awsPlacementGroupSpecTags...)

	for _, tag := range infra.Status.PlatformStatus.AWS.ResourceTags {
		mergedList = append(mergedList, machinev1beta1.TagSpecification{Name: tag.Key, Value: tag.Value})
	}

	return mergedList
}

// buildPlacementGroupTagList compile a list of ec2 tags from AWSPlacementGroup provider spec and infrastructure object platform spec
func buildPlacementGroupTagList(awsPlacementGroup string, awsPlacementGroupSpecTags []machinev1beta1.TagSpecification, infra *configv1.Infrastructure) []*ec2.Tag {
	rawTagList := []*ec2.Tag{}

	mergedTags := mergeInfrastructureAndAWSPlacementGroupSpecTags(awsPlacementGroupSpecTags, infra)

	clusterID := infra.Status.InfrastructureName

	for _, tag := range mergedTags {
		// AWS tags are case sensitive, so we don't need to worry about other casing of "Name"
		if !strings.HasPrefix(tag.Name, "kubernetes.io/cluster/") && tag.Name != "Name" {
			rawTagList = append(rawTagList, &ec2.Tag{Key: aws.String(tag.Name), Value: aws.String(tag.Value)})
		}
	}
	rawTagList = append(rawTagList, []*ec2.Tag{
		{Key: aws.String("kubernetes.io/cluster/" + clusterID), Value: aws.String("owned")},
		{Key: aws.String("Name"), Value: aws.String(awsPlacementGroup)},
	}...)

	return removeDuplicatedTags(rawTagList)
}

// Scan machine tags, and return a deduped tags list. The first found value gets precedence.
func removeDuplicatedTags(tags []*ec2.Tag) []*ec2.Tag {
	m := make(map[string]bool)
	result := []*ec2.Tag{}

	// look for duplicates
	for _, entry := range tags {
		if _, value := m[*entry.Key]; !value {
			m[*entry.Key] = true
			result = append(result, entry)
		}
	}
	return result
}

// isAWS4xxError will determine if the passed error is an AWS error with a 4xx status code.
func isAWS4xxError(err error) bool {
	if _, ok := err.(awserr.Error); ok {
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			if reqErr.StatusCode() >= 400 && reqErr.StatusCode() < 500 {
				return true
			}
		}
	}
	return false
}

// checkOrCreatePlacementGroup checks for the existence of a placement group on AWS and validates its config
// it proceeds to create one if such group doesn't exist.
func checkOrCreatePlacementGroup(client awsclient.Client, pg *machinev1.AWSPlacementGroup, infra *configv1.Infrastructure) error {
	placementGroups, err := client.DescribePlacementGroups(&ec2.DescribePlacementGroupsInput{
		GroupNames: []*string{aws.String(pg.Name)},
	})
	if err != nil && !isAWS4xxError(err) {
		// Ignore a 400 error as AWS will report an unknown placement group as a 400.
		return fmt.Errorf("could not describe placement groups: %v", err)
	}

	// more than one placement group matching
	if len(placementGroups.PlacementGroups) > 1 {
		return fmt.Errorf("%s: failed to check placement group: expected 1 placement group, got %d",
			pg.Name, len(placementGroups.PlacementGroups))
	}

	// placement group already exists
	if len(placementGroups.PlacementGroups) == 1 {
		if err := validatePlacementGroupConfig(pg, placementGroups.PlacementGroups[0]); err != nil {
			return err
		}
		setObservedConfiguration(pg, placementGroups.PlacementGroups[0])
		return nil
	}

	// we build a tag list for the placement group by inheriting user defined tags from infra
	tagList := buildPlacementGroupTagList(pg.Name, []machinev1beta1.TagSpecification{}, infra)

	// no placement group by that name existed, so we create one
	createPlacementGroupInput := &ec2.CreatePlacementGroupInput{
		GroupName: aws.String(pg.Name),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String(ec2.ResourceTypePlacementGroup),
				Tags:         tagList,
			},
		},
	}
	switch pg.Spec.ManagementSpec.Managed.GroupType {
	case machinev1.AWSSpreadPlacementGroupType:
		createPlacementGroupInput.SetStrategy(ec2.PlacementStrategySpread)
	case machinev1.AWSClusterPlacementGroupType:
		createPlacementGroupInput.SetStrategy(ec2.PlacementStrategyCluster)
	case machinev1.AWSPartitionPlacementGroupType:
		createPlacementGroupInput.SetStrategy(ec2.PlacementStrategyPartition)
		if pg.Spec.ManagementSpec.Managed.Partition != nil && pg.Spec.ManagementSpec.Managed.Partition.Count != 0 {
			createPlacementGroupInput.SetPartitionCount(int64(pg.Spec.ManagementSpec.Managed.Partition.Count))
		}
	default:
		return fmt.Errorf("unknown placement strategy %q: valid values are %s, %s, %s",
			pg.Spec.ManagementSpec.Managed.GroupType,
			machinev1.AWSSpreadPlacementGroupType,
			machinev1.AWSClusterPlacementGroupType,
			machinev1.AWSPartitionPlacementGroupType)
	}

	out, err := client.CreatePlacementGroup(createPlacementGroupInput)
	if err != nil {
		errStr := fmt.Sprintf("unable to create placement group: %v", err)
		// TODO(question): do we need to set conditions with these kinds of errors?
		//condition := metav1.Condition{
		//	Type:    "Ready",
		//	Status:  metav1.ConditionFalse,
		//	Reason:  AWSPlacementGroupCreationFailedConditionReason,
		//	Message: errStr,
		//}
		//pg.Status.Conditions = setCondition(pg.Status.Conditions, condition)
		return fmt.Errorf(errStr)
	}

	// set successful condition for placement group creation
	condition := metav1.Condition{Type: "Ready",
		Status: metav1.ConditionTrue,
		Reason: AWSPlacementGroupCreationSucceededConditionReason,
	}
	pg.Status.Conditions = setCondition(pg.Status.Conditions, condition)

	klog.Infof("%s: successfully created placement group %s", pg.Name, *out.PlacementGroup.GroupId)
	return nil
}

// validatePlacementGroupConfig validates that the configuration of the existing placement group
// matches the configuration of the AWSPlacementGroup spec
func validatePlacementGroupConfig(pg *machinev1.AWSPlacementGroup, placementGroup *ec2.PlacementGroup) error {
	if placementGroup == nil {
		return fmt.Errorf("found nil placement group")
	}
	if aws.StringValue(placementGroup.GroupName) != pg.Name {
		return fmt.Errorf("placement group name mismatch: wanted: %q, got: %q", pg.Name, aws.StringValue(placementGroup.GroupName))
	}
	var expectedPlacementGroupType string
	switch pg.Spec.ManagementSpec.Managed.GroupType {
	case machinev1.AWSSpreadPlacementGroupType:
		expectedPlacementGroupType = ec2.PlacementStrategySpread
	case machinev1.AWSClusterPlacementGroupType:
		expectedPlacementGroupType = ec2.PlacementStrategyCluster
	case machinev1.AWSPartitionPlacementGroupType:
		expectedPlacementGroupType = ec2.PlacementStrategyPartition
	default:
		return fmt.Errorf("unknown placement strategy %q: valid values are %s, %s and %s", pg.Spec.ManagementSpec.Managed.GroupType, machinev1.AWSSpreadPlacementGroupType, machinev1.AWSClusterPlacementGroupType, machinev1.AWSPartitionPlacementGroupType)
	}

	if aws.StringValue(placementGroup.Strategy) != expectedPlacementGroupType {
		// TODO: use invalidconfig
		//return mapierrors.InvalidMachineConfiguration("mismatch between configured placement group type and existing placement group type: wanted: %q, got: %q", expectedPlacementGroupType, aws.StringValue(placementGroup.Strategy))
		return fmt.Errorf("mismatch between configured placement group type and existing placement group type: wanted: %q, got: %q", expectedPlacementGroupType, aws.StringValue(placementGroup.Strategy))
	}

	if pg.Spec.ManagementSpec.Managed.GroupType == machinev1.AWSPartitionPlacementGroupType && pg.Spec.ManagementSpec.Managed.Partition != nil {
		if pg.Spec.ManagementSpec.Managed.Partition.Count != 0 && int64(pg.Spec.ManagementSpec.Managed.Partition.Count) != aws.Int64Value(placementGroup.PartitionCount) {
			// TODO: use invalidconfig
			//return mapierrors.InvalidMachineConfiguration("mismatch between configured placement group partition count and existing placement group partition count: wanted: %d, got: %d", placement.Spec.ManagementSpec.Managed.Partition.Count, aws.Int64Value(placementGroup.PartitionCount))
			return fmt.Errorf("mismatch between configured placement group partition count and existing placement group partition count: wanted: %d, got: %d", pg.Spec.ManagementSpec.Managed.Partition.Count, aws.Int64Value(placementGroup.PartitionCount))
		}
	}
	return nil
}

// deletePlacementGroup deletes the placement group for the machine
func deletePlacementGroup(client awsclient.Client, pg *machinev1.AWSPlacementGroup, infra *configv1.Infrastructure) error {

	// Check that the placement group exists
	placementGroups, err := client.DescribePlacementGroups(&ec2.DescribePlacementGroupsInput{
		GroupNames: []*string{aws.String(pg.Name)},
	})
	if err != nil && !isAWS4xxError(err) {
		// Ignore a 400 error as AWS will report an unknown placement group as a 400.
		return fmt.Errorf("could not describe placement groups: %v", err)
	}

	if len(placementGroups.PlacementGroups) > 1 {
		return fmt.Errorf("expected 1 placement group for name %q, got %d", pg.Name, len(placementGroups.PlacementGroups))
	}

	if len(placementGroups.PlacementGroups) == 0 {
		// This is the normal path, the named placement group doesn't exist
		return nil
	}

	placementGroup := placementGroups.PlacementGroups[0]

	clusterID := infra.Status.InfrastructureName

	// Check that the placement group has a cluster tag
	found := false
	for _, tag := range placementGroup.Tags {
		if aws.StringValue(tag.Key) == "kubernetes.io/cluster/"+clusterID && aws.StringValue(tag.Value) == "owned" {
			found = true
			break
		}
	}
	if !found {
		// TODO: should we return an error here?
		return fmt.Errorf("%s: AWSPlacementGroup wasn't created by machine-api and cannot be removed", pg.Name)
	}

	// check that the placement group contains no instances
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("placement-group-name"), Values: []*string{aws.String(pg.Name)}},
			{Name: aws.String("placement-group-name"), Values: []*string{aws.String(pg.Name)}},
		},
	})
	if err != nil {
		return fmt.Errorf("could not get the number of instances in the placement group %s: %v", pg.Name, err)
	}

	var instanceCount int
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			// we ignore Terminated instances
			// these should not count towards instances actively in a placement group
			if aws.StringValue(instance.State.Name) != "terminated" {
				instanceCount++
			}
		}
	}
	if instanceCount > 0 {
		return fmt.Errorf("%s: AWSPlacementGroup still constains %d instances and cannot be removed", pg.Name, instanceCount)
	}

	// only one placement group with the given name exists and it is empty, so we remove it
	deletePlacementGroupInput := &ec2.DeletePlacementGroupInput{GroupName: aws.String(pg.Name)}
	if _, err := client.DeletePlacementGroup(deletePlacementGroupInput); err != nil {
		return fmt.Errorf("unable to delete placement group: %v", err)
	}

	klog.Infof("%s: successfully deleted placement group %s", pg.Name, *placementGroup.GroupId)
	return nil
}

// reflectObservedConfiguration checks for the existence of a placement group on AWS and if that's the case
// it syncs its config with the ObservedConfiguration in the Status of the object
func reflectObservedConfiguration(client awsclient.Client, pg *machinev1.AWSPlacementGroup) error {

	placementGroups, err := client.DescribePlacementGroups(&ec2.DescribePlacementGroupsInput{
		GroupNames: []*string{aws.String(pg.Name)},
	})
	if err != nil && !isAWS4xxError(err) {
		// Ignore a 400 error as AWS will report an unknown placement group as a 400.
		return fmt.Errorf("could not describe placement groups: %v", err)
	}

	if len(placementGroups.PlacementGroups) < 1 {
		klog.Infof("%s: no matching placement group", pg.Name)
		return nil
	}

	if len(placementGroups.PlacementGroups) > 1 {
		return fmt.Errorf("expected 1 placement group for name %q, got %d", pg.Name, len(placementGroups.PlacementGroups))
	}

	setObservedConfiguration(pg, placementGroups.PlacementGroups[0])

	return nil
}

// setObservedConfiguration sets the configuration observed from the AWS placement group to
// the ObservedConfiguration filed in the Status of the object
func setObservedConfiguration(pg *machinev1.AWSPlacementGroup, placementGroup *ec2.PlacementGroup) {
	pg.Status.ObservedConfiguration.GroupType = machinev1.AWSPlacementGroupType(strings.Title(aws.StringValue(placementGroup.Strategy)))
	pg.Status.ObservedConfiguration.Partition = &machinev1.AWSPartitionPlacement{Count: int32(aws.Int64Value(placementGroup.PartitionCount))}

	condition := metav1.Condition{Type: "Ready"}

	var equal bool
	if pg.Spec.ManagementSpec.Managed != nil {
		equal = reflect.DeepEqual(pg.Status.ObservedConfiguration, *pg.Spec.ManagementSpec.Managed)
	}

	if equal {
		condition.Status = metav1.ConditionTrue
		condition.Reason = AWSPlacementGroupConfigurationInSyncConditionReason
	} else {
		condition.Status = metav1.ConditionFalse
		condition.Reason = AWSPlacementGroupConfigurationMismatchConditionReason
	}

	pg.Status.Conditions = setCondition(pg.Status.Conditions, condition)
}
