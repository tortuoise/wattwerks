package wattwerks

import (
        "appengine"
        "appengine/datastore"
        "fmt"
        "reflect"
        "time"
        _"golang.org/x/net/context"
        _"net/http"
)

type DS struct {
        //r *http.Request
        ctx appengine.Context
}

type DSErr struct {
        When time.Time
        What string
}

var _ CustDatabase = &DS{}
var _ error = &DSErr{}

var (
        entities = map[string]string{"Cust": "cstmr", "Good": "product", "Cart": "cart", "Order": "order"}

)

func (ds *DS) datastoreKey(id int64) (*datastore.Key) {

        c := ds.ctx
        if id > 9999 && id < 99999 {
                return  datastore.NewKey(c, "cstmr", "", id, nil)
        }
        return  datastore.NewKey(c, "good", "", id, nil)

}

func (ds *DS) dsKey(t reflect.Type, id interface{}) (*datastore.Key) {

        c := ds.ctx
        if entity, ok := entities[t.Name()]; ok {
                return  datastore.NewKey(c, entity, "", id.(int64), nil)
        }
        return nil

}

func (ds *DS) ListCusts() ([]*Cust, error) {

        c := ds.ctx
        cs := make([]*Cust, 0)
        q := datastore.NewQuery("cstmr").Order("Id")
        ks, err := q.GetAll(c, &cs)
        if err != nil {
                return nil, fmt.Errorf("Get custy list error")
        }
        for i, k := range ks {
                cs[i].Id = k.IntID()
        }
        return cs, nil

}

func (ds *DS) AddCust(custy *Cust) (int64, error) {

        c := ds.ctx
        k := ds.datastoreKey(custy.Id)
        _, err := datastore.Put(c, k, custy)
        if err != nil {
                return 0, fmt.Errorf("Add custy error")
        }
        return k.IntID(), nil

}

func (ds *DS) Add(v interface{}) (int64, error) {

        c := ds.ctx
        k := ds.dsKey(reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem().Field(0).Interface())
        if k == nil {
                return 0, fmt.Errorf("Add custy error - key create error")
        }
        _, err := datastore.Put(c, k, v)
        if err != nil {
                return 0, fmt.Errorf("Add custy error - datastore put error")
        }
        return k.IntID(), nil

}

func (ds *DS) GetCust(id int64) (*Cust, error) {

        c := ds.ctx
        k := ds.datastoreKey(id)
        cst := &Cust{}
        err := datastore.Get(c, k, cst)
        if err != nil {
                return nil, fmt.Errorf("Get by id error")
        }
        return cst, nil

}

func (ds *DS) GetCustwEmail(email string) (*Cust, error)  {

        c := ds.ctx
        q := datastore.NewQuery("cstmr").Filter("Email", email)
        cst := &Cust{}
        _, err := q.GetAll(c, cst)
        if err != nil {
                return nil, fmt.Errorf("Get by email error")
        }
        return cst, nil

}

func (ds *DS) GetCustKey(email string) (*Cust, *datastore.Key, error) {

        c := ds.ctx
        q := datastore.NewQuery("cstmr").Filter("Email", email)
        cst := &Cust{}
        k, err := q.GetAll(c, cst)
        if err != nil {
                return nil, nil, fmt.Errorf("Get by email error %v", err)
        }
        return cst, k[0], nil

}

func (ds *DS) Get(v interface{}) error {

        if reflect.TypeOf(v).Kind() != reflect.Ptr {
                //return fmt.Errorf("Get error: pointer reqd. as argument")
                return DSErr{When: time.Now(), What: "Get error: pointer reqd",}
        }
        c := ds.ctx
        id := reflect.ValueOf(v).Elem().Field(0).Interface()
        if id.(int64) == 0 {
                return fmt.Errorf("Get error - id not set")
        }
        k := ds.dsKey(reflect.TypeOf(v).Elem(), id)
        if k == nil {
                return fmt.Errorf("Get error - key create error")
        }

        if err := datastore.Get(c, k, reflect.ValueOf(v).Interface()); err != nil {
                return fmt.Errorf("Get error - datastore get error")
        }
        return nil

}

func (ds *DS) UpdateCust(custy *Cust) error {

        c := ds.ctx
        k := ds.datastoreKey(custy.Id)
        _, err := datastore.Put(c, k, custy)
        if err != nil {
                return fmt.Errorf("Updating error %v", err)
        }
        return nil

}

func (ds *DS) DeleteCust(id int64) error {

        c := ds.ctx
        k := ds.datastoreKey(id)
        if err := datastore.Delete(c, k); err != nil {
                return fmt.Errorf("Deleting error")
        }
        return nil

}

func (ds *DS) Close() error {

        return nil

}

func (dse DSErr) Error() string {

        return fmt.Sprintf("%v: %v", dse.When, dse.What)

}
