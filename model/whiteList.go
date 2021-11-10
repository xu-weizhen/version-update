//处理设备白名单
package model

import (
	"container/list"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

const MAX_MAP_NUM = 10
const MAX_ID_NUM = 100000

var lock_of_list sync.Mutex //互斥锁1
var lock_of_load sync.Mutex

type Index struct {
	Aid                 int
	Device_platform     string
	Update_version_code string
}

type IDList struct { //加载到内存的白名单
	Head    Index               //白名单的表头
	Content map[string]struct{} //白名单内容
}

var Id_list = list.New()

//var hs map[Index]*list.Element
var Hash sync.Map

func check_ID(v *Version, res *NewVersion, db *sql.DB) bool {
	temp := Index{
		v.Aid,
		v.Device_platform,
		res.Update_version_code,
	}
	fmt.Printf("temp的内容为%v", temp.Update_version_code)
	fmt.Printf("id的内容为%v", v.Device_id)
	_, exist := Hash.Load(temp)
	if exist { //版本对应的白名单已经加载到内存中
		return is_Here(temp, v)
	} else {
		//首先再次判断是否在内存中
		lock_of_load.Lock()
		fmt.Println("执行到了这里1")
		if _, exist := Hash.Load(temp); exist { //在内存中，则解锁，再执行is_here函数查询设备id是否存在
			lock_of_load.Unlock()
			fmt.Println("执行到了这里2")
			return is_Here(temp, v)
		} else { //不在内存，需要从数据库中加载device_id_list到内存中
			//首先从数据库加载device_id_list
			stmt, _ := db.Prepare(`select device_id_list from  device_id where aid=? AND platform=? AND update_version_code=?`)
			defer stmt.Close()
			var str string
			stmt.QueryRow(temp.Aid, temp.Device_platform, temp.Update_version_code).Scan(&str)
			//fmt.Printf("%v", str)
			strlist := strings.Split(str, " ")
			fmt.Printf("长度是：%v", len(strlist))

			//将获取到的白名单list写入一个map数组中
			Idlist_map := make(map[string]struct{}, 100000)
			for i := 0; i < 100000; i++ {
				Idlist_map[strlist[i]] = struct{}{}
			}

			Node := IDList{
				temp,
				Idlist_map,
			}
			//先判断当前链表是否满了吗(10个)，满了就除去链表头的map，再向链表尾部插入新的map
			lock_of_list.Lock()
			if Id_list.Len() < 2 {
				Id_list.PushBack(Node)
				//hs[temp] = Id_list.Back()
				Hash.Store(temp, Id_list.Back())
				fmt.Println("执行到了这里3")
			} else {
				//delete(hs, Id_list.Front().Value.(IDList).Head) //删去索引
				fmt.Println("执行到了这里4")
				Hash.Delete(Id_list.Front().Value.(IDList).Head)
				Id_list.Remove(Id_list.Front())
				Id_list.PushBack(Node)
				//hs[temp] = Id_list.Back()
				Hash.Store(temp, Id_list.Back())
				//怎么删除map上的元素
			}

			lock_of_list.Unlock()

			lock_of_load.Unlock()
			//判断在不在
			if _, t := Idlist_map[v.Device_id]; t {
				fmt.Println("执行到了这里5")
				return true
			} else {
				fmt.Println("执行到了这里6")
				return false
			}
		}
	}
	fmt.Println("执行到了这里7")
	return false
}

func is_Here(temp Index, v *Version) bool {
	//_, ok := hs[temp].Value.(IDList).Content[v.Device_id]
	t1, _ := Hash.Load(temp)
	_, ok := t1.(*list.Element).Value.(IDList).Content[v.Device_id]
	if ok {
		fmt.Println("执行到了这里8")
		lock_of_list.Lock()
		//Id_list.MoveToBack(hs[temp])
		t2, _ := Hash.Load(temp)
		Id_list.MoveToBack(t2.(*list.Element)) //将map列表移动到链表末尾
		lock_of_list.Unlock()
		return true //在白名单中
	} else {
		fmt.Println("执行到了这里9")
		return false
	}
}
