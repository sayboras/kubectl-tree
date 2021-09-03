package main

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
)

type Type string // True False Unknown or ""
type Reason string

func extractStatus(obj unstructured.Unstructured) (Type, Reason) {
	jsonVal, _ := json.Marshal(obj.Object["status"])
	klog.V(6).Infof("status for object=%s/%s: %s", obj.GetKind(), obj.GetName(), string(jsonVal))
	statusF, ok := obj.Object["status"]
	if !ok {
		return "", ""
	}
	statusV, ok := statusF.(map[string]interface{})
	if !ok {
		return "", ""
	}
	conditionsF, ok := statusV["conditions"]
	if !ok {
		return "", ""
	}
	conditionsV, ok := conditionsF.([]interface{})
	if !ok {
		return "", ""
	}

	customType := false
	for _, cond := range conditionsV {
		condM, ok := cond.(map[string]interface{})
		if !ok {
			return "", ""
		}
		condType, ok := condM["type"].(string)
		if !ok {
			return "", ""
		}
		if condType == "Ready" {
			condStatus, _ := condM["status"].(string)
			if condStatus == "True" {
				condReason, _ := condM["reason"].(string)
				return Type(condType), Reason(condReason)
			}
		}
		if condType == "PermanentFailure" || condType == "Failure" || condType == "InterfaceChangeApplied" {
			customType = true
			condStatus, _ := condM["status"].(string)
			if condStatus == "True" {
				message, _ := condM["message"].(string)
				return Type(condType), Reason(message)
			}
		}
	}
	if customType {
		return Type("In Progress"), Reason("")
	}
	return "", ""
}
