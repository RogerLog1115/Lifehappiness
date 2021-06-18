package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type GetData struct {
	ID				bson.ObjectId	`json:"id"`
	LocationId		string	`bson:"locationId" json:"locationId"`
	LocationName	string	`bson:"locationName" json:"locationName"`
	Address			string	`bson:"address" json:"naddress"`
	Phone			string	`bson:"phone" json:"nphone"`
	Type			string	`bson:"type" json:"ntype"`
	PaymentNotes	string	`bson:"paymentNotes" json:"paymentNotes"`
	LimitNotes		string	`bson:"limitNotes" json:"limitNotes"`
}

type GetTaifuData struct {
	ID 		bson.ObjectId	`bson:"_id,omitempty" json:"id"`
	TFName	string			`bson:"公墓名稱" json:"tfName,公墓名稱"`
	TFType	string			`bson:"區別" json:"tfType,區別"`
	Address	string			`bson:"地址或鄰近地址" json:"address,地址或鄰近地址"`
	Contact	string			`bson:"聯絡人" json:"contact,聯絡人"`
	Phone	string			`bson:"電話" json:"fphone,電話"`
}

type FLegalData struct {
	ID 				bson.ObjectId	`bson:"_id" json:"id"`
	FacilityName	string			`bson:"facilityName" json:"facilityName"`
	FuPhone			string			`bson:"phone" json:"fuPhone"`
	FuUrl			string			`bson:"url"`
	FuEmail			string			`bson:"email" json:"fuEmail"`
	Address			string			`bson:"address" json:"address"`
}

type status struct {
	Status string `json:"status"`
}

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Fatal error config file:%v\n", err)
	}
}

func main() {
	r := gin.Default()
	r.Use(Cors())

	//r.LoadHTMLFiles("index.html")
	//db := DB()
	//data := &GetData{}
	r.LoadHTMLFiles("./asset/index.html")
	r.LoadHTMLGlob("./asset/pages/*")
	w := r.Group("/web")
	{
		w.StaticFS("/",http.Dir("./asset"))
	}
	a := r.Group("/api")
	{
		a.GET("/getnature",GetNatureData)
		a.GET("/getta",GetTaData)
		a.GET("/getlegal",GetFLData)
	}
	r.Run(":80")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func DB() *mgo.Database {
	logrus.Info("Server On")

	database := viper.GetString("DATABASE")
	session, err := mgo.Dial("database:27017")
	if err != nil {
		logrus.Error(err)
	}
	session.SetMode(mgo.Monotonic, true)
	if sessionErr := session.Ping();sessionErr != nil {
		logrus.Error(sessionErr)
	}

	session.SetPoolLimit(4096)
	db := session.DB(database)
	return db
}

func GetNatureData(g *gin.Context){
	s := []GetData{}
	d := DB()
	col := d.C("naturelist")
	err := col.Find(nil).Sort("locationId").All(&s)
	if err != nil{
		logrus.Println(err)
	}
	g.JSON(http.StatusOK,s)
}
func GetTaData(g *gin.Context){
	s := []GetTaifuData{}
	d := DB()
	col := d.C("tailegalfu")
	err := col.Find(nil).All(&s)
	if err != nil{
		logrus.Println(err)
	}
	g.JSON(http.StatusOK,s)
}
func GetFLData(g *gin.Context){
	s := []FLegalData{}
	d := DB()
	col := d.C("fufacility")
	err := col.Find(nil).Select(bson.M{"facilityName":1,"phone":1,"url":1,"email":1,"address":1}).All(&s)
	if err != nil{
		logrus.Println(err)
	}
	g.JSON(http.StatusOK,s)
}