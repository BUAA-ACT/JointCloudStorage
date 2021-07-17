package dao

import "testing"

const (
	Version = "v0.2"

	CollectionCloud  = "Cloud"
	CollectionUser   = "User"
	CollectionFile   = "File"
	MigrationAdvice2 = "MigrationAdvice"
)

func TestDao_InsertMigrationAdvice(t *testing.T) {
	// Init DAO instance
	var err error
	db := GetDatabaseInstance()
	m := MigrationAdvice{
		UserId: "zhangjh",
		StoragePlanOld: StoragePlan{
			N:           2,
			K:           1,
			StorageMode: "Replica",
			Clouds: []Cloud{
				{
					CloudID:      "aliyun-qingdao",
					Endpoint:     "oss-cn-qingdao.aliyuncs.com",
					AccessKey:    "",
					SecretKey:    "",
					StoragePrice: 0.12,
					TrafficPrice: 0.25,
					Availability: 0.9995,
					Status:       "UP",
					Location:     "120.382109,36.075311",
					Address:      "localhost:8282",
					CloudName:    "阿里云-青岛",
					ProviderName: "aliyun",
				},
				{
					CloudID:      "aliyun-hohhot",
					Endpoint:     "oss-cn-huhehaote.aliyuncs.com",
					AccessKey:    "",
					SecretKey:    "",
					StoragePrice: 0.1,
					TrafficPrice: 0.5,
					Availability: 0.9999,
					Status:       "UP",
					Location:     "111.744578,40.873867",
					Address:      "localhost:8082",
					CloudName:    "阿里云-呼和浩特",
					ProviderName: "aliyun",
				},
			},
		},
		StoragePlanNew: StoragePlan{
			N:           2,
			K:           1,
			StorageMode: "Replica",
			Clouds: []Cloud{
				{
					CloudID:      "aliyun-hangzhou",
					Endpoint:     "oss-cn-hangzhou.aliyuncs.com",
					AccessKey:    "",
					SecretKey:    "",
					StoragePrice: 0.12,
					TrafficPrice: 0.25,
					Availability: 0.9995,
					Status:       "UP",
					Location:     "120.188938,30.301958",
					Address:      "localhost:8182",
					CloudName:    "阿里云-杭州",
					ProviderName: "aliyun",
				},
				{
					CloudID:      "aliyun-hohhot",
					Endpoint:     "oss-cn-huhehaote.aliyuncs.com",
					AccessKey:    "",
					SecretKey:    "",
					StoragePrice: 0.1,
					TrafficPrice: 0.5,
					Availability: 0.9999,
					Status:       "UP",
					Location:     "111.744578,40.873867",
					Address:      "localhost:8082",
					CloudName:    "阿里云-呼和浩特",
					ProviderName: "aliyun",
				},
			},
		},
		CloudsOld: []Cloud{
			{
				CloudID:      "aliyun-qingdao",
				Endpoint:     "oss-cn-qingdao.aliyuncs.com",
				AccessKey:    "",
				SecretKey:    "",
				StoragePrice: 0.12,
				TrafficPrice: 0.25,
				Availability: 0.9995,
				Status:       "UP",
				Location:     "120.382109,36.075311",
				Address:      "localhost:8282",
				CloudName:    "阿里云-青岛",
				ProviderName: "aliyun",
			},
		},
		CloudsNew: []Cloud{
			{
				CloudID:      "aliyun-hangzhou",
				Endpoint:     "oss-cn-hangzhou.aliyuncs.com",
				AccessKey:    "",
				SecretKey:    "",
				StoragePrice: 0.12,
				TrafficPrice: 0.25,
				Availability: 0.9995,
				Status:       "UP",
				Location:     "120.188938,30.301958",
				Address:      "localhost:8182",
				CloudName:    "阿里云-杭州",
				ProviderName: "aliyun",
			},
		},
		Cost: 0,
	}
	err = db.InsertMigrationAdvice(m)
	t.Log(err.Error())

}
