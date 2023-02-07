package gstmt

// import (
// 	"fmt"
// 	"gosql/gdo"
// 	"gosql/stmt"
// 	"sync"
// 	"tool/cntr"

// 	"github.com/pkg/errors"
// 	"google.golang.org/protobuf/reflect/protoreflect"
// )

// var null protoreflect.ProtoMessage = nil

// type IGBA interface {
// 	Init() error
// 	HasBuffer() bool
// 	Execute() error
// }

// type GBA[PM protoreflect.ProtoMessage] struct {
// 	sgl    *GoSql
// 	gid    byte
// 	tid    byte
// 	datas  map[int64]PM
// 	buffer map[stmt.DbOp]*cntr.Set[int64]
// 	mu     sync.Mutex
// 	id     int64

// 	//////////////////////////////////////////////////
// 	// 客製化函式
// 	//////////////////////////////////////////////////
// 	genFunc   func() PM
// 	getIdFunc func(pm PM) int64
// }

// func NewGBA[PM protoreflect.ProtoMessage](gid byte, tid byte, genFunc func() PM, getIdFunc func(pm PM) int64) *GBA[PM] {
// 	g := &GBA[PM]{
// 		gid:   gid,
// 		tid:   tid,
// 		datas: map[int64]PM{},
// 		buffer: map[stmt.DbOp]*cntr.Set[int64]{
// 			stmt.DbDelete: cntr.NewSet[int64](),
// 			stmt.DbInsert: cntr.NewSet[int64](),
// 			stmt.DbUpdate: cntr.NewSet[int64](),
// 		},
// 		genFunc:   genFunc,
// 		getIdFunc: getIdFunc,
// 	}
// 	return g
// }

// func (g *GBA[PM]) Init() error {
// 	var err error
// 	g.sgl, err = GetGoSql(g.gid)
// 	if err != nil {
// 		return errors.Wrapf(err, "取得 GoSql(%d) 時發生錯誤", g.gid)
// 	}
// 	var i, count int32
// 	pms := g.newContainer(100)
// 	count, err = g.sgl.QueryMulti(g.tid, pms, nil)
// 	if err != nil {
// 		return errors.Wrapf(err, "讀取 ProtoTable(%d) 時發生錯誤", g.tid)
// 	}
// 	var id int64
// 	var pm PM
// 	var ok bool
// 	for i = 0; i < count; i++ {
// 		pm = (*pms)[i].(PM)
// 		id = g.getIdFunc(pm)
// 		if _, ok = g.datas[id]; !ok {
// 			g.datas[id] = pm
// 			fmt.Printf("(g *GBA[PM]) Init | id: %d, data: %+v\n", id, pm)
// 		}
// 	}
// 	return nil
// }

// func (g *GBA[PM]) Set(pm PM) {
// 	g.id = g.getIdFunc(pm)
// 	g.datas[g.id] = pm
// }

// func (g *GBA[PM]) Get(id int64) (PM, error) {
// 	if pm, ok := g.datas[id]; ok {
// 		return pm, nil
// 	}
// 	return null.(PM), errors.New(fmt.Sprintf("沒有 id 為 %d 的數據\n", id))
// }

// //////////////////////////////////////////////////
// // Insert
// //////////////////////////////////////////////////
// func (g *GBA[PM]) Insert(pm PM) error {
// 	g.id = g.getIdFunc(pm)
// 	if _, ok := g.datas[g.id]; ok {
// 		return errors.New(fmt.Sprintf("id 為 %d 的數據已存在\n", g.id))
// 	}

// 	g.mu.Lock()
// 	defer g.mu.Unlock()

// 	// 加入數據
// 	g.datas[g.id] = pm
// 	// 加入插入緩存
// 	g.buffer[stmt.DbInsert].Add(g.id)
// 	return nil
// }

// //////////////////////////////////////////////////
// // Query
// //////////////////////////////////////////////////
// func (g *GBA[PM]) Query(id int64) PM {
// 	if pm, ok := g.datas[id]; ok {
// 		return pm
// 	}
// 	return null.(PM)
// }

// func (g *GBA[PM]) Filter(pms *[]PM, filter func(pm PM) bool) (int, error) {
// 	var i, length int = 0, len(*pms)
// 	for _, pm := range g.datas {
// 		if filter(pm) {
// 			(*pms)[i] = pm
// 			i++
// 			if i == length {
// 				break
// 			}
// 		}
// 	}
// 	return i, nil
// }

// func (g *GBA[PM]) newContainer(size int) *[]protoreflect.ProtoMessage {
// 	pms := make([]protoreflect.ProtoMessage, size)
// 	for i := 0; i < size; i++ {
// 		pms[i] = g.genFunc()
// 	}
// 	return &pms
// }

// //////////////////////////////////////////////////
// // Update
// //////////////////////////////////////////////////
// func (g *GBA[PM]) Update(pm PM) error {
// 	g.id = g.getIdFunc(pm)

// 	if _, ok := g.datas[g.id]; !ok {
// 		err := g.Insert(pm)
// 		if err != nil {
// 			return errors.Wrap(err, fmt.Sprintf("沒有 id 為 %d 的數據\n", g.id))
// 		}
// 		return errors.New(fmt.Sprintf("沒有 id 為 %d 的數據\n", g.id))
// 	}

// 	g.mu.Lock()
// 	defer g.mu.Unlock()

// 	// 更新數據
// 	g.datas[g.id] = pm

// 	// Insert 的緩存當中已有，當執行 Insert 時便會將最新數據加入，無須再額外執行 Update
// 	if g.buffer[stmt.DbInsert].Contains(g.id) {
// 		return nil
// 	}

// 	if !g.buffer[stmt.DbUpdate].Contains(g.id) {
// 		// 加入更新緩存
// 		g.buffer[stmt.DbUpdate].Add(g.id)
// 	}

// 	return nil
// }

// //////////////////////////////////////////////////
// // Delete
// //////////////////////////////////////////////////
// func (g *GBA[PM]) Delete(id int64) error {
// 	g.id = id
// 	if _, ok := g.datas[g.id]; !ok {
// 		return errors.New(fmt.Sprintf("沒有 id 為 %d 的數據\n", g.id))
// 	}

// 	needDelete := true
// 	g.mu.Lock()
// 	defer g.mu.Unlock()

// 	if g.buffer[stmt.DbInsert].Contains(g.id) {
// 		g.buffer[stmt.DbInsert].Remove(g.id)
// 		needDelete = false
// 	} else if g.buffer[stmt.DbUpdate].Contains(g.id) {
// 		g.buffer[stmt.DbUpdate].Remove(g.id)
// 	}

// 	delete(g.datas, g.id)

// 	if needDelete {
// 		g.buffer[stmt.DbUpdate].Add(g.id)
// 	}
// 	return nil
// }

// //////////////////////////////////////////////////
// // GBA
// //////////////////////////////////////////////////
// func (g *GBA[PM]) HasBuffer() bool {
// 	if g.buffer[stmt.DbUpdate].Length() > 0 {
// 		return true
// 	}
// 	if g.buffer[stmt.DbInsert].Length() > 0 {
// 		return true
// 	}
// 	if g.buffer[stmt.DbDelete].Length() > 0 {
// 		return true
// 	}
// 	return false
// }

// func (g *GBA[PM]) Execute() error {
// 	if g.buffer[stmt.DbDelete].Length() > 0 {
// 		for uid := range g.buffer[stmt.DbDelete].Elements {
// 			g.sgl.DeleteBy(g.tid, gdo.WS().Eq("player_id", uid))
// 		}
// 		g.buffer[stmt.DbDelete].Clear()
// 	}

// 	if g.buffer[stmt.DbUpdate].Length() > 0 {
// 		for uid := range g.buffer[stmt.DbUpdate].Elements {
// 			g.sgl.Update(g.tid, g.datas[uid], gdo.WS().Eq("player_id", uid))
// 		}
// 		g.buffer[stmt.DbUpdate].Clear()
// 	}

// 	if g.buffer[stmt.DbInsert].Length() > 0 {
// 		for uid := range g.buffer[stmt.DbInsert].Elements {
// 			g.sgl.Insert(g.tid, []protoreflect.ProtoMessage{g.datas[uid]})
// 			g.sgl.ExcuteInsert(g.tid)
// 			g.buffer[stmt.DbInsert].Remove(uid)
// 		}
// 	}

// 	return nil
// }
