package gen

import (
	"context"
	"errors"
	"fmt"
	"id-gen-db/client"
	"sync"
	"time"
)

type mallocResult interface {
	GetMax() int64
	GetCur() int64
}

type buffer struct {
	cur int64
	max int64
}

func (f buffer) GetMax() int64 {
	return f.max
}
func (f buffer) GetCur() int64 {
	return f.cur
}

type IDGen struct {
	locker       sync.Mutex
	cur          int64
	max          int64
	buffer       *buffer
	mallocLocker sync.Mutex
	inMalloc     bool
}

var idGen *IDGen

func GetIDGen() *IDGen {
	return idGen
}

func GenInit() {
	gen := &IDGen{}
	err := gen.malloc()
	if err != nil {
		panic(err)
	}
	idGen = gen
}

func (g *IDGen) NextID() (int64, error) {
	g.locker.Lock()
	defer g.locker.Unlock()
	o := time.After(2 * time.Millisecond)
	d := make(chan struct{})
	go func() {
		select {
		case <-o:
			fmt.Println("超时 gen")
			fmt.Printf("%v\n", g)
		case <-d:
		}
	}()

	if g.cur == g.max {
		if g.buffer == nil {
			return 0, errors.New("wrong buffer")
		}
		g.cur = g.buffer.cur
		g.max = g.buffer.max
		g.buffer = nil
	}
	if g.max-g.cur < 1000 && g.buffer == nil {
		go g.malloc()
	}
	g.cur++
	close(d)
	return g.cur, nil
}

func (g *IDGen) malloc() error {
	if g.inMalloc {
		return nil
	}
	max := g.max
	fmt.Printf("start malloc :%v\n", max)
	defer func() {
		fmt.Printf("end malloc :%v\n", max)
	}()
	g.mallocLocker.Lock()
	defer g.mallocLocker.Unlock()
	if g.buffer != nil {
		return nil
	}
	g.inMalloc = true
	defer func() { g.inMalloc = false }()
	buf, err := mallocFromMysql(context.Background(), "test")

	if err != nil {
		return err
	}
	g.buffer = &buffer{buf.GetCur(), buf.GetMax()}
	return nil
}

func mallocFromMysql(ctx context.Context, bus string) (mallocResult, error) {
	db := client.GetMysqlDb()

	value := []interface{}{}
	value = append(value, bus)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `update gen set max=max+step where bus = ?`, value...)
	if err != nil {
		return nil, err
	}
	row := tx.QueryRowContext(ctx, `select max, (max - step) as cur from gen where bus = ?`, value...)
	var max, cur int64
	err = row.Scan(&max, &cur)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	return buffer{cur, max}, nil
}
