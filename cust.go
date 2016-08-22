package wattwerks

import (
        "appengine/datastore"
        "time"
)

type Cust struct {
        Id int64 `json:"id"` // 10,000-99,999
        Email string `json:"email"`
        Password string `json:"password"`
        Firstname string `json:"firstname"`
        Lastname string `json:"lastname"`
        Birthdate time.Time `json:"birthdate"`
        Meterid int64 `json:"meterid"`
        Addr string `json:"addr"`
        Street string `json:"street"`
        Area string `json:"area"`
        City string `json:"city"`
        Postcode string `json:"postcode"`
        State string `json:"state"`
        Country string `json:"country"`
        Phone string `json:"phone"`
        Indi bool `json:"indi"` //Individual or not (company)
        Company string `json:"company"`
        Vat string `json:"vat"`
        Pmtmethod string `json:"pmtmethod"`
        Crtid int64 `json:"cartid"` //current cart 100,000-999,999
        //Crt Cart `json:"cart"`
        //Rdrs []Order `json:"rdrs"`
}

type CustDatabase interface {

        ListCusts() ([]*Cust, error)

        AddCust(custy *Cust) (int64, error)   //create

        GetCust(id int64) (*Cust, error)  //retrieve by id

        GetCustwEmail(email string) (*Cust, error)  //retrieve by email

        GetCustKey(email string) (*Cust, *datastore.Key, error)

        UpdateCust(custy *Cust) error //update

        DeleteCust(id int64) error //delete

        Close() error
}
