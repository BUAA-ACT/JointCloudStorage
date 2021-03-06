package dao

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	RoleHost  = "HOST"
	RoleGuest = "GUEST"
)

const (
	AdviceStatusPending = "PENDING"
)

// Dao encapsulates database operations.
type Dao struct {
	client          *mongo.Client
	database        string
	cloudCollection string
	fileCollection  string
	userCollection  string
	migrationAdvice string
	keyCollection   string
}

type Cloud struct {
	CloudID      string  `bson:"cloud_id"`
	Endpoint     string  `bson:"endpoint"`
	AccessKey    string  `bson:"access_key" `
	SecretKey    string  `bson:"secret_key" `
	StoragePrice float64 `bson:"storage_price"`
	TrafficPrice float64 `bson:"traffic_price"`
	Availability float64 `bson:"availability"`
	Status       string  `bson:"status"`
	Location     string  `bson:"location"`
	Address      string  `bson:"address"`
	CloudName    string  `bson:"cloud_name"`
	ProviderName string  `bson:"provider_name"`
	Bucket       string  `json:"Bucket" bson:"bucket"`
}

type User struct {
	UserId            string             `bson:"user_id"`
	Email             string             `bson:"email"`
	Password          string             `bson:"password"`
	Nickname          string             `bson:"nickname"`
	Role              string             `bson:"role"`
	Avatar            string             `bson:"avatar"`
	LastModified      time.Time          `bson:"last_modified"`
	Preference        Preference         `bson:"preference"`
	StoragePlan       StoragePlan        `bson:"storage_plan"`
	DataStats         DataStats          `bson:"data_stats"`
	AccessCredentials []AccessCredential `bson:"access_credentials"`
	Status            string             `bson:"status"`
}

type Preference struct {
	Vendor       int            `bson:"vendor"`
	StoragePrice float64        `bson:"storage_price"`
	TrafficPrice float64        `bson:"traffic_price"`
	Availability float64        `bson:"availability"`
	Latency      map[string]int `bson:"latency"`
}

type StoragePlan struct {
	N            int     `bson:"n"`
	K            int     `bson:"k"`
	StorageMode  string  `bson:"storage_mode"`
	Clouds       []Cloud `bson:"clouds"`
	StoragePrice float64 `bson:"storage_price"`
	TrafficPrice float64 `bson:"traffic_price"`
	Availability float64 `bson:"availability"`
}

type DataStats struct {
	Volume          int64            `bson:"volume"`
	UploadTraffic   map[string]int64 `bson:"upload_traffic"`
	DownloadTraffic map[string]int64 `bson:"download_traffic"`
}

type AccessCredential struct {
	CloudID  string `bson:"cloud_id"`
	UserID   string `bson:"user_id"`
	Password string `bson:"password"`
}

type File struct {
	FileID            string    `bson:"file_id"`
	FileName          string    `bson:"file_name"`
	Owner             string    `bson:"owner"`
	Size              int64     `bson:"size"`
	LastModified      time.Time `bson:"last_modified"`
	SyncStatus        string    `bson:"sync_status"`
	LastReconstructed time.Time `bson:"last_reconstructed"`
	ReconstructStatus string    `bson:"reconstruct_status"`
	DownloadUrl       string    `bson:"download_url"`
}

type MigrationAdvice struct {
	UserId         string      `bson:"user_id"`
	StoragePlanOld StoragePlan `bson:"storage_plan_old"`
	StoragePlanNew StoragePlan `bson:"storage_plan_new"`
	CloudsOld      []Cloud     `bson:"clouds_old"`
	CloudsNew      []Cloud     `bson:"clouds_new"`
	Cost           float64     `bson:"cost"`
	Status         string      `bson:"status"`
}

type AccessKey struct {
	UserID     string    `json:"UserID,omitempty" bson:"user_id"`
	AccessKey  string    `json:"AccessKey,omitempty" bson:"access_key"`
	SecretKey  string    `json:"SecretKey,omitempty" bson:"secret_key"`
	Comment    string    `json:"Comment,omitempty" bson:"comment"`
	CreateTime time.Time `json:"CreateTime,omitempty" bson:"create_time"`
	Available  bool      `json:"Available,omitempty" bson:"available"`
}

// NewDao constructs a data access object (Dao).
func NewDao(mongoURI, database, cloudCollection, userCollection, fileCollection, migrationAdvice, keyCollection string) (*Dao, error) {
	dao := &Dao{
		database:        database,
		cloudCollection: cloudCollection,
		userCollection:  userCollection,
		fileCollection:  fileCollection,
		migrationAdvice: migrationAdvice,
		keyCollection:   keyCollection,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	dao.client = client
	err = dao.ensureIndex("cloud_id", true, dao.cloudCollection)
	if err != nil {
		return nil, err
	}

	// TODO: ensure user and file

	return dao, nil
}

func (d *Dao) ensureIndex(index string, unique bool, collection string) error {
	col := d.client.Database(d.database).Collection(collection)
	idx := mongo.IndexModel{
		Keys: bson.M{
			index: 1,
		},
		Options: &options.IndexOptions{
			Unique: &unique,
		},
	}

	_, err := col.Indexes().CreateOne(context.TODO(), idx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCloud insert new cloud info to database.
func (d *Dao) UpdateCloud(cloud Cloud) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"cloud_id": cloud.CloudID,
		},
		bson.M{
			"$set": bson.M{
				"storage_price": cloud.StoragePrice,
				"traffic_price": cloud.TrafficPrice,
				"availability":  cloud.Availability,
				"status":        cloud.Status,
				"location":      cloud.Location,
				"access_key":    cloud.AccessKey,
				"secret_key":    cloud.SecretKey,
				"address":       cloud.Address,
				"bucket":        cloud.Bucket,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAllClouds return the info of given bucket.
func (d *Dao) GetAllClouds() ([]Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var clouds []Cloud
	cur, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Cloud
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		// ??????????????????
		//elem.AccessKey = ""
		//elem.SecretKey = ""
		clouds = append(clouds, elem)
	}

	return clouds, nil
}

func (d *Dao) GetOtherClouds(cid string) ([]Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var clouds []Cloud
	cur, err := col.Find(context.TODO(), bson.M{"cloud_id": bson.M{"$ne": cid}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Cloud
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		//// ??????????????????
		//elem.AccessKey = ""
		//elem.SecretKey = ""
		clouds = append(clouds, elem)
	}

	return clouds, nil
}

func (d *Dao) GetCloud(cid string) (Cloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var cloud Cloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": cid}).Decode(&cloud)

	//// ??????????????????
	//cloud.AccessKey = ""
	//cloud.SecretKey = ""
	return cloud, err
}

func (d *Dao) GetCloudNum() (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	num, err := col.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return 0, err
	} else {
		return int(num), nil
	}
}
func (d *Dao) GetFile(fid string) (File, error) {
	col := d.client.Database(d.database).Collection(d.fileCollection)

	var file File
	err := col.FindOne(context.TODO(), bson.M{"file_id": fid}).Decode(&file)
	return file, err
}

func (d *Dao) GetUser(uid string) (User, error) {
	col := d.client.Database(d.database).Collection(d.userCollection)

	var user User
	err := col.FindOne(context.TODO(), bson.M{"user_id": uid}).Decode(&user)
	return user, err
}

func (d *Dao) GetAllUser() ([]User, error) {
	col := d.client.Database(d.database).Collection(d.userCollection)

	var users []User
	cur, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		users = append(users, elem)
	}

	return users, nil
}

func (d *Dao) InsertMigrationAdvice(adv MigrationAdvice) error {
	col := d.client.Database(d.database).Collection(d.migrationAdvice)
	_, err := col.ReplaceOne(
		context.TODO(),
		bson.M{
			"user_id": adv.UserId,
		},
		adv,
		&options.ReplaceOptions{Upsert: aws.Bool(true)},
	)
	if err != nil {
		return err
	}

	return nil
}

func bool2pointer(b bool) *bool {
	return &b
}

func (d *Dao) InsertCloud(cloud Cloud) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.InsertOne(
		context.TODO(),
		cloud,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) InsertUser(user User) error {
	col := d.client.Database(d.database).Collection(d.userCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"user_id": user.UserId,
		},
		bson.M{
			"$set": user,
		},
		&options.UpdateOptions{
			Upsert: bool2pointer(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) InsertFiles(files []File) error {
	fs := make([]interface{}, len(files))
	for i := range files {
		fs[i] = files[i]
	}

	col := d.client.Database(d.database).Collection(d.fileCollection)
	for _, file := range files {
		_, err := col.UpdateOne(
			context.TODO(),
			bson.M{
				"file_id": file.FileID,
			},
			bson.D{
				{"$set", file},
			},
			&options.UpdateOptions{
				Upsert: bool2pointer(true),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Dao) DeleteFiles(files []File) error {
	var fs []string
	for _, v := range files {
		fs = append(fs, v.FileID)
	}

	col := d.client.Database(d.database).Collection(d.fileCollection)
	_, err := col.DeleteMany(
		context.TODO(),
		bson.M{
			"file_id": bson.M{"$in": fs},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) ChangeVolume(uid string, op string, files []File) error {
	var sum int64
	for _, v := range files {
		sum += v.Size
	}
	if op == "Delete" {
		sum = -sum
	}

	col := d.client.Database(d.database).Collection(d.userCollection)
	_, err := col.UpdateOne(
		context.TODO(),
		bson.M{
			"user_id": uid,
		},
		bson.M{
			"$inc": bson.M{"data_stats.volume": sum},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dao) DeleteUser(uid string) error {
	// ?????????????????????????????????
	col := d.client.Database(d.database).Collection(d.fileCollection)
	_, err := col.DeleteMany(
		context.TODO(),
		bson.M{
			"owner": uid,
		},
	)
	if err != nil {
		return err
	}

	// ????????????
	col = d.client.Database(d.database).Collection(d.userCollection)
	_, err = col.DeleteOne(
		context.TODO(),
		bson.M{
			"user_id": uid,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

/*******************************************************
* ??????accesskey???secretKey??????????????????
 ********************************************************/

//??????,????????????????????????
func (dao *Dao) KeyUpsert(ak AccessKey) error {
	col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	operation := bson.M{
		"$set": ak,
	}
	option := options.UpdateOptions{
		Upsert: bool2pointer(true),
	}
	_, err := col.UpdateOne(context.TODO(), filter, operation, &option)
	if err != nil {
		return err
	}

	return nil
}

//??????key
func (dao *Dao) DeleteKey(ak AccessKey) error {
	col := dao.client.Database(dao.database).Collection(dao.keyCollection)

	filter := bson.M{
		"access_key": ak.AccessKey,
	}

	_, err := col.DeleteOne(context.TODO(), filter)

	if err != nil {
		return err
	}
	return nil
}

/*
 * ????????????????????????????????????voteCloud
 */
type VoteCloud struct {
	Id      string `bson:"cloud_id" json:"id"`
	Cloud   Cloud  `bson:"cloud" json:"cloud"`
	VoteNum int    `bson:"vote_num" json:"vote_num"`
	Address string `bson:"address" json:"address"`
}

func (d *Dao) InsertVoteCloud(cloud VoteCloud) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	_, err := col.InsertOne(
		context.TODO(),
		cloud,
	)
	if err != nil {
		return err
	}
	return nil
}

//get the number of clouds whose id is cid
func (d *Dao) CloudsCount(cid string) (int64, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	count, err := col.CountDocuments(context.TODO(), bson.M{"cloud_id": cid})
	if err != nil {
		return count, err
	} else {
		return count, nil
	}
}

//delete the cloud
func (d *Dao) DeleteVoteCloud(id string) error {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	_, err := col.DeleteOne(context.TODO(), bson.M{"cloud_id": id})
	if err != nil {
		return err
	} else {
		return nil
	}
}

//add vote number
func (d *Dao) AddVoteNum(vote int, id string) (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	res, err := col.UpdateOne(
		context.TODO(),
		bson.M{"cloud_id": id},
		bson.M{
			"$inc": bson.M{"vote_num": vote},
		})
	if err != nil {
		return int(res.ModifiedCount), err
	} else {
		return int(res.ModifiedCount), nil
	}
}

//Get struct voteCloud by id
func (d *Dao) GetVoteCloud(id string) (VoteCloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

//Get all voteCloud in collection voteCloud
func (d *Dao) GetAllVoteCloud() ([]VoteCloud, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)

	var result []VoteCloud
	cur, err := col.Find(context.TODO(), bson.M{})
	defer cur.Close(context.TODO())
	if err != nil {
		return result, err
	}

	for cur.Next(context.TODO()) {
		var cloud VoteCloud
		if err = cur.Decode(&cloud); err != nil {
			return result, err
		}
		result = append(result, cloud)
	}
	return result, nil
}

//Get the vote number of the cloud with id
func (d *Dao) GetVoteNumber(id string) (int, error) {
	col := d.client.Database(d.database).Collection(d.cloudCollection)
	var result VoteCloud
	err := col.FindOne(context.TODO(), bson.M{"cloud_id": id}).Decode(&result)
	if err != nil {
		return -1, err
	} else {
		return result.VoteNum, nil
	}
}
