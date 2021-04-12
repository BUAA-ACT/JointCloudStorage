db.Cloud.insertMany(
    [
        {
            "cloud_id": "aliyun-hohhot",
            "storage_price": 0.12,
            "traffic_price": 0.5,
            "availability": 0.9999,
            "status": "UP",
            "endpoint": "oss-cn-huhehaote.aliyuncs.com",
            "access_key": "LTAI4G3PCfrg7aXQ6EvuDo25",
            "secret_key": "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0",
			"location": "111.744578,40.873867",
            "bucket": "jcspan-huhehaote",
            "address": "aliyun-hohhot.jointcloudstorage.cn"
        },
        {
            "cloud_id": "aliyun-hangzhou",
            "storage_price": 0.12,
            "traffic_price": 0.5,
            "availability": 0.99995,
            "status": "UP",
            "endpoint": "oss-cn-hangzhou.aliyuncs.com",
            "access_key": "LTAI4G3PCfrg7aXQ6EvuDo25",
            "secret_key": "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0",
            "location": "120.188938,30.301958",
            "bucket": "jcspan-hangzhou",
            "address": "aliyun-hangzhou.jointcloudstorage.cn",
        },
        {
            "cloud_id": "aliyun-qingdao",
            "storage_price": 0.12,
            "traffic_price": 0.5,
            "availability": 0.99995,
            "status": "UP",
            "endpoint": "oss-cn-qingdao.aliyuncs.com",
            "access_key": "LTAI4G3PCfrg7aXQ6EvuDo25",
            "secret_key": "5bmnIvUqvuuAG1j6QuWuhJ73MWAHE0",
            "location": "120.382109,36.075311",
            "bucket": "jcspan-qingdao",
            "address": "aliyun-qingdao.jointcloudstorage.cn",
        }
    ]
)