package main

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat/combin"
	"shaoliyin.me/jcspan/dao"
)

func storagePlan(param GetStoragePlanParam, clouds []dao.Cloud) GetStoragePlanData {
	// 初始化参数
	N := len(clouds)     // 可用云服务数量
	nMin := param.Vendor // 存储方案所包含云服务的数量下界
	//nMax := minInt(N, 5) // 存储方案所包含云服务的数量上界
	nMax := nMin

	sMin := 999.9 // 当前最小存储成本
	tMin := 999.9 // 当前最小流量成本

	var storageFirst dao.StoragePlan
	var trafficFirst dao.StoragePlan

	self, others := splitClouds(clouds)

	for n := nMin; n <= nMax; n++ {
		for _, cb := range combin.Combinations(N-1, n-1) {
			// 从其他云中选出n-1个云，再加上自己
			cls := []dao.Cloud{self}
			for _, i := range cb {
				cls = append(cls, others[i])
			}
			plan := dao.StoragePlan{
				N:           n,
				StorageMode: ECMode,
				Clouds:      cls,
			}

			// 遍历可能的k值
			for k := (n + 1) / 2; k <= n; k++ {
				if k == n {
					// Trick: 纠删码模式下k不能等于n，所以可以把n映射到1，即多副本模式
					plan.K = 1
				} else {
					plan.K = k
				}

				if plan.K == 1 {
					plan.StorageMode = ReplicaMode
				}

				// 计算并验证存储成本
				s := calStoragePrice(plan)
				if s > param.StoragePrice {
					continue
				}

				// 计算并验证流量成本
				t := calTrafficPrice(plan, true)
				if t > param.TrafficPrice {
					continue
				}

				// 计算并验证可用性
				a := calAvailability(plan)
				if a < param.Availability {
					continue
				}

				// TODO：计算并验证延迟要求

				plan.StoragePrice = s
				plan.TrafficPrice = t
				plan.Availability = a

				if s < sMin {
					sMin = s
					storageFirst = plan
				}

				if t < tMin {
					tMin = t
					trafficFirst = plan
				}
			}
		}
	}

	return GetStoragePlanData{
		StoragePriceFirst: storageFirst,
		TrafficPriceFirst: trafficFirst,
	}
}

func calStoragePrice(plan dao.StoragePlan) float64 {
	var price float64

	for _, c := range plan.Clouds {
		price += c.StoragePrice
	}
	price /= float64(plan.K)

	return price
}

func calTrafficPrice(plan dao.StoragePlan, resort bool) float64 {
	clouds := plan.Clouds
	if resort {
		sort.Slice(clouds, func(i, j int) bool {
			if clouds[i].TrafficPrice != clouds[j].TrafficPrice {
				return clouds[i].TrafficPrice < clouds[j].TrafficPrice
			} else if clouds[i].StoragePrice != clouds[j].StoragePrice {
				return clouds[i].StoragePrice < clouds[j].StoragePrice
			} else {
				return clouds[i].Availability > clouds[j].Availability
			}
		})
	}

	if plan.K == 1 {
		return clouds[0].TrafficPrice
	} else {
		var price float64
		for i := 0; i < plan.K; i++ {
			price += clouds[i].TrafficPrice
		}
		price /= float64(plan.K)
		return price
	}
}

func calAvailability(plan dao.StoragePlan) float64 {
	var res float64

	if plan.K == 1 {
		res = 1
		for _, c := range plan.Clouds {
			res *= 1 - c.Availability
		}
		return 1 - res
	} else {
		var avg float64
		for _, c := range plan.Clouds {
			avg += c.Availability
		}
		avg /= float64(len(plan.Clouds))

		for i := plan.K; i <= plan.N; i++ {
			// TODO: 准确值
			res += math.Pow(avg, float64(i)) * math.Pow((1-avg), float64(plan.N-i)) * float64(combin.Binomial(plan.N, i))
		}
		return res
	}
}

func downloadPlan(plan dao.StoragePlan, clouds []dao.Cloud) GetDownloadPlanData {
	resp := GetDownloadPlanData{
		StorageMode: plan.StorageMode,
	}

	cloudMap := make(map[string]dao.Cloud)
	for _, v := range clouds {
		cloudMap[v.CloudID] = v
	}

	for i, v := range plan.Clouds {
		if c := cloudMap[v.CloudID]; c.Status == "UP" {
			resp.Clouds = append(resp.Clouds, c)
			resp.Index = append(resp.Index, i)
		}
		if len(resp.Clouds) >= plan.K {
			break
		}
	}
	return resp
}

func reSchedule(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		<-t.C
		requestID := uuid.New().String()
		logInfo("Starting to reschedule storage plans", requestID)
		users, err := db.GetAllUser()
		if err != nil {
			logError(err, requestID, "GetAllUser failed")
			continue
		}

		clouds, err := db.GetAllClouds()
		if err != nil {
			logError(err, requestID, "GetAllClouds failed")
			continue
		}

		for _, u := range users {
			// 跳过新用户
			if u.Preference.Vendor == 0 || u.Role == dao.RoleGuest {
				continue
			}

			// 计算新方案
			plan := storagePlan(GetStoragePlanParam(u.Preference), clouds).StoragePriceFirst

			// 对比新旧方案
			if plan.K != u.StoragePlan.K || plan.N != u.StoragePlan.N {
				continue
			}

			reordered, deleted, added := transform(u.StoragePlan.Clouds, plan.Clouds)
			if len(deleted) == 0 || len(added) == 0 {
				continue
			}
			plan.Clouds = reordered
			adv := dao.MigrationAdvice{
				UserId:         u.UserId,
				StoragePlanOld: u.StoragePlan,
				StoragePlanNew: plan,
				CloudsOld:      deleted,
				CloudsNew:      added,
				Cost:           42,
			}

			// 写入数据库
			err = db.InsertMigrationAdvice(adv)
			if err != nil {
				logError(err, requestID, "InsertMigrationAdvice failed", adv)
				continue
			}
			logInfo("Created new migration advice", requestID, adv)
		}
		logInfo("Finish reschedule storage plans", requestID)
	}
}

func splitClouds(clouds []dao.Cloud) (dao.Cloud, []dao.Cloud) {
	var self dao.Cloud
	var others []dao.Cloud

	for _, v := range clouds {
		if v.CloudID == *flagCloudID {
			self = v
		} else {
			others = append(others, v)
		}
	}

	return self, others
}

func transform(old []dao.Cloud, new []dao.Cloud) (reordered, deleted, added []dao.Cloud) {
	index := make([]int, len(old))
	for i := range index {
		index[i] = -1
	}

	for i, v := range new {
		for j, w := range old {
			if v.CloudID == w.CloudID {
				index[i] = j
				break
			}
		}
	}

	reordered = make([]dao.Cloud, len(old))
	for i, v := range new {
		if index[i] == -1 {
			added = append(added, v)
		} else {
			reordered[index[i]] = v
		}
	}
	j := 0
	for i, v := range reordered {
		if v.CloudID == "" {
			deleted = append(deleted, old[i])
			reordered[i] = added[j]
			j++
		}
	}

	return
}
