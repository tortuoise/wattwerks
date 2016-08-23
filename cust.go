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

type Cart struct {
        Id int64 `json:"id"` // 100,000-999,999
        Cstid int64 `json:"cstid"`  //completed cart default 0
        Ids []int64 `json:"items"` //Good Ids
        Qntts []int `json:"qntts"`
        Prchsd bool `json:"purchased"`
        Ttls []float64 `json:"ttls"`
        Txs float64 `json:"txs"`
        Ttl float64 `json:"total"`
}

type Order struct{
        Id int64 `json:"id"` //1,000,000-9,999,999
        TxId string `json:"txid"`
        TxRefNo string `json:"txrefno"`
        TxStatus string `json:"txstatus"`
        TxMsg string `json:"txmsg"`
        PgTxnNo string `json:"pgtxnno"`
        IssuerRefNo string `json:"issuerrefno"`
        AuthIdCode string `json:"authidcode"`
        PgRespCode string `json:"pgrespcode"`
        OriginalAmount string `json:"originalamount "`
        AdjustedAmount string `json:"adjustedamount "`
        DpRuleName string `json:"dprulename "`
        DpRuleType string `json:"dpruletype "`
        Amount string `json:"amount "`
        TransactionAmount string `json:"transactionamount "`
        PaymentMode string `json:"paymentmode "`
        TxGateway string `json:"txgateway "`
        IssuerCode string `json:"issuercode "`
        TxnDateTime string `json:"txndatetime "`
        IsCOD string `json:"iscod "`
        //Details of un-registered buyer
        FirstName string `json:"firstname "`
        LastName string `json:"lastname "`
        AddressStreet1 string `json:"addressstreet1 "`
        AddressStreet2 string `json:"addressstreet2 "`
        AddressFirstName string `json:"addressfirstname "`
        AddressZip string `json:"addresszip "`
        AddressState string `json:"addressstate "`
        AddressCountry string `json:"addresscountry "`
        MobileNo string `json:"mobileno "`
        Email string `json:"email "`
        //user account if buyer registered
        //Buyer *Cust //`json:"buyer"`
        //Cart Cart `json:"cart"`
        ReqTime time.Time `json:"reqTime"`
        SecSign string `json:"secSign"`
        Cstid int64 `json:"cstid"`  //completed cart buyer
        Crtid int64 `json:"cartid"` //completed cart id
        Notified bool `json:"boolean"`
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
