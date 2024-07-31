package gotil

type Data struct {
	Id int `json:"id"`
}

func (*Data) TableName() string {
	return "data"
}

func main() {
	c, _ := NewElasticSearch("")
	c.Save(&Data{})

	var datas []Data
	c.Search(&datas, "", "", "")
}
