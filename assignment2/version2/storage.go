package main

import (
	"github.com/globalsign/mgo"
)

type mongoDBStorage struct {
	session *mgo.Session
}

type hub struct {
	tld   string `bson:"tld"`
	link  string `bson:"link"`
	count int    `bson:"count"`
}

type sitePage struct {
	tld            string `bson:"tld"`
	link           string `bson:"link"`
	fetchTime      string `bson:"fetchTime"`
	title          string `bson:"title"`
	htmlPageSource string `bson:"htmlPageSource"`
	mainText       string `bson:"mainText"`
	mainTextSHA1   string `bson:"mainTextSHA1"`
}

func (storage *mongoDBStorage) getHubData(tld string) {

}

func (storage *mongoDBStorage) getFetchedData(tld string) {

}

func (storage *mongoDBStorage) restore(tld string) {
	// get Hub
	// get fetched

	// hub, fetched - hub -> return map
}

func (storage *mongoDBStorage) insert(collectionName string, data interface{}) {
	session := storage.session.Copy()
	defer session.Close()

	collection := session.DB(conf.MongoDB.Database).C(collectionName)

	err := collection.Insert(data)
	if err != nil {
		panic(err)
	}
}

func (storage *mongoDBStorage) deinit() {
	storage.session.Close()
}

func (storage *mongoDBStorage) ensureIndex(collectionName string, key ...string) {
	session := storage.session.Copy()
	defer session.Close()

	collection := session.DB(conf.MongoDB.Database).C(collectionName)

	err := collection.EnsureIndexKey(key...)
	if err != nil {
		panic(err)
	}
}

func (storage *mongoDBStorage) init() {
	session, err := mgo.Dial(conf.MongoDB.URL)
	storage.session = session
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
}

// func demo() {
// // http: //www.itfanr.cc/2017/06/28/golang-connect-to-mongodb/
// 	// 创建连接
// 	session, err := mgo.Dial(conf.MongoDB.Url)
// 	storage.session = session
// 	if err != nil {
// 		panic(err)
// 	}
// 	session.SetMode(mgo.Monotonic, true)

// 	storage.database = session.DB(conf.MongoDB.Database)
// 	c := storage.database.C("test")

// 	c := session.DB(conf.MongoDB.Database).C("test")

// 	err = c.EnsureIndexKey("Name")
// 	if err != nil {
// 		panic(err)
// 	}

// 	storage.database = session.DB(conf.MongoDB.Database) //数据库名称
// 	c := storage.database.C("test")                      // 集合名称

// 	// c := session.DB(conf.MongoDB.Database).C("test")

// 	err = c.EnsureIndexKey("Name")
// 	if err != nil {
// 		panic(err)
// 	}
// 	//插入数据
// 	err = c.Insert(&Person{"Tommy", "123456"}, &Person{"Hanleilei", "98765"},
// 		&Person{"喜洋洋", "98765"}, &Person{"灰太狼", "46577"},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//查询并赋值 Find().One()
// 	result := Person{}
// 	err = c.Find(bson.M{"name": "Tommy"}).One(&result)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//输出
// 	fmt.Println("Phone ", result.Phone)
// 	//集合中元素数量 Count()
// 	countNum, err := c.Count()
// 	fmt.Println("obj numbers ", countNum)
// 	//查询多条数据 Find().Iter()
// 	var onep = Person{}
// 	iter := c.Find(nil).Iter()
// 	for iter.Next(&onep) {
// 		fmt.Println("姓名 ", onep.Name)
// 	}
// 	//查询多条数据 Find().All()
// 	var personAll []Person
// 	err = c.Find(nil).All(&personAll)
// 	for i := 0; i < len(personAll); i++ {
// 		fmt.Println("Person ", personAll[i].Name, personAll[i].Phone)
// 	}
// 	//更新数据 Update()
// 	abc := Person{}
// 	err = c.Find(bson.M{"name": "Tommy"}).One(&abc)
// 	fmt.Println("Tommy phone is ", abc.Phone)
// 	err = c.Update(bson.M{"name": "Tommy"}, bson.M{"$set": bson.M{"phone": "10086"}})
// 	err = c.Find(bson.M{"name": "Tommy"}).One(&abc)
// 	fmt.Println("Tommy phone is ", abc.Phone)
// 	//删除数据 Remove()
// 	fmt.Println(c.Count())
// 	err = c.Remove(bson.M{"phone": "46577"})
// 	fmt.Println(c.Count())
// 	// c.DropCollection()
// 	fmt.Println(c.Count())
// 	fmt.Println("end")
// }
