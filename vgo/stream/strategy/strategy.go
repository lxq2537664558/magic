package strategy

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

var strategyes *Strategy

type Strategy struct {
	// grlock    sync.RWMutex
	// aglock    sync.RWMutex
	sync.RWMutex
	GroupRaws map[string]*GroupRaw
	AllGroup  map[string]*Group
	db        *bolt.DB
	dbname    string
	// bucket     *bolt.Bucket
	bucketName []byte
}

func (sg *Strategy) Init() {
	db, err := bolt.Open(sg.dbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	sg.db = db
	sg.GroupRaws = make(map[string]*GroupRaw)
	sg.AllGroup = make(map[string]*Group)
	// if groups exists , loadGroups
	if sg.isExist() {
		sg.loadGroups()
	}
}

func (sg *Strategy) isExist() bool {
	err := sg.db.View(func(tx *bolt.Tx) error {
		if tx.Bucket(sg.bucketName) == nil {
			return fmt.Errorf("%s groups_not_exist", string(sg.bucketName))
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

func (sg *Strategy) loadGroups() error {
	if !sg.isExist() {
		return fmt.Errorf("groups_not_exist")
	}

	sg.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(sg.bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			g := &GroupRaw{}
			json.Unmarshal(v, g)
			log.Printf("loadGroups group message is %s, %v\n", string(k), g)
			sg.Lock()
			sg.GroupRaws[string(k)] = g
			sg.Unlock()
		}
		sg.Lock()
		// load rawGroups into allGroup
		for k, v := range sg.GroupRaws {
			fmt.Println(k, v)
			// 初始化该group和parent线上的所有group
			firstG, ok := sg.AllGroup[k]
			// 当前group一旦初始化过,直接返回
			if ok {
				continue
			}
			firstG = &Group{}
			firstG.Alerts = v.Alerts
			firstG.Users = v.Users
			firstG.ID = v.ID
			sg.AllGroup[k] = firstG
			parent := v.Parent
			g := firstG
			for {
				// 如果没有parent，则返回
				if parent == "" {
					break
				}
				// 寻找父亲
				pr1, ok := sg.GroupRaws[parent]
				if !ok {
					log.Fatal("can find parent, parent name is ", parent)
					return fmt.Errorf("can find parent, parent name is %s ", parent)
				}
				pp1, ok := sg.AllGroup[parent]
				// 若第一个父亲已经存在，指针赋值后直接返回
				if ok {
					// 增加子节点
					pp1.AddChild(g.ID)
					g.Parent = pp1
					break
				}
				pp1 = &Group{}
				pp1.Alerts = pr1.Alerts
				pp1.Users = pr1.Users
				pp1.ID = pr1.ID
				g.Parent = pp1
				sg.AllGroup[parent] = pp1
				parent = pr1.Parent
				// 增加子节点
				pp1.AddChild(g.ID)
				g = pp1
			}
		}
		sg.Unlock()
		//

		// free mem
		sg.Lock()
		// load rawGroups into allGroup
		for k, _ := range sg.GroupRaws {
			delete(sg.GroupRaws, k)
		}
		sg.Unlock()
		return nil
	})
	return nil
}

func (sg *Strategy) Show() {
	for k, v := range sg.AllGroup {
		fmt.Println("Show ->>> ", k, v.Parent)
	}
}

func (sg *Strategy) Add(gr *GroupRaw) error {
	return sg.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(sg.bucketName)
		if err != nil {
			return fmt.Errorf("create bucket error: %s", err)
		}
		g, err := gr.GetNewGroup(gr)
		if err != nil {
			log.Println("Add group ", gr.ID, "failed, err message is", err)
			return err
		}
		sg.Lock()
		sg.AllGroup[gr.ID] = g
		sg.Unlock()
		bv, err := json.Marshal(gr)
		if err != nil {
			return err
		}
		return b.Put([]byte(gr.ID), bv)
	})
}

func (sg *Strategy) delete(gid string, bucket *bolt.Bucket) {
	// delete group and child group
	for {
		if g, ok := sg.AllGroup[gid]; ok {
			for _, child := range g.Child {
				sg.delete(child, bucket)
			}
			delete(sg.AllGroup, gid)
			bucket.Delete([]byte(gid))
		}
		break
	}
}

func (sg *Strategy) Del(gid string) error {
	return sg.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(sg.bucketName)
		if bucket == nil {
			return nil
		}
		sg.Lock()
		sg.delete(gid, bucket)
		sg.Unlock()
		// return bucket.Delete([]byte(gid))
		return nil
	})
}

func (sg *Strategy) Close() error {
	return sg.db.Close()
}

func NewStrategy(dbname string, bucketname string) *Strategy {
	sg := &Strategy{dbname: dbname, bucketName: []byte(bucketname)}
	strategyes = sg
	return sg
}

// test test
func test() {

	strategyes = NewStrategy("my.db", "groups")
	strategyes.Init()
	fmt.Println("Alert", strategyes)

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

	group1 := &GroupRaw{
		ID:     "group1",
		Alerts: make(map[string]*Alert),
		Parent: "group0",
		Users:  []*User{u},
	}
	group1.Alerts["mem.mem_usage"] = alert0

	group2 := &GroupRaw{
		ID:     "group2",
		Alerts: make(map[string]*Alert),
		Parent: "group1",
		Users:  []*User{u},
	}
	group2.Alerts["mem.mem_usage"] = alert0

	err := strategyes.Add(group0)
	if err != nil {
		log.Println("Add group0 err ", err)
	}
	err = strategyes.Add(group1)
	if err != nil {
		log.Println("Add group1 err ", err)
	}
	err = strategyes.Add(group2)
	if err != nil {
		log.Println("Add group2 err ", err)
	}
	time.Sleep(1 * time.Second)

	strategyes.Show()

	fmt.Println("Del Start")
	err = strategyes.Del("group0")
	if err != nil {
		log.Println("Del group0 err ", err)
	}
	fmt.Println("Del End")

	strategyes.Show()
	log.Println("Del group0 ok")

	err = strategyes.Del("group1")
	if err != nil {
		log.Println("Del group1 err ", err)
	}
	log.Println("Del group1 ok")

	err = strategyes.Del("group2")
	if err != nil {
		log.Println("Del group2 err ", err)
	}
	log.Println("Del group2 ok")
	// time.Sleep(5 * time.Second)

	strategyes.Close()
}
