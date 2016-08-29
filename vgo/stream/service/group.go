package service

import (
	"encoding/json"
	"fmt"
	//	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var hostGroups = make(map[string][]*Group)

type Group struct {
	ID     string
	Alerts map[string]*Alert
	Parent *Group
	Users  []*User
	Hosts  []string
}

type GroupRaw struct {
	ID     string
	Alerts map[string]*Alert
	Parent string
	Users  []*User
	Hosts  []string
}

type Alert struct {
	Type  int
	Value [2]float64
	Count int
}

type User struct {
	Type       uint8
	Phone      string
	Operatorid string
}

var rawGroups = make(map[string]*GroupRaw)
var allGroup = make(map[string]*Group)

var db *bolt.DB

func main() {
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initGroups()

	writeGroups()

	loadGroups()

	initHostGroups()

	buf, err := json.Marshal(allGroup)
	fmt.Println(err, string(buf))
}

func initHostGroups() {
	//	for k, v := range allGroup {

	//	}

}

func loadGroups() {
	//first load into rawGroups
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("groups"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			g := &GroupRaw{}
			json.Unmarshal(v, g)
			rawGroups[string(k)] = g
		}

		return nil
	})

	// load rawGroups into allGroup
	for k, v := range rawGroups {
		fmt.Println(k)
		// 初始化该group和parent线上的所有group
		firstG, ok := allGroup[k]

		// 当前group一旦初始化过,直接返回
		if ok {
			continue
		}
		firstG = &Group{}
		firstG.Alerts = v.Alerts
		firstG.Users = v.Users
		firstG.ID = v.ID
		allGroup[k] = firstG

		parent := v.Parent
		g := firstG

		for {
			// 如果没有parent，则返回
			if parent == "" {
				break
			}

			// 寻找父亲
			pr1, _ := rawGroups[parent]
			pp1, ok := allGroup[parent]

			// 若第一个父亲已经存在，指针赋值后直接返回
			if ok {
				g.Parent = pp1
				break
			}

			pp1 = &Group{}
			pp1.Alerts = pr1.Alerts
			pp1.Users = pr1.Users
			pp1.ID = pr1.ID

			g.Parent = pp1

			allGroup[parent] = pp1

			parent = pr1.Parent
			g = pp1
		}

	}
}
func writeGroups() {
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("groups"))
		if err != nil {
			log.Fatalln("create bucket error: ", err)
		}

		for k, v := range rawGroups {
			bv, _ := json.Marshal(v)
			b.Put([]byte(k), bv)
		}

		return nil
	})
}
func initGroups() {
	u := &User{
		Phone:      "15880261185",
		Operatorid: "994757521",
	}

	alert0 := &Alert{
		Type:  1,
		Value: [2]float64{60, 80},
		Count: 10,
	}

	group0 := &GroupRaw{
		ID:     "group0",
		Alerts: make(map[string]*Alert),
		Parent: "",
		Users:  []*User{u},
	}

	group0.Alerts["mem.mem_usage"] = alert0

	alert1 := &Alert{
		Type:  1,
		Value: [2]float64{60, 80},
		Count: 10,
	}

	group1 := &GroupRaw{
		ID:     "group1",
		Alerts: make(map[string]*Alert),
		Parent: "group0",
		Users:  []*User{u},
	}

	group1.Alerts["cpu.cpu_usage"] = alert1

	alert2 := &Alert{
		Type:  1,
		Value: [2]float64{50, 70},
		Count: 10,
	}

	group2 := &GroupRaw{
		ID:     "group2",
		Alerts: make(map[string]*Alert),
		Parent: "group1",
		Users:  []*User{u},
	}

	group2.Alerts["cpu.cpu_usage"] = alert2

	alert3 := &Alert{
		Type:  1,
		Value: [2]float64{40, 60},
		Count: 10,
	}

	group3 := &GroupRaw{
		ID:     "group3",
		Alerts: make(map[string]*Alert),
		Parent: "group0",
		Users:  []*User{u},
	}

	group3.Alerts["cpu.cpu_usage"] = alert3

	alert4 := &Alert{
		Type:  1,
		Value: [2]float64{40, 60},
		Count: 10,
	}

	group4 := &GroupRaw{
		ID:     "group4",
		Alerts: make(map[string]*Alert),
		Parent: "group2",
		Users:  []*User{u},
	}

	group4.Alerts["cpu.cpu_usage"] = alert4

	rawGroups["group0"] = group0
	rawGroups["group1"] = group1
	rawGroups["group2"] = group2
	rawGroups["group3"] = group3
	rawGroups["group4"] = group4
}
