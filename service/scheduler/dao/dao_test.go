package dao

import (
	"shaoliyin.me/jcspan/entity"
	"testing"
)

const (
	Version = "v0.2"

	MongoURI         = "mongodb://localhost:27017"
	DaoTestEnv       = "dao_test"
	CollectionCloud  = "Cloud"
	CollectionUser   = "User"
	CollectionFile   = "File"
	MigrationAdvice2 = "CollectionMigrationAdvice"
)

func TestDao_InsertMigrationAdvice(t *testing.T) {
	// Init DAO instance
	var err error

	databases := map[string]map[string]*CollectionConfig{
		DaoTestEnv: {
			MigrationAdvice2: nil,
			CollectionFile:   nil,
			CollectionUser:   nil,
			CollectionCloud:  nil,
		},
	}
	err = NewDao(MongoURI, databases)
	if err != nil {
		t.Fatalf("failed when create dao")
	}
	adviceCol := databases[DaoTestEnv][MigrationAdvice2].CollectionHandler
	m := entity.MigrationAdvice{
		UserId: "zhangjh",
		StoragePlanOld: entity.StoragePlan{
			N:           2,
			K:           1,
			StorageMode: "Replica",
			Clouds: []entity.Cloud{
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
		StoragePlanNew: entity.StoragePlan{
			N:           2,
			K:           1,
			StorageMode: "Replica",
			Clouds: []entity.Cloud{
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
		CloudsOld: []entity.Cloud{
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
		CloudsNew: []entity.Cloud{
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

	err = InsertMigrationAdvice(adviceCol, m)
	if err != nil {
		t.Log(err.Error())
	}

}
