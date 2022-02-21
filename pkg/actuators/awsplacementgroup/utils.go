package awsplacementgroup

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// setCondition sets the condition on the existing Contions slice and returns the new slice.
// If Conditions does not already have a condition with the specified type,
// a condition will be added to the slice
// If the AWSPlacementGroup does already have a condition with the specified type,
// the condition will be updated if either of the following are true.
func setCondition(conditions []metav1.Condition, condition metav1.Condition) []metav1.Condition {
	now := metav1.Now()

	if existingCondition := findCondition(conditions, condition); existingCondition == nil {
		condition.LastTransitionTime = now
		conditions = append(conditions, condition)
	} else {
		updateExistingCondition(&condition, existingCondition)
	}

	return conditions
}

func findCondition(conditions []metav1.Condition, condition metav1.Condition) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == condition.Type {
			return &conditions[i]
		}
	}
	return nil
}

func updateExistingCondition(newCondition, existingCondition *metav1.Condition) {
	if !shouldUpdateCondition(newCondition, existingCondition) {
		return
	}

	if existingCondition.Status != newCondition.Status {
		existingCondition.LastTransitionTime = metav1.Now()
	}
	existingCondition.Status = newCondition.Status
	existingCondition.Reason = newCondition.Reason
	existingCondition.Message = newCondition.Message
}

func shouldUpdateCondition(newCondition, existingCondition *metav1.Condition) bool {
	return newCondition.Reason != existingCondition.Reason || newCondition.Message != existingCondition.Message
}
