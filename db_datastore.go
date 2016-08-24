package wattwerks

import (
        "appengine"
        "appengine/datastore"
        "fmt"
        "log"
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
        entities = map[string]string{"Counter": "counter", "Cust": "cstmr", "Good": "product", "GoodDeets": "productdeets",  "Cart": "cart", "Order": "order"}

)

func (ds *DS) datastoreKey(id int64) (*datastore.Key) {

        c := ds.ctx
        if id > 9999 && id < 99999 {
                return  datastore.NewKey(c, "cstmr", "", id, nil)
        }
        return  datastore.NewKey(c, "good", "", id, nil)

}

func (ds *DS) dsKey(t reflect.Type, id ...interface{}) (*datastore.Key) {

        c := ds.ctx
        if entity, ok := entities[t.Name()]; ok {
                switch t.Name() {
                        case "Counter":
                                return  datastore.NewKey(c, entity, "thekey", 0, nil)
                        default:
                                if len(id) > 0 {
                                        return  datastore.NewKey(c, entity, "", id[0].(int64), nil)
                                } else {
                                        return  datastore.NewIncompleteKey(c, entity, nil)  // shouldn't get here
                                }
                }
        }
        return nil

}

func (ds *DS) dsChildKey(t reflect.Type, id int64, pk *datastore.Key) *datastore.Key {

        c := ds.ctx
        if entity, ok := entities[t.Name()]; ok {
                return  datastore.NewKey(c, entity, "", id, pk)
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

/*func (ds *DS) Add(v interface{}) (int64, error) {

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

}*/

func (ds *DS) Add(v interface{}, n ...int64) (int64, error) {

        c := ds.ctx
        var k *datastore.Key
        if len(n) == 0 { // for Counter
                //k = ds.dsKey(reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem().Field(0).Interface())
                k = ds.dsKey(reflect.TypeOf(v).Elem())
        } else { // if extra args are provided use as key ID
                k = ds.dsKey(reflect.TypeOf(v).Elem(), n[0])
        }
        if k == nil {
                return 0, fmt.Errorf("Add error - key create error - unknown entity")
        }
        _, err := datastore.Put(c, k, v)
        if err != nil {
                return 0, fmt.Errorf("Add error - datastore put error: %v", err)
        }
        return k.IntID(), nil

}

func (ds *DS) AddwParent(parent interface{}, child interface{}, offset int64) (error) {

        if reflect.TypeOf(parent).Kind() != reflect.Ptr || reflect.TypeOf(child).Kind() != reflect.Ptr{
                return DSErr{When: time.Now(), What: "Get error: pointers reqd",}
        }

        c := ds.ctx

        pt := reflect.TypeOf(parent).Elem()
        if _, ok := pt.FieldByName("Id"); !ok {
                return DSErr{When: time.Now(), What: "Add w parent error: parent lacks Id",}
        }

        pv := reflect.ValueOf(parent).Elem().Field(0).Interface().(int64)
        log.Println(pv)
        pk := ds.dsKey(pt, pv)
        ck := ds.dsChildKey(reflect.TypeOf(child).Elem(), pv + offset, pk)

        if pk == nil || ck == nil {
                return DSErr{When: time.Now(), What: "Add w parent error: during key creation",}
        }
        pk, err := datastore.Put(c, pk, parent)
        if err != nil {
                return DSErr{When: time.Now(), What: "Add w parent error: during parent put",}
        }
        if pk != nil {
                _, err := datastore.Put(c, ck, child)
                if err != nil{
                        return DSErr{When: time.Now(), What: "AddwParent error: during child put"}
                }
        }
        return nil
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
                return DSErr{When: time.Now(), What: "Get error: pointer reqd",}
        }

        var id int64
        var k *datastore.Key
        // check whether Id field is available in struct - if it is, it shouldn't be 0
        if _, ok := reflect.TypeOf(v).Elem().FieldByName("Id"); ok {
                id = reflect.ValueOf(v).Elem().Field(0).Interface().(int64)  // could also use FieldByName("Id") instead of Field(0)
                if id == 0 { // shouldn't be zero
                        return DSErr{When: time.Now(), What: "Get error: id not set",}
                }
                k = ds.dsKey(reflect.TypeOf(v).Elem(), id)  //complete key
        } else {
                k = ds.dsKey(reflect.TypeOf(v).Elem())
        }

        c := ds.ctx
        if k == nil {
                return fmt.Errorf("Get error - key create error")
        }

        //if err := datastore.Get(c, k, reflect.ValueOf(v).Interface()); err != nil {
        if err := datastore.Get(c, k, v); err != nil {
                return fmt.Errorf("Get error - datastore get error: %v, key kind: %v",err, k.Kind())
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
