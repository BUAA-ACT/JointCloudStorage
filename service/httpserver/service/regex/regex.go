package regex

import (
	"cloud-storage-httpserver/args"
	"cloud-storage-httpserver/model"
	"cloud-storage-httpserver/service/tools"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
)

func emailRegex(email interface{}) (string, bool) {
	if reflect.TypeOf(email).Kind() != reflect.String {
		return "", false
	}
	return email.(string), true
}

func verifyCodeRegex(code interface{}) (string, bool) {
	if reflect.TypeOf(code).Kind() != reflect.String {
		return "", false
	}
	return code.(string), true
}

func passwordRegex(code interface{}) (string, bool) {
	if reflect.TypeOf(code).Kind() != reflect.String {
		return "", false
	}
	return code.(string), true
}

func vendorRegex(vendor interface{}) (uint64, bool) {
	if reflect.TypeOf(vendor).Kind() != reflect.String {
		return 0, false
	}
	realVendor, err := strconv.ParseUint(vendor.(string), 10, 64)
	return realVendor, !tools.PrintError(err)
}

func storagePriceRegex(storagePrice interface{}) (float64, bool) {
	if reflect.TypeOf(storagePrice).Kind() != reflect.String {
		return 0, false
	}
	realStoragePrice, err := strconv.ParseFloat(storagePrice.(string), 64)
	return realStoragePrice, !tools.PrintError(err)
}

func trafficPriceRegex(trafficPrice interface{}) (float64, bool) {
	if reflect.TypeOf(trafficPrice).Kind() != reflect.String {
		return 0, false
	}
	realTrafficPrice, err := strconv.ParseFloat(trafficPrice.(string), 64)
	return realTrafficPrice, !tools.PrintError(err)
}

func availabilityRegex(availability interface{}) (float64, bool) {
	if reflect.TypeOf(availability).Kind() != reflect.String {
		return 0, false
	}
	realAvailability, err := strconv.ParseFloat(availability.(string), 64)
	return realAvailability, !tools.PrintError(err)
}

func latencyRegex(latency interface{}) (*map[string]uint64, bool) {
	if reflect.TypeOf(latency).Kind() != reflect.Map {
		return nil, false
	}
	var success = true
	userLatency := map[string]uint64{}
	for field, value := range latency.(map[string]interface{}) {
		if reflect.TypeOf(value).Kind() != reflect.Map {
			formatValue := fmt.Sprintf("%v", value)
			uintValue, err := strconv.ParseUint(formatValue, 10, 64)
			success = success && !tools.PrintError(err)
			userLatency[field] = uintValue
		}
	}
	return &userLatency, success
}

func storagePlanRegex(plan interface{}) (*model.StoragePlan, bool) {
	if reflect.TypeOf(plan).Kind() != reflect.Map {
		return nil, false
	}
	var storagePlan model.StoragePlan
	err := mapstructure.Decode(plan, &storagePlan)
	return &storagePlan, !tools.PrintError(err)
}

func cloudRegex(cloud interface{}) (*model.Cloud, bool) {
	if reflect.TypeOf(cloud).Kind() != reflect.Map {
		return nil, false
	}
	var realCloud model.Cloud
	err := mapstructure.Decode(cloud, &realCloud)
	return &realCloud, !tools.PrintError(err)
}

func statusRegex(status interface{}) (bool, bool) {
	if reflect.TypeOf(status).Kind() != reflect.String {
		return false, false
	}
	var realStatus bool
	realStatus, err := strconv.ParseBool(status.(string))
	return realStatus, !tools.PrintError(err)
}

func voteResultRegex(voteResult interface{}) (bool, bool) {
	if reflect.TypeOf(voteResult).Kind() != reflect.String {
		return false, false
	}
	var realVoteResult bool
	realVoteResult, err := strconv.ParseBool(voteResult.(string))
	return realVoteResult, !tools.PrintError(err)
}

func CheckRegex(value interface{}, field string) (interface{}, bool) {
	var formatValue interface{}
	if reflect.TypeOf(value).Kind() != reflect.Map {
		formatValue = fmt.Sprintf("%v", value)
	} else {
		formatValue = value
	}
	switch field {
	case args.FieldWordEmail, args.FieldWordNewEmail:
		return emailRegex(formatValue)
	case args.FieldWordVerifyCode:
		return verifyCodeRegex(formatValue)
	case args.FieldWordPassword, args.FieldWordNewPassword:
		return passwordRegex(formatValue)
	case args.FieldWordVendor:
		return vendorRegex(formatValue)
	case args.FieldWordStoragePrice:
		return storagePriceRegex(formatValue)
	case args.FieldWordTrafficPrice:
		return trafficPriceRegex(formatValue)
	case args.FieldWordAvailability:
		return availabilityRegex(formatValue)
	case args.FieldWordStoragePlan:
		return storagePlanRegex(formatValue)
	case args.FieldWordLatency:
		return latencyRegex(formatValue)
	case args.FieldWordStatus:
		return statusRegex(formatValue)
	case args.FieldWordVoteResult:
		return voteResultRegex(formatValue)
	case args.FieldWordCloud:
		return cloudRegex(formatValue)
	default:
		return formatValue, true
	}
}
