package wattwerks
import (
        _"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	_"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"regexp"
	"appengine"
	"appengine/datastore"
	"appengine/blobstore"
	//"appengine/log"
	//"appengine/user"
	"appengine/mail"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"strconv"
	"strings"
	"net/url"
)
// datastore structs Good-4 Cust-5 Cart-6 Order-7 GoodDeets-8 CartsPurchased-9
        type Sheep struct {
                Name string `json:"name"`
                Email string `json:"email"`
                Addr string `json:"addr"`
                Phon string `json:"phon"`
                Area string `json:"area"`
                Note string `json:"note"`
                JDte time.Time `json:"jdte"`
                Confirmed bool `json:"conf"`
        }
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
                /*Details of un-registered buyer*/
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
                /*user account if buyer registered*/
                //Buyer *Cust //`json:"buyer"`
                //Cart Cart `json:"cart"`
                ReqTime time.Time `json:"reqTime"`
                SecSign string `json:"secSign"`
                Cstid int64 `json:"cstid"`  //completed cart buyer
                Crtid int64 `json:"cartid"` //completed cart id
                Notified bool `json:"boolean"`
        }
        type Counter struct {
                TtlGd int64 // 1000-9999
                NxtGd int64 // next lowest available ID
                TtlCst int64 // 10,000-99,999
                NxtCst int64
                TtlRdr int64
                NxtRdr int64
                HlsRdr []int64 // incomplete order ids that get returned from PG wo completion/confirmation
        }
        type Category struct {
                Name string `json:"category"`
                Subcategories []string `json:"subcategories"`
        }
        type Good struct {
                Id int64 `json:"id,string"` //1,000-9,999
                Code string `json:"code"`
                Category string `json:"category"`
                Subcategory string `json:"subCategory"`
                Brand string `json:"brand"`
                Desc string `json:"desc"`
                Price float64 `json:"price,string"`
                Url string `json:"url"`
                Urlimg string `json:"urlImg"`
                Featured bool `json:"featured,string"`
                Hidden bool `json:"hidden,string"`
                Deets GoodDeets `json:"goodDeets"`
        }
        type GoodDeets struct {
                //Id int64 `json:"id"`
                DescDetails string `json:"descDetails"`//path to file with details
                Tax float64 `json:"tax,string"`//percent
                Price float64 `json:"price,string"`
                Stock int `json:"stock,string"`
                Related []int64 `json:"related,string"`
                Prices []float64 `json:"prices,string"`
                Volumes []int `json:"volumes,string"`
                //PriceVolume map[int]float64 `json:"priceVolume"`
                ParameterNames []string `json:"parameterNames"`
                ParameterValues []string `json:"parameterValues"`
                //Parameters map[string]string `json:"parameters"`
                Features []string `json:"features"`
                Items []string `json:"items"`
                UrlImgs1 string `json:"urlImgs1"`
                UrlImgs2 string `json:"urlImgs2"`
                UrlImgs3 string `json:"urlImgs3"`
                UrlFile string `json:"urlFile"`
        }
//rendering structs
        type Render struct { //for most purposes
                Message string `json:"message"`
                Cstmr Cust `json:"cstmr"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
        }
        type Render1 struct { //for account lists
                Message string `json:"message"`
                Cstmr Cust `json:"cstmr"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
                Cstmrs []Cust `json:"cstmrs"`
        }
        type Render2 struct { //for upload url
                UpURL *url.URL `json:"upurl,string"`
                Cstmr Cust `json:"cstmr"`
                Categories []Category `json:"categories,string"`
        }
        type Render3 struct { //for order lists
                Message string `json:"message"`
                Cstmr Cust `json:"cstmr"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
                Crt Cart `json:"items"`
                Orders []Order `json:"orders"`
        }
        type Render4 struct { //for cart items
                Message string `json:"message"`
                Cstmr Cust `json:"cstmr"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
                Crt Cart `json:"items"`
        }
        type Render40 struct { //for new order
                Message string `json:"message"`
                Cstmr Cust `json:"cstmr"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
                Crt Cart `json:"items"`
                Rdr Order `json:"rdr"`
        }
        type Render5 struct { //for catalog templates
                Ctgry string `json:"ctgry"`
                Sctgry string `json:"sctgry"`
                Cstmr Cust `json:"cstmr"`
                Crt Cart `json:"crt"`
                Goods []Good `json:"goods,string"`
                Categories []Category `json:"categories,string"`
                Subcategories []string `json:"subcategories,string"`
        }
        type Render6 struct { //for catalog templates
                Ctgry string `json:"ctgry"`
                Sctgry string `json:"sctgry"`
                Cstmr Cust `json:"cstmr"`
                Crt Cart `json:"crt"`
                Gd Good `json:"good,string"`
                Categories []Category `json:"categories,string"`
                Deets GoodDeets `json:"gooddeets"`
        }
//globals
var(
	ckstore = sessions.NewCookieStore([]byte("wttwrks_cks_scrt"))
	tmpl_prm = template.Must(template.ParseGlob("templates/prm/*"))
	tmpl_cmn = template.Must(template.ParseGlob("templates/cmn/*"))
	tmpl_adm_cmn = template.Must(template.ParseGlob("templates/adm/cmn/*"))
	tmpl_adm_gds_edt = template.Must(template.ParseFiles("templates/adm/goods_edit", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_gds_ntr = template.Must(template.ParseFiles("templates/adm/goods_entry", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_gds_pld = template.Must(template.ParseFiles("templates/adm/goods_upload", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_gds_ntr_new = template.Must(template.ParseFiles("templates/adm/goods_entry_new", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_gds_lst = template.Must(template.ParseFiles("templates/adm/goods_list", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_acc_lst = template.Must(template.ParseFiles("templates/adm/acc_list", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_adm_rdrs = template.Must(template.ParseFiles("templates/adm/orders", "templates/adm/cmn/body", "templates/adm/cmn/right", "templates/adm/cmn/center", "templates/adm/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc = template.Must(template.ParseFiles("templates/acc/account", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_reg = template.Must(template.ParseFiles("templates/acc/register", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_lgn = template.Must(template.ParseFiles("templates/acc/login", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_frgt = template.Must(template.ParseFiles("templates/acc/forgot", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_rmnd = template.Must(template.ParseFiles("templates/acc/reminder", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_edt = template.Must(template.ParseFiles("templates/acc/edit", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_rtrns = template.Must(template.ParseFiles("templates/acc/returns", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_rdrs = template.Must(template.ParseFiles("templates/acc/orders", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_crt = template.Must(template.ParseFiles("templates/acc/cart", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_chckt = template.Must(template.ParseFiles("templates/acc/chckt", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_lgn_err = template.Must(template.ParseFiles("templates/acc/lgn_err", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_reg_err = template.Must(template.ParseFiles("templates/acc/reg_err", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_edt_err = template.Must(template.ParseFiles("templates/acc/edt_err", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_acc_rnd_err = template.Must(template.ParseFiles("templates/acc/rnd_err", "templates/acc/cmn/body", "templates/acc/cmn/right", "templates/acc/cmn/center", "templates/acc/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_gds_lst = template.Must(template.ParseFiles("templates/cat/goods_list", "templates/cat/cmn/body", "templates/cat/cmn/top", "templates/cat/cmn/left", "templates/cat/cmn/right", "templates/cat/cmn/center", "templates/cat/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_cat_cat = template.Must(template.ParseFiles("templates/cat/ctgry", "templates/cat/cmn/body", "templates/cat/cmn/top", "templates/cat/cmn/left", "templates/cat/cmn/right", "templates/cat/cmn/center", "templates/cat/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_cat_sub_cat = template.Must(template.ParseFiles("templates/cat/sctgry", "templates/cat/cmn/body", "templates/cat/cmn/top", "templates/cat/cmn/left", "templates/cat/cmn/right", "templates/cat/cmn/center", "templates/cat/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	tmpl_cat_gds_dts = template.Must(template.ParseFiles("templates/cat/gds_dts", "templates/cat/cmn/body", "templates/cat/cmn/top", "templates/cat/cmn/left", "templates/cat/cmn/right", "templates/cat/cmn/center", "templates/cat/cmn/search", "templates/cmn/base", "templates/cmn/head", "templates/cmn/menu", "templates/cmn/footer"))
	validEmail = regexp.MustCompile("^.*@.*\\.(com|org|in|mail|io)$")
	validPath = regexp.MustCompile(`^/(reg|ta|how|flock|confirm|product|goods)?/?(.*)$`)
	cksnbl = []byte(`Please enable cookies`)// to &lt; a href=\"/account/\"&gt; Continue &lt; /a &gt;`)
	acc_errs = map[int]*template.Template{1:tmpl_acc_lgn_err, 2:tmpl_acc_reg_err, 3:tmpl_acc_edt_err, 4:tmpl_acc_rnd_err}
	//gd_errs = map[int]*template.Template{1:tmpl_cat_stk_err, 2:tmpl_cat_prc_err, 3:tmpl_cat_edt_err}
	states = map[string]string{"1475":"Andaman and Nicobar Islands","1476":"Andhra Pradesh","1477":"Arunachal Pradesh","1478":"Assam","1479":"Bihar","1480":"Chandigarh","3970":"Chhattisgarh","1481":"Dadra and Nagar Haveli","1482":"Daman and Diu","1483":"Delhi","1484":"Goa","1485":"Gujarat","1486":"Haryana","1487":"Himachal Pradesh","1488":"Jammu and Kashmir","3971":"Jharkhand","1489":"Karnataka","1490":"Kerala","1491":"Lakshadweep Islands","1492":"Madhya Pradesh","1493":"Maharashtra","1494":"Manipur","1495":"Meghalaya","1496":"Mizoram","1497":"Nagaland","1498":"Orissa","1499":"Pondicherry","1500":"Punjab","1501":"Rajasthan","1502":"Sikkim","1503":"Tamil Nadu","1504":"Tripura","1505":"Uttar Pradesh","3972":"Uttarakhand","1506":"West Bengal"}
	accessKey = "6ZUW285028D82TW1ELKT";
	secretKey = []byte("d20fc029ba8300c06443cb94a83b36db88fa7c67");
        merchantId = []byte("592cfpvxpr")
        ccy = []byte("INR")
)
// utils
        func renderTemplate(w http.ResponseWriter, tmpll string, s *Sheep) {
                err := tmpl_prm.ExecuteTemplate(w, tmpll+".html", s)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                }
        }

        func renderTemplate1(w http.ResponseWriter, tmpll string, s []*Sheep) {
                err := tmpl_prm.ExecuteTemplate(w, tmpll+".html", s)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                }
        }

// handlers for campaign registration
        const confirmMessage = `Thanks for registering. Please confirm your email address by clicking on the link: %s`
        func confirm(c *(appengine.Context),  email string) {
                url := fmt.Sprintf("%s%s","http://wattwerks.appspot.com/confirm/", email)
                msg := &mail.Message {
                        Sender: "Wattwerks Support <support@wattwerks.appspotmail.com>",
                        To: []string{email},
                        Subject: "Confirm your registration",
                        Body: fmt.Sprintf(confirmMessage, url),
                }
                if err := mail.Send(*c,msg); err != nil {
                        (*c).Errorf("Couldn't send email: %v", err)
                }
        }
        func handleConfirmPage(w http.ResponseWriter, r *http.Request ) {
                if r.Method != "GET" {
                        http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
                        return
                }
                m := validPath.FindStringSubmatch(r.URL.Path)
                if m == nil {
                        http.NotFound(w,r)
                        //renderTemplate(w, "err", nil)
                        return
                }
                c := appengine.NewContext(r)
                q := datastore.NewQuery("sheep").Filter("Email=", m[2]).Limit(1)
                var s0 Sheep
                //for s := q.Run(c) ; ; {
                if s := q.Run(c) ; s != nil {
                        key, err := s.Next(&s0)
                        //if err == datastore.Done {
                                //break
                        //}
                        if err != nil && err != datastore.Done {
                                renderTemplate(w, "err", nil)
                                return
                        }
                        if key != nil {
                          s0.Confirmed = true 
                          if _, err := datastore.Put(c, key, &s0); err != nil {
                                  renderTemplate(w, "err", nil)
                                  return
                          }
                          renderTemplate(w, "confirm", &s0)
                          return
                        }
                }
                renderTemplate(w, "err", &s0)
                return
        }

        func handleFlockPage(w http.ResponseWriter, r *http.Request) {
                if r.Method != "GET" {
                        http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
                        return
                }
                c := appengine.NewContext(r)
                q0 := datastore.NewQuery("sheep")
                var s0 []*Sheep

                key, err := q0.GetAll(c, &s0)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                if key != nil {
                        renderTemplate1(w, "flock", s0)
                        return
                }
                renderTemplate1(w, "flock", nil)
        }

        func handleMainPage(w http.ResponseWriter, r *http.Request) {
                if r.Method != "GET" {
                        http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
                        return
                }
                if r.URL.Path != "/promo/" {
                        http.NotFound(w, r)
                        return
                }
                c := appengine.NewContext(r)
                q0 := datastore.NewQuery("Sheep").Filter("Name=", "Guest")
                var s0 []*Sheep

                key, err := q0.GetAll(c, &s0)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                if key != nil {
                        renderTemplate(w, "reg", s0[0])
                        return
                }
                renderTemplate(w, "reg", &Sheep{Name: "Guest",Email: "guest@user.com",Addr: "123 JP Nagar", Phon: "123456", Area: "123", Note: "Rock on", JDte: time.Now()})
        }

        func handleRegPage(w http.ResponseWriter, r *http.Request ) {
                if r.Method != "GET" {
                        http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
                        return
                }
                c := appengine.NewContext(r)
                //fmt.Sprintf("%s", "Creating sheep w name: " + r.FormValue("name"))
                //validating email
                m := validEmail.FindStringSubmatch(r.FormValue("email"))
                if m == nil {
                        //http.Error(w, "Email invalid", http.StatusNotAcceptable)
                        renderTemplate(w, "err", nil)
                        return
                }
                confirm(&c, r.FormValue("email"))
                u1 := Sheep { Name: r.FormValue("name"), Email: r.FormValue("email"), Addr: r.FormValue("addr"),
                                Area: r.FormValue("area"), Note: r.FormValue("note"), JDte: time.Now() }
                if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"sheep", nil), &u1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                renderTemplate(w, "ta", &u1)
        }

        func handleHowPage(w http.ResponseWriter, r *http.Request ) {
                if r.Method != "GET" {
                        http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
                        return
                }
                renderTemplate(w, "how", nil)
        }

// handlers for goods admin
        func handleGoodsSearch(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                // get posted form data
                err := r.ParseForm()
                handle(err)
                v := r.Form
                srchtrm := v.Get("srchTrm")
                handle(err)
                q0 := datastore.NewQuery("product").Filter("Brand =", srchtrm)
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cc, err := getCstmr1(c, w, r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                key, err := q0.GetAll(c, &s0)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if key != nil {
                        data := Render {"", cc, s0, c0}
                        //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                        err = tmpl_gds_lst.ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                }
                data := Render {"No results", cc, s0, c0}
                err = tmpl_gds_lst.ExecuteTemplate(w,"base", data)
                handle(err)
        }

        func handleGoodsList(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                q0 := datastore.NewQuery("product")
                var s0 []Good
                var c0 []Category
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cc, err := getCstmr1(c, w, r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                key, err := q0.GetAll(c, &s0)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if key != nil {
                        data := Render {"", cc, s0, c0 }
                        //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                        err = tmpl_adm_gds_lst.ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                }
                cc, err = getCstmr(c, "Guest")
                hndl(err, "hndlGoodsList0")
                data := Render {"", cc, s0, c0 }
                err = tmpl_adm_gds_lst.ExecuteTemplate(w,"base", data)
                hndl(err, "hndlGoodsList1")
        }

        func handleGoodsListJSON(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                q0 := datastore.NewQuery("product")
                var s0 []Good
                key, err := q0.GetAll(c, &s0)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if key != nil {
                        gds, err := json.Marshal(s0)
                        handle(err)
                        w.Header().Set("Content-Type", "text/plain")
                        io.WriteString(w, string(gds))
                        //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                        //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", data)
                        return
                }
                //err = tmpl_adm_gds_lst.ExecuteTemplate(w,"admin", "")
                //handle(err)
                w.Header().Set("Content-Type", "text/plain")
                io.WriteString(w, "No goods in data store")
        }

        func handleGoodEntry(w http.ResponseWriter, r *http.Request ) { // this function responds to GET request by displaying form for single good entry
                c := appengine.NewContext(r)
                var s0 []Good
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cc, err := getCstmr1(c, w, r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                data := Render {"", cc, s0, c0 }
                err = tmpl_adm_gds_ntr.ExecuteTemplate(w,"base", data)
                handle(err)
        }

        func handleGoodsEntry(w http.ResponseWriter, r *http.Request ) { // this function responds to GET request by displaying form for mulitple goods file upload
                c := appengine.NewContext(r)
                uploadURL, err := blobstore.UploadURL(c, "/admin/goods/upload", nil)
                handle(err)
                cc, err := getCstmr1(c, w, r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                data := Render2 {uploadURL, cc, c0}
                err = tmpl_adm_gds_pld.ExecuteTemplate(w,"base", data)
                handle(err)
        }

        func handleGoodCreate(w http.ResponseWriter, r *http.Request ) { // this function handles conventional GET request, creates struct with form values and stores in datastore
                //c := appengine.NewContext(r)
                //pd1 := GoodDeets{ DescDetails: "/goods/docs/sma_sb24", Price: 6000.00, Tax: 0.155, Stock: 12, Related:[]int64{}, Prices: []float64{}, Volumes: []int{}, /*PriceVolume: map[int]float64{1:9000.00, 10:7500.00}, Parameters: map[string]string{"Power":"240W","":""},*/ ParameterNames: []string{}, ParameterValues: []string{}, Features: []string{}, Items: []string{}, UrlImgs1:"/imgs/small/sma_240sb.png", UrlImgs2: "/imgs/small/sma_240sb.png", UrlImgs3: "/imgs/small/sma_240sb.png", UrlFile: "/imgs/small/sma_240sb.png"}
                /*p1 := Good{ Code:"SMA_SB240", Category:"Inverter", Subcategory: "Micro", Brand:"SMA", Desc: "SunnyBoy 240W", Price: 6000.00, Url: "/goods/sma_sb240", Urlimg:"/goods/images/med/sma_sb240.png", Featured: true, Hidden:false, Deets: pd1 }
                if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"product", nil), &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"productdeets", nil), &pd1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                err := tmpl_adm_gds_ntr.ExecuteTemplate(w,"front", "")
                handle(err)*/
                c := appengine.NewContext(r)
                // get counter and find next available id. if problem getting next id then start over at 1
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 0
                        gi.NxtGd = 1001
                        gi.TtlCst = 0
                        gi.NxtCst = 10001
                        gi.TtlRdr = 0
                        gi.NxtRdr = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                http.Error(w, err.Error(), 500)
                                return
                        }
                }
                // get posted form data
                err := r.ParseForm()
                handle(err)
                v := r.Form
                price, err := strconv.ParseFloat(v.Get("price"), 64)
                handle(err)
                //tax, err := strconv.ParseFloat(v.Get("tax"), 64)
                //handle(err)
                stock, err := strconv.Atoi(v.Get("stock"))
                handle(err)
                featured, err := strconv.ParseBool(v.Get("featured"))
                handle(err)
                hidden, err := strconv.ParseBool(v.Get("hidden"))
                handle(err)
                //create GoodDeets and Good
                pd1 := GoodDeets{ DescDetails: v.Get("descDetails"), Price: price, Tax: 0.155, Stock: stock, Related:[]int64{}, Prices: []float64{}, Volumes: []int{}, ParameterNames: []string{}, ParameterValues: []string{}, Features: []string{}, Items: []string{}, UrlImgs1: v.Get("urlImgs") , UrlImgs2: v.Get("urlImgs"), UrlImgs3: v.Get("urlImgs"), UrlFile: v.Get("urlFile")}
                p1 := Good{ Id: gi.NxtGd, Code:v.Get("code"), Category: v.Get("category"), Subcategory: v.Get("subcategory"), Brand: v.Get("brand"), Desc: v.Get("desc"), Price: price, Url: v.Get("url"), Urlimg:v.Get("urlImg"), Featured: featured, Hidden: hidden, Deets: pd1 }
                pk := datastore.NewKey(c, "product", "", gi.NxtGd, nil)
                handle(err)
                if _, err := datastore.Put(c, pk, &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                if _, err := datastore.Put(c, datastore.NewKey(c,"productdeets","", gi.NxtGd+10000000, pk), &pd1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                gi.NxtGd++
                gi.TtlGd++
                if _, err := datastore.Put(c, ik, gi); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                handleGoodsList(w,r)
        }

        func handleGoodsCreate(w http.ResponseWriter, r *http.Request ) { // this function handles conventional GET request, creates struct with form values and stores in datastore
                c := appengine.NewContext(r)
                blobs, _, err := blobstore.ParseUpload(r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                file := blobs["file"]
                if len(file) == 0 {
                        log.Fatal("No file uploaded")
                        http.Redirect(w,r,"/admin/goods", http.StatusFound)
                        return
                }
                rdr := blobstore.NewReader(c, file[0].BlobKey)
                ubs := make([]byte, 100000)
                n, err := rdr.Read(ubs)
                ubs = ubs[:n]
                var goods []Good
                err = json.Unmarshal(ubs, &goods)
                handle(err)
                // get counter and find next available id. if problem getting next id then start over at 1
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 0
                        gi.NxtGd = 1001
                        gi.TtlCst = 0
                        gi.NxtCst = 10001
                        gi.TtlRdr = 0
                        gi.NxtRdr = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                http.Error(w, err.Error(), 500)
                                return
                        }
                }
                for _,good := range goods {
                        pk := datastore.NewKey(c, "product", "", gi.NxtGd, nil)
                        good.Id = gi.NxtGd
                        handle(err)
                        if _, err := datastore.Put(c, pk, &good); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                                handle(err)
                                return
                        }
                        if _, err := datastore.Put(c, datastore.NewKey(c,"productdeets","", gi.NxtGd+10000000, pk), &good.Deets); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                                handle(err)
                                return
                        }
                        gi.NxtGd++
                        gi.TtlGd++
                }
                if _, err := datastore.Put(c, ik, gi); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                //http.Redirect(w, r, "
                handleGoodsList(w,r)
        }

        func handleGoodCreateJS(w http.ResponseWriter, r *http.Request ) { // this function handles form  post data, creates struct and stores in datastore
                c := appengine.NewContext(r)
                // get counter and find next available id. if problem getting next id then start over at 1
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 0
                        gi.NxtGd = 1001
                        gi.TtlCst = 0
                        gi.NxtCst = 10001
                        gi.TtlRdr = 0
                        gi.NxtCst = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                http.Error(w, err.Error(), 500)
                                return
                        }
                }
                // get posted form data
                err := r.ParseForm()
                handle(err)
                v := r.Form
                price, err := strconv.ParseFloat(v.Get("price"), 64)
                handle(err)
                //tax, err := strconv.ParseFloat(v.Get("tax"), 64)
                //handle(err)
                stock, err := strconv.Atoi(v.Get("stock"))
                handle(err)
                featured, err := strconv.ParseBool(v.Get("featured"))
                handle(err)
                hidden, err := strconv.ParseBool(v.Get("hidden"))
                handle(err)
                //create GoodDeets and Good
                pd1 := GoodDeets{ DescDetails: v.Get("descDetails"), Price: price, Tax: 0.155, Stock: stock, Related:[]int64{}, Prices: []float64{}, Volumes: []int{}, ParameterNames: []string{}, ParameterValues: []string{}, Features: []string{}, Items: []string{}, UrlImgs1: v.Get("urlImgs") , UrlImgs2: v.Get("urlImgs"), UrlImgs3: v.Get("urlImgs"), UrlFile: v.Get("urlFile")}
                p1 := Good{ Id: gi.NxtGd, Code:v.Get("code"), Category: v.Get("category"), Subcategory: v.Get("subcategory"), Brand: v.Get("brand"), Desc: v.Get("desc"), Price: price, Url: v.Get("url"), Urlimg:v.Get("urlImg"), Featured: featured, Hidden: hidden, Deets: pd1 }
                pk := datastore.NewKey(c, "product", "", gi.NxtGd, nil)
                handle(err)
                if _, err := datastore.Put(c, pk, &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                if _, err := datastore.Put(c, datastore.NewKey(c,"productdeets","", gi.NxtGd+10000000, pk), &pd1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                gi.NxtGd++
                gi.TtlGd++
                if _, err := datastore.Put(c, ik, gi); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                err = tmpl_adm_gds_ntr.ExecuteTemplate(w,"goods_entry_success", "")
                handle(err)
        }

        func handleGoodCreateJS1(w http.ResponseWriter, r *http.Request ) { // this function handles json post data, marshals into struct and stores in datastore
                c := appengine.NewContext(r)
                // get counter and find next available id. if problem getting next id then start over at 1
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 0
                        gi.NxtGd = 1001
                        gi.TtlCst = 0
                        gi.NxtCst = 10001
                        gi.TtlRdr = 0
                        gi.NxtCst = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                http.Error(w, err.Error(), 500)
                                return
                        }
                }
                // decode json input into product
                var pd1 GoodDeets
                p1 := &Good{Deets:pd1}
                err := json.NewDecoder(r.Body).Decode(p1)
                handle(err)
                // store product in datastore
                pk := datastore.NewKey(c, "product", "", gi.NxtGd, nil)
                handle(err)
                if _, err := datastore.Put(c, pk, p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        //err := templates.ExecuteTemplate(w,"problem", "")
                        //handle(err)
                        return
                }
                if _, err := datastore.Put(c, datastore.NewKey(c,"productdeets","", gi.NxtGd+10000000, pk), &pd1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                gi.NxtGd++
                gi.TtlGd++
                if _, err := datastore.Put(c, ik, gi); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                err = tmpl_adm_gds_ntr_new.ExecuteTemplate(w,"goods_entry_success", "")
                handle(err)
        }

        func handleGoodEdit(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                id,err := strconv.ParseInt(vars["id"], 10, 64)
                handle(err)
                pk := datastore.NewKey(c, "product", "", id, nil)
                cp := new(Good)
                if err := datastore.Get(c, pk, cp); err != nil {
                                http.Error(w, err.Error(), 500)
                                return
                }
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                data := struct {
                        This *Good `json:"good,string"`
                        Categories []Category `json:"categories,string"`
                }{
                        cp, c0,
                }
                err = tmpl_adm_gds_edt.ExecuteTemplate(w,"base", data)
                handle(err)
        }

        /*
        func handleGoodError(c appengine.Context, w http.ResponseWriter, r *http.Request, errtyp int, errmsg []byte) {
                var s0 []Good
                var c0 []Category
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cs, err := getCstmr(c, "Guest")
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                session, err := ckstore.Get(r, "admin_path")
                handle(err)
                session.Values["lggd"] = "Guest"
                session.Save(r, w)
                data := Render{string(errmsg), cs, s0,c0}
                err = gd_errs[errtyp].ExecuteTemplate(w, "base", data)
                handle(err)
                return
        }

        func handleGoodEditPostJS(w http.ResponseWriter, r *http.Request ) { // this function handles json post data, marshals into struct and stores in datastore
                c := appengine.NewContext(r)
                // get posted form data
                err := r.ParseForm()
                handle(err)
                v := r.Form
                id, err := strconv.ParseInt(v.Get("Id"), 10, 64)
                handle(err)
                price, err := strconv.ParseFloat(v.Get("price"), 64)
                handle(err)
                //tax, err := strconv.ParseFloat(v.Get("tax"), 64)
                //handle(err)
                stock, err := strconv.Atoi(v.Get("stock"))
                handle(err)
                featured, err := strconv.ParseBool(v.Get("featured"))
                handle(err)
                hidden, err := strconv.ParseBool(v.Get("hidden"))
                handle(err)
                //create Gooddeets and Good
                pd1 := GoodDeets{ DescDetails: v.Get("descDetails"), Price: price, Tax: 0.155, Stock: stock, Related:[]int64{}, Prices: []float64{}, Volumes: []int{}, ParameterNames: []string{}, ParameterValues: []string{}, Features: []string{}, Items: []string{}, UrlImgs1: v.Get("urlImgs") , UrlImgs2: v.Get("urlImgs"), UrlImgs3: v.Get("urlImgs"), UrlFile: v.Get("urlFile")}
                p1 := Good{ Id: id, Code:v.Get("code"), Category: v.Get("category"), Subcategory: v.Get("subcategory"), Brand: v.Get("brand"), Desc: v.Get("desc"), Price: price, Url: v.Get("url"), Urlimg:v.Get("urlImg"), Featured: featured, Hidden: hidden, Deets: pd1 }
                pk := datastore.NewKey(c, "product", "", id, nil)
                handle(err)
                if _, err := datastore.Put(c, pk, &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                if _, err := datastore.Put(c, datastore.NewKey(c,"productdeets","", id+100000, pk), &pd1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        handle(err)
                        return
                }
                err = goods_edit.ExecuteTemplate(w,"goods_edit_success", "")
                handle(err)
        }*/

// handlers for account admin
        func handleCookieCheck(w http.ResponseWriter, r *http.Request) { //only ever called from redirect
                c := appengine.NewContext(r)
                var s0 []Good
                var c0 []Category
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cs, err := getCstmr(c, "Guest")
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                for i,ckp := range w.Header() {
                        log.Println(i, ckp)
                }
                _, err = r.Cookie("COOKIES")
                if err != nil {
                        data := Render4{string(cksnbl),cs,s0,c0, getCart(c,0,&cs)}
                        err = tmpl_acc_lgn.ExecuteTemplate(w,"base", data)
                        return
                }
                http.Redirect(w, r, "/account/login", 301) // cookies enabled - resume normal flow
                return
        }

        func handleAccount(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if session.Values["lggd"] != nil { // returning user
                        if session.Values["lggd"].(string) == "Guest"{  // guest
                                cs, err := getCstmr(c, "Guest")
                                if err != nil {
                                        srvErr(c, w, err)
                                        return
                                }
                                data := Render4{"",cs,s0,c0, getCart(c,session.Values["crtd"].(int64),&cs)}
                                session.Values["lggd"] = "Guest"
                                session.Save(r, w)
                                err = tmpl_acc_reg.ExecuteTemplate(w,"base", data)
                                handle(err)
                                return
                        }
                        var cs Cust //user
                        cs, err = getCstmr(c, session.Values["lggd"].(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        data := Render4{"",cs,s0,c0, getCart(c, session.Values["crtd"].(int64), &cs)}
                        err := tmpl_acc.ExecuteTemplate(w,"base", data)
                        handle(err)
                        return
                }
                ckchck(w, r) //if session is nil check whether cookies enabled
        }

        func handleAccountwTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template , rndr interface{}) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if session.Values["lggd"] != nil { // returning user
                        cs, err := getCstmr(c, session.Values["lggd"].(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        crtd := session.Values["crtd"].(int64)
                        if session.Values["lggd"] != "Guest"{  // guest
                                switch rndr := rndr.(type) {
                                case Cart:
                                        s0 = getGoods(c, rndr.Ids)
                                        data := Render4{"",cs,s0,c0,rndr}
                                        err = tmpl.ExecuteTemplate(w,"base", data)
                                        handle(err)
                                        return
                                case []Order:
                                        data := Render3{"",cs,s0,c0, getCart(c, crtd, &cs), rndr}
                                        err = tmpl.ExecuteTemplate(w,"base", data)
                                        handle(err)
                                        return
                                case Render2:
                                        data := Render4{"",cs,s0,c0, getCart(c, 0, &cs)}
                                        err = tmpl.ExecuteTemplate(w,"base", data)
                                        handle(err)
                                        return
                                case Render40:
                                        s0 = getGoods(c, rndr.Crt.Ids)
                                        rndr.Message = ""
                                        rndr.Cstmr = cs
                                        rndr.Goods = s0
                                        rndr.Categories = c0
                                        err = tmpl.ExecuteTemplate(w, "base", rndr)
                                        hndl(err, "handleAccountwTemplate5")
                                        return
                                default:
                                        data := Render4{"",cs,s0,c0, getCart(c, 0, &cs)}
                                        err = tmpl.ExecuteTemplate(w,"base", data)
                                        handle(err)
                                        return
                                }
                        }
                        data := Render4{"",cs,s0,c0,getCart(c,0,&cs)}
                        err = tmpl.ExecuteTemplate(w,"base", data)
                        handle(err)
                        return
                }
                ckchck(w, r) //if session is nil check whether cookies enabled
        }

        func handleAccountRegister(w http.ResponseWriter, r *http.Request ) {
                /*session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                for i,x := range r.Header {
                        log.Println(i, x )
                }
                ck := r.RemoteAddr + "  " + r.Header["User-Agent"][0]
                log.Println(ck)*/
                handleAccount(w,r)
        }

        func handleAccountLogin(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
		
                if session.Values["lggd"] != nil && session.Values["lggd"].(string) != "Guest" { // returning user
                        var cs Cust
                        cs, err = getCstmr(c, session.Values["lggd"].(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        data := Render4{"",cs,s0,c0, getCart(c, session.Values["crtd"].(int64), &cs)}
                        err := tmpl_acc.ExecuteTemplate(w,"base", data)
                        hndl(err,"handleAccountLogin0")
                        return
                }
                cs, err := getCstmr(c, "Guest")
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
		if session.Values["crtd"] == nil {
                	session.Values["crtd"] = int64(0)
		}
                data := Render4{"",cs,s0,c0, getCart(c, session.Values["crtd"].(int64), &cs)}
                session.Values["lggd"] = "Guest"
                session.Values["crtd"] = int64(0)
                session.Save(r, w)
                err = tmpl_acc_lgn.ExecuteTemplate(w,"base", data)
                hndl(err,"handleAccountLogin1")
        }

        func handleAccountLoginP(w http.ResponseWriter, r *http.Request ) { // this function handles POST request, creates struct with form values and stores in datastore
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                handle(err)
                gdsession, err := ckstore.Get(r, "goods_path")
                handle(err)
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountLoginP0")
                //validation
                  var ok bool
                  v := r.Form
                  em := v.Get("email")
                  //validating email
                  m := validEmail.FindStringSubmatch(em)
                  if m == nil {
                          handleAccountError(c, w, r, 2, []byte("Invalid email"))
                          return
                  }
                  pw := v.Get("pswd")
                  if pw == "" {
                          handleAccountError(c, w, r, 1, []byte("Password incorrect"))
                          return
                  }
                  xs,err := getCstmre(c,em)
                  hndl(err,"handleAccountLoginP1")
                  if xs.Firstname != "Guest" {
                          if xs.Password == pw {
                                  ok = true
                          } else {
                                  ok = false
                                  log.Println(xs.Email + " Incorrect password")
                                  handleAccountError(c,w,r,1,[]byte("Incorrect password"))
                                  return
                          }
                  }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if ok {
                        crtd := mkCart(c, xs.Email)
                        session.Values["lggd"] = xs.Firstname
                        session.Values["crtd"] = crtd
                        gdsession.Values["lggd"] = xs.Firstname
                        gdsession.Values["crtd"] = crtd
			session.Options = &sessions.Options{
						Path:"/account/",
					}
			gdsession.Options = &sessions.Options{
						Path:"/goods/",
					}
                        sessions.Save(r, w)
                        data := Render4 {"Logged in", xs, s0, c0, getCart(c, crtd, &xs)}
                        err = tmpl_acc.ExecuteTemplate(w, "base", data)
                        hndl(err, "handleAccountLoginP1")
                        return
                }
                // check if cstmr already has session
                cc, err := getCstmr1(c, w, r)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if cc.Firstname != "Guest" {
                        data := Render {"Already logged in", cc, s0, c0}
                        err = tmpl_acc.ExecuteTemplate(w, "base", data)
                        hndl(err, "handleAccountLoginP2")
                        return
                }
                session.Values["lggd"] = "Guest"
                session.Save(r, w)
                handleAccountLogin(w,r)
        }

        func handleAccountEdit(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session,err := ckstore.Get(r, "account_path")
                hndl(err, "handleAccountEdit0")
                lggd := session.Values["lggd"].(string)
                //crtd := session.Values["crtd"].(int64)
                switch lggd {
                        case "Guest", "":
                                handleAccount(w,r)
                                return
                        default:
                                cstmr,err := getCstmr(c, lggd)
                                hndl(err, "handleAccountEdit1")
                                if cstmr.Firstname != "Guest" || cstmr.Firstname != "" {
                                        handleAccountwTemplate(w, r, tmpl_acc_edt, cstmr)
                                        return
                                }
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                }
        }

        func handleAccountEditP(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                hndl(err, "handleAccountEditP0")
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountEditP1")
                //validation
                  v := r.Form
                  fn := v.Get("firstname")
                  ln := v.Get("lastname")
                  em := v.Get("email")
                  ph := v.Get("phone")
                  ad := v.Get("address_1")
                  ci := v.Get("city")
                  st := v.Get("state")
                  pc := v.Get("postcode")
                  ct := v.Get("indi")
                  pm := v.Get("pmt_method")
                  cm := v.Get("company")
                  ag := v.Get("agree")
                  pw := v.Get("password")
                  cn := v.Get("confirm")
                  //meterid, err := strconv.Atoi(v.Get("meterid"),64)
                  if fn=="" || ln=="" || em=="" || ph=="" || ad=="" || ci=="" || st=="" || pc=="" || ct=="" || pm=="" || ag=="" {
                          handleAccountError(c, w, r, 3, []byte("Missing required data"))
                          return
                  }
                  //validating email
                  m := validEmail.FindStringSubmatch(r.FormValue("email"))
                  if m == nil {
                          handleAccountError(c, w, r, 3, []byte("Invalid email"))
                          return
                  }
                  xs, err := getCstmre(c, em)
                  if xs.Firstname == "Guest" {
                          handleAccountError(c, w, r, 3, []byte("Invalid first name"))
                          return
                  }
                  if pw == "" || pw != cn || pw != xs.Password {
                          handleAccountError(c, w, r, 3, []byte("Password mismatch"))
                          return
                  }
                  var meterid int64
                  if v.Get("meterid") != "" {
                          meterid, err = strconv.ParseInt(v.Get("meterid"),10,64)
                          hndl(err, "handleAccountEditP2")
                  } else {
                          meterid = 0
                  }
                  indi, err := strconv.ParseBool(v.Get("indi"))
                  hndl(err, "handleAccountEditP3")
                  if !indi && cm=="" {
                          handleAccountError(c, w, r, 3, []byte("Missing required company data"))
                          return
                  }
                //amend Cust
                p1 := Cust{ Id: xs.Id, Email:em, Firstname: fn, Lastname: ln, Company: cm, Phone: ph, Addr: ad, City: ci, State: states[st], Postcode: pc, Indi: indi, Meterid: meterid, Password:pw, Pmtmethod: pm, Crtid: xs.Crtid}
                pk := datastore.NewKey(c, "cstmr", "", xs.Id, nil)
                hndl(err, "handleAccountCreateP4")
                if _, err := datastore.Put(c, pk, &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        hndl(err, "handleAccountCreateP5")
                        return
                }
                session.Values["lggd"] = p1.Firstname
                session.Values["crtd"] = p1.Crtid
                session.Save(r, w)
                handleAccountEdit(w,r)
        }

        func handleAccountCart(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                lggd := session.Values["lggd"]
                if lggd != "Guest" {
                        cs, err := getCstmr(c, lggd.(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        crtd := session.Values["crtd"]
                        var crt Cart
                        if crtd != nil {
                                crt = getCart(c, crtd.(int64), &cs) //must make sure cart id belongs to user
                                if err != nil {
                                        srvErr(c, w, err)
                                        return
                                }
                        } else {
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                        }
                        handleAccountwTemplate(w, r, tmpl_acc_crt, crt)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountCartAdd(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                gd,err := strconv.ParseInt(vars["id"], 10, 64)
                hndl(err,"handleAccountCartAdd0")
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                lggd := session.Values["lggd"]
                if lggd != "Guest" {
                        //cs, err := getCstmr(c, lggd.(string))
                        //if err != nil {
                        //	srvErr(c, w, err)
                        //	return
                        //}
                        crtd := session.Values["crtd"]
                        var crt Cart
                        if crtd != nil {
                                //crt = getCart(c, crtd.(int64), &cs) //must make sure cart id belongs to user
                                //if err != nil {
                                //	srvErr(c, w, err)
                                //	return
                                //}
                                crt = addCart(c, crtd.(int64), gd)
                                handle(err)
                        } else {
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                        }
                        handleAccountwTemplate(w, r, tmpl_acc_crt, crt)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountCartPdt(w http.ResponseWriter, r *http.Request) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r,"account_path")
                hndl(err,"handleAccountCartPdt0")
                //validation
                //ok := false
                err = r.ParseForm()
                hndl(err, "handleAccountCartPdt1")
                v := r.PostForm
                qs := make([]int,0)
                for _,q := range v {
                        qq,err := strconv.Atoi(q[0])
                        hndl(err,"handleAccountPtd2")
                        qs = append(qs, qq)
                }
                lggd := session.Values["lggd"].(string)
                if lggd != "Guest" {
                        crtd := session.Values["crtd"]
                        var crt Cart
                        if crtd != nil {
                                crt = pdtCart(c, crtd.(int64), qs)
                        } else {
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                        }
                        handleAccountwTemplate(w, r, tmpl_acc_crt, crt)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountCartRm(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                gd,err := strconv.ParseInt(vars["id"], 10, 64)
                hndl(err,"handleAccountCartRm0")
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                lggd := session.Values["lggd"]
                if lggd != "Guest" {
                        //cs, err := getCstmr(c, lggd.(string))
                        //if err != nil {
                        //	srvErr(c, w, err)
                        //	return
                        //}
                        crtd := session.Values["crtd"]
                        var crt Cart
                        if crtd != nil {
                                //crt = getCart(c, crtd.(int64), &cs) //must make sure cart id belongs to user
                                //if err != nil {
                                //	srvErr(c, w, err)
                                //	return
                                //}
                                crt = rmCart(c, crtd.(int64), gd)
                                hndl(err,"handleAccountCartRm1")
                        } else {
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                        }
                        handleAccountwTemplate(w, r, tmpl_acc_crt, crt)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountChckt(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                lggd := session.Values["lggd"]
                if lggd != "Guest" {
                        cs, err := getCstmr(c, lggd.(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        crtd := session.Values["crtd"]
                        var crt Cart
                        if crtd != nil {
                                crt = getCart(c, crtd.(int64), &cs) // must make sure cart id belongs to user
                                if err != nil {
                                        srvErr(c, w, err)
                                        return
                                }
                        } else {
                                handleAccountError(c, w, r, 1, []byte("Sorry there's been a problem. Please login."))
                                return
                        }
                        id := nxtOrderId(c)
                        hsh := append( strconv.AppendInt( strconv.AppendFloat( append( make([]byte,0), merchantId...), crt.Ttl, 'f', 0, 64), id, 10), ccy...)
                        secSign := string(hMac(hsh, secretKey))
                        rdr := Order{ Id: id, FirstName: cs.Firstname, LastName: cs.Lastname, Email: cs.Email, ReqTime: time.Now(), SecSign: secSign}
                        data := Render40{Crt: crt, Rdr: rdr}
                        _,err = addOrder(c, rdr)
                        hndl(err, "handleAccountChckt1")
                        handleAccountwTemplate(w, r, tmpl_acc_chckt, data)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountLogout(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                gdsession, err := ckstore.Get(r, "goods_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cs, err := getCstmr(c, "Guest")
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                data := Render4{"",cs,s0,c0, getCart(c, 0, &cs)}
                session.Values["lggd"] = "Guest"
                session.Values["crtd"] = int64(0)
                session.Options = &sessions.Options{
                                        Path:"/account/",
                                }
                gdsession.Values["lggd"] = "Guest"
                gdsession.Values["crtd"] = int64(0)
                gdsession.Options = &sessions.Options{
                                        Path:"/goods/",
                                }
                err = sessions.Save(r, w)
		handle(err)
                err = tmpl_acc_lgn.ExecuteTemplate(w,"base", data)
                handle(err)
        }

        func handleAccountCreateP(w http.ResponseWriter, r *http.Request ) { // this function handles POST request, creates struct with form values and stores in datastore
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                hndl(err, "handleAccountCreateP0")
                // get counter and find next available id. if problem getting next id then start over at 1
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 0
                        gi.NxtGd = 1001
                        gi.TtlCst = 0
                        gi.NxtCst = 10001
                        gi.TtlRdr = 0
                        gi.NxtRdr = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                http.Error(w, err.Error(), 500)
                                return
                        }
                }
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountCreateP1")
                //validation
                  v := r.Form
                  fn := v.Get("firstname")
                  ln := v.Get("lastname")
                  em := v.Get("email")
                  ph := v.Get("phone")
                  ad := v.Get("address_1")
                  ci := v.Get("city")
                  st := v.Get("state")
                  pc := v.Get("postcode")
                  ct := v.Get("indi")
                  pm := v.Get("pmt_method")
                  cm := v.Get("company")
                  ag := v.Get("agree")
                  pw := v.Get("password")
                  cn := v.Get("confirm")
                  //meterid, err := strconv.Atoi(v.Get("meterid"),64)
                  if fn=="" || ln=="" || em=="" || ph=="" || ad=="" || ci=="" || st=="" || pc=="" || ct=="" || pm=="" || ag=="" {
                          handleAccountError(c, w, r, 2, []byte("Missing required data"))
                          return
                  }
                  //validating email
                  m := validEmail.FindStringSubmatch(r.FormValue("email"))
                  if m == nil {
                          handleAccountError(c, w, r, 2, []byte("Invalid email"))
                          return
                  }
                  xs, err := getCstmre(c, em)
                  if xs.Firstname != "Guest" {
                          handleAccountError(c, w, r, 2, []byte("Email already exists in system"))
                          return
                  }
                  if pw != "" && pw != cn {
                          handleAccountError(c, w, r, 2, []byte("Password mismatch"))
                          return
                  }
                  var meterid int64
                  if v.Get("meterid") != "" {
                          meterid, err = strconv.ParseInt(v.Get("meterid"),10,64)
                          hndl(err, "handleAccountCreateP2")
                  } else {
                          meterid = 0
                  }
                  indi, err := strconv.ParseBool(v.Get("indi"))
                  hndl(err, "handleAccountCreateP3")
                  if !indi && cm=="" {
                          handleAccountError(c, w, r, 2, []byte("Missing required company data"))
                          return
                  }
                //create Cust
                p1 := Cust{ Id: gi.NxtCst, Email:em, Firstname: fn, Lastname: ln, Company: cm, Phone: ph, Addr: ad, City: ci, State: states[st], Postcode: pc, Indi: indi, Meterid: meterid, Password:pw, Pmtmethod: pm, Crtid: gi.NxtCst+1e5}
                pk := datastore.NewKey(c, "cstmr", "", gi.NxtCst, nil)
                hndl(err, "handleAccountCreateP4")
                if _, err := datastore.Put(c, pk, &p1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        hndl(err, "handleAccountCreateP5")
                        return
                }
                c1 := Cart{ Id: gi.NxtCst+1e5, Ids: []int64{}, Qntts: []int{}}
                pk = datastore.NewKey(c, "cart", "", gi.NxtCst+1e5, nil)
                if _, err := datastore.Put(c, pk, &c1); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        hndl(err, "handleAccountCreateP6")
                        return
                }
                gi.NxtCst++
                gi.TtlCst++
                if _, err := datastore.Put(c, ik, gi); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        err := tmpl_cmn.ExecuteTemplate(w,"problem", "")
                        hndl(err, "handleAccountCreateP7")
                        return
                }
                session.Values["lggd"] = p1.Firstname
                session.Values["crtd"] = p1.Crtid
                session.Options = &sessions.Options{
                                        Path:"/account/",
                                }
                session.Save(r, w)
                http.Redirect(w, r, "/account/", 307)
                //handleAccount(w,r)
        }

        func handleAccountList(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var cc Cust
                if session.Values["lggd"] != nil && session.Values["lggd"].(string) != "Guest" { // returning user
                        cc, err = getCstmr(c, session.Values["lggd"].(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                } else {
                        cc, err = getCstmr(c, "Guest")
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                }
                q0 := datastore.NewQuery("cstmr")
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                var cs []Cust
                key, err := q0.GetAll(c, &cs)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                if key != nil {
                        data := Render1 {"",  cc, s0, c0, cs}
                        //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                        err = tmpl_adm_acc_lst.ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                }
                err = tmpl_adm_acc_lst.ExecuteTemplate(w,"base", "")
                handle(err)
        }

        func handleAccountForgot(w http.ResponseWriter, r *http.Request ) {
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
		if session.Values["lggd"] != nil && session.Values["lggd"] != "Guest" {
			handleAccountEdit(w,r)
			return
		}
                handleAccountwTemplate(w, r, tmpl_acc_frgt, 0)
                return
        }

        func handleAccountForgotP(w http.ResponseWriter, r *http.Request ) { // this function handles post submission of email and firstname for pswd retrieval
                c := appengine.NewContext(r)
                // get posted form data
                err := r.ParseForm()
                hndl(err, "handleAccountForgotP0")
                //validation
                  var ok bool
                  v := r.Form
                  em := (v.Get("email"))
                  //validating email
                  m := validEmail.FindStringSubmatch(em)
                  if m == nil {
                          handleAccountError(c, w, r, 2, []byte("Invalid email"))
                          return
                  }
                  fs := strings.TrimSpace(v.Get("firstname"))
                  if fs == "" {
                          handleAccountError(c, w, r, 1, []byte("Invalid first name"))
                          return
                  }
                  xs,err := getCstmre(c,em)
                  hndl(err,"handleAccountForgotP1")
                  if xs.Firstname != "Guest" && strings.ToLower(xs.Firstname) == strings.ToLower(fs) {
                          ok = true
		  } else {
                          ok = false
                          handleAccountError(c,w,r,1,[]byte("We don't have that email/firstname registered"))
                          return
		  }
		if ok == true { // send email with password
			reminderEmail(&c, em, xs.Password)	
			logout(w, r)
			handleAccountwTemplate(w,r, tmpl_acc_rmnd, 0)
			return
		}
                handleAccountwTemplate(w, r, tmpl_acc_frgt, 0)
                return
        }
        
	const reminder = `You requested your password, which is: %s` 
	func reminderEmail(c *(appengine.Context),  e string, p string) {
                msg := &mail.Message {
                        Sender: "Wattwerks Support <support@wattwerks.appspotmail.com>",
                        To: []string{e},
                        Subject: "Your Wattwerks password",
                        Body: fmt.Sprintf(reminder, p),
                }
                if err := mail.Send(*c,msg); err != nil {
                        (*c).Errorf("Couldn't send email: %v", err)
                }
        }

        func handleAccountOrders(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                lggd := session.Values["lggd"]
                if lggd != "Guest" {
                        cs, err := getCstmr(c, lggd.(string))
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        rdrs, ks, err := getOrders(c, cs.Id)
                        hndl(err, "handleAccountOrders0")
                        if ks != nil {
                                handleAccountwTemplate(w, r, tmpl_acc_rdrs, rdrs)
                                return
                        }
                        handleAccountwTemplate(w, r, tmpl_acc_rdrs, rdrs)
                        return
                }
                handleAccount(w,r)
        }

        func handleAccountReturns(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                hndl(err, "handleAccountReturns0")
                lggd := session.Values["lggd"].(string)
                if lggd != "Guest" {
                        cs, err := getCstmr(c, lggd)
                        hndl(err, "handleAccountReturns1")
                        rdrs, err := getCmpltOrders(c, cs.Email)
                        hndl(err, "handleAccountReturns2")
                        handleAccountwTemplate(w, r, tmpl_acc_rtrns, rdrs)
                        return
                }
                handleAccount(w,r)
        }

        func handleOrdersList(w http.ResponseWriter, r *http.Request ) {
                c := appengine.NewContext(r)
                _, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var rdrs []Order
                q0 := datastore.NewQuery("order")
                _, err = q0.GetAll(c, &rdrs)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                handleAccountwTemplate(w, r, tmpl_adm_rdrs, rdrs)
                return
        }

        func handleAccountOrderP(w http.ResponseWriter, r *http.Request ) { // CITRUS EMULATOR
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                crtd := session.Values["crtd"].(int64)
                hndl(err, "handleAccountOrderP00")
                lggd := session.Values["lggd"]
                cs, err := getCstmr(c, lggd.(string))
                hndl(err, "handleAccountOrderP0")
                //crt := getCart(c, crtd, &cs)
                //hndl(err, "handleAccountOrderP1")
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountOrderP2")
                //validation
                v := r.Form
                pd := make([]string,0)
                pd = append(pd, v.Get("crtid")) // 0
                pd = append(pd, v.Get("merchantTxnId")) // 1
                pd = append(pd, v.Get("secSignature")) // 2
                pd = append(pd, v.Get("reqtime")) // 3
                pd = append(pd, v.Get("orderAmount")) // 4
                pd = append(pd, v.Get("firstname")) // 5
                pd = append(pd, v.Get("lastname"))  // 6
                pd = append(pd, v.Get("email")) // 7
                pd = append(pd, v.Get("currency")) // 8
                pd = append(pd, v.Get("addressStreet1")) // 9
                pd = append(pd, v.Get("addressCity")) // 10
                pd = append(pd, v.Get("addressZip")) // 11
                pd = append(pd, v.Get("addressState")) // 12
                pd = append(pd, v.Get("addressCountry")) // 13
                pd = append(pd, v.Get("phoneNumber")) // 14
                pd = append(pd, v.Get("paymentMode")) // 15
                pd = append(pd, v.Get("returnUrl")) // 16 
                pd = append(pd, v.Get("notifyUrl")) // 17
                pd = append(pd, v.Get("templateCode")) // 18
                pd = append(pd, v.Get("reqtime")) // 18
                for i,p := range pd {
                        if i < 9 && p == "" {
                                log.Println("Missing data: ", i)
                                handleAccountError(c, w, r, 4, []byte("Missing required data:"))
                                return
                        }
                }
                data := append(merchantId, []byte(pd[4] + pd[1] + pd[8])...)
                secSign,_ := hex.DecodeString(pd[2])
                flag := checkMac(data, secSign, secretKey)
                pdi, err := strconv.ParseInt(pd[1],10,64)
                hndl(err, "handleAccountOrderP4")
                pdc,err := strconv.ParseInt(pd[0],10,64)
                hndl(err, "handleAccountOrderP4")
                if pdc != crtd {
                        handleAccountError(c, w, r, 4, []byte("Supplied cart id doesn't match session cart id"))
                        return
                }
                if ! flag {
                        handleAccountError(c, w, r, 4, []byte("Hmac error"))
                        return
                }
                rdr := Order{
                        Id:pdi,  //1,000,000-9,999,999
                        TxId:pd[1],
                        TxRefNo:pd[2],
                        TxStatus:"Incompete",
                        TxMsg:"",
                        PgTxnNo:"",
                        IssuerRefNo:"",
                        AuthIdCode:"",
                        PgRespCode:"",
                        OriginalAmount:pd[4],
                        AdjustedAmount:pd[4],
                        DpRuleName:"",
                        DpRuleType:"",
                        Amount:pd[4],
                        TransactionAmount:pd[4],
                        PaymentMode:pd[15],
                        TxGateway:"",
                        IssuerCode:"",
                        TxnDateTime:"",
                        IsCOD:"false",
                        /*Details of un-registered buyer*/
                        FirstName:pd[5],
                        LastName:pd[6],
                        Email:pd[7],
                        AddressStreet1:pd[9],
                        AddressStreet2:pd[9],
                        AddressFirstName:"",
                        AddressZip:pd[11],
                        AddressState:pd[12],
                        AddressCountry:pd[13],
                        MobileNo:pd[14],
                        /*user account if buyer registered*/
                        //Buyer *Cust //`json:"buyer"`
                        //Cart Cart `json:"cart"`
                        Cstid: cs.Id,
                        Crtid: crtd,
                        //SecSign: string(secSign),
                        Notified:false,
                }
                rdr,_,err = pdtOrder(c, rdr)
                hndl(err, "handleAccountOrderP3")
                handleAccountOrders(w, r)
                return
        }

        func handleCitrusP(w http.ResponseWriter, r *http.Request) {
                c := appengine.NewContext(r)
                //session,err := ckstore.Get(r, "wattwerks_cookies")
                //hndl(err, "handleCitrusP0")
                err := r.ParseForm()
                hndl(err, "handleCitrusP1")
                v := r.Form
                data := ""
                txnId := v.Get("TxId")
                txnStatus := v.Get("TxStatus")
                txnDateTime := v.Get("txnDateTime")
                amount := v.Get("amount")
                txRefNo := v.Get("TxRefNo")
                pgTxnId := v.Get("pgTxnNo")
                issuerRefNo := v.Get("issuerRefNo")
                authIdCode := v.Get("authIdCode")
                firstName := v.Get("firstName")
                lastName := v.Get("lastName")
                email := v.Get("email")
                pgRespCode := v.Get("pgRespCode")
                zipCode := v.Get("addressZip")
                reqSignature := v.Get("signature")
                //signature := ""
                flag := false
                if txnId != "" {
                        data += txnId
                }
                if txnStatus != "" {
                        data += txnStatus
                }
                if amount != "" {
                        data += amount
                }
                if pgTxnId != "" {
                        data += pgTxnId
                }
                if issuerRefNo != "" {
                        data += issuerRefNo
                }
                if authIdCode != "" {
                        data += authIdCode
                }
                if firstName != "" {
                        data += firstName
                }
                if lastName != "" {
                        data += lastName
                }
                if pgRespCode != "" {
                        data += pgRespCode
                }
                if zipCode != "" {
                        data += zipCode
                }
                //com.citruspay.util.HMACSignature  sigGenerator = new com.citruspay.util.HMACSignature (); 
                //signature = sigGenerator.generateHMAC(data, key);
                //if(reqSignature !=null && !reqSignature.equalsIgnoreCase("") &&!signature.equalsIgnoreCase(reqSignature)){
                flag = checkMac([]byte(data), []byte(reqSignature), secretKey)
                if flag {
                        byr,err := getCstmre(c, email)
                        hndl(err, "handleCitrusPost01")
                        td,err := strconv.ParseInt(txnId,10,64)
                        hndl(err, "handleCitrusPost02")
                        x,k,err := getOrder(c, td)
                        hndl(err, "handleCitrusPost1")
                        if byr.Firstname != "Guest" {
                                if txnStatus != "" && txnStatus == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleCitrusPost2")
                                        } else {
                                                k,err = addOrder(c, Order{TxnDateTime: txnDateTime, TxRefNo: txRefNo, TxId: txnId, TxStatus: txnStatus, Amount: amount})
                                                hndl(err, "handleCitrusPost3")
                                        }
                                }
                        } else { // checkout as guest
                                if txnStatus != "" && txnStatus == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleCitrusPost4")
                                        } else {
                                                k,err = addOrder(c, Order{TxnDateTime: txnDateTime, TxRefNo: txRefNo, TxId: txnId, TxStatus: txnStatus, Amount: amount, FirstName: firstName, LastName: lastName, Email: email})
                                                hndl(err, "handleCitrusPost5")
                                        }
                                }
                        }
                } else {
                        handleAccountError(c, w, r, 4, []byte("HMAC Error"))
                }
                return
        }

        func handleAccountNewOrderP(w http.ResponseWriter, r *http.Request ) { // CITRUS return
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                crtd := session.Values["crtd"].(int64)
                hndl(err, "handleAccountNewOrderP0")
                lggd := session.Values["lggd"]
                cs, err := getCstmr(c, lggd.(string))
                hndl(err, "handleAccountOrderP1")
                //crt := getCart(c, crtd, &cs)
                //hndl(err, "handleAccountOrderP1")
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountOrderP2")
                //validation
                v := r.Form
                pd := make([]string,0)
                pd = append(pd, v.Get("TxId")) // 0
                pd = append(pd, v.Get("TxStatus")) // 1
                pd = append(pd, v.Get("amount")) // 2
                pd = append(pd, v.Get("pgTxnNo")) // 3
                pd = append(pd, v.Get("issuerRefNo"))  // 4
                pd = append(pd, v.Get("authIdCode")) // 5
                pd = append(pd, v.Get("firstName")) // 6
                pd = append(pd, v.Get("lastName")) // 7
                pd = append(pd, v.Get("pgRespCode")) // 8
                pd = append(pd, v.Get("addressZip")) // 9
                pd = append(pd, v.Get("txnDateTime")) // 10
                pd = append(pd, v.Get("TxRefNo")) // 11
                pd = append(pd, v.Get("email")) // 12
                pd = append(pd, v.Get("signature")) // 13
                flag := false
                for _,p := range pd {
                        if p == "" {
                                handleAccountError(c, w, r, 4, []byte("Missing required data"))
                                return
                        }
                }
                pdi, err := strconv.ParseInt(pd[0],10,64)
                hndl(err, "handleAccountNewOrderP3")
                data := pd[0]+pd[1]+pd[2]+pd[3]+pd[4]+pd[5]+pd[6]+pd[7]+pd[8]+pd[9]
                secSign,_ := hex.DecodeString(pd[13])
                flag = checkMac([]byte(data), secSign, secretKey)
                //pdi,err := strconv.ParseInt(pd[0],10,64)
                //if pdi != crtd {
                //        handleAccountError(c, w, r, 4, []byte("Supplied cart id doesn't match session cart id"))
                //        return
                //}
                var rdr Order
                if flag {
                        byr,err := getCstmre(c, pd[12])
                        hndl(err, "handleCitrusPost01")
                        x,k,err := getOrder(c, pdi)
                        hndl(err, "handleCitrusPost1")
                        if k != nil {
                                x.Id = pdi  //1,000,000-9,999,999
                                x.TxId = pd[0]
                                x.TxRefNo = pd[11]
                                x.TxStatus = pd[1]
                                x.TxMsg = ""
                                x.PgTxnNo = pd[3]
                                x.IssuerRefNo = pd[4]
                                x.AuthIdCode = pd[5]
                                x.PgRespCode = pd[8]
                                x.OriginalAmount = pd[2]
                                x.AdjustedAmount = pd[2]
                                x.DpRuleName = ""
                                x.DpRuleType = ""
                                x.Amount = pd[2]
                                x.TransactionAmount = pd[2]
                                x.PaymentMode = pd[15]
                                x.TxGateway = ""
                                x.IssuerCode = ""
                                x.TxnDateTime = pd[10]
                                x.IsCOD = "false"
                                /*Details of un-registered buyer*/
                                x.FirstName = pd[6]
                                x.LastName = pd[7]
                                x.Email = pd[12]
                                x.AddressStreet1 = ""
                                x.AddressStreet2 = ""
                                x.AddressFirstName = ""
                                x.AddressZip = pd[9]
                                x.AddressState = ""
                                x.AddressCountry = ""
                                x.MobileNo = ""
                                /*user account if buyer registered*/
                                //Buyer *Cust //`json:"buyer"`
                                //Cart Cart `json:"cart"`
                                x.Cstid = cs.Id
                                x.Crtid = crtd
                                x.SecSign = string(secSign)
                                x.Notified = false
                        } else {
                                rdr = Order{
                                        Id:pdi,  //1,000,000-9,999,999
                                        TxId:pd[0],
                                        TxRefNo:pd[11],
                                        TxStatus:pd[1],
                                        TxMsg:"",
                                        PgTxnNo: pd[3],
                                        IssuerRefNo: pd[4],
                                        AuthIdCode: pd[5],
                                        PgRespCode: pd[8],
                                        OriginalAmount: pd[2],
                                        AdjustedAmount: pd[2],
                                        DpRuleName:"",
                                        DpRuleType:"",
                                        Amount:pd[2],
                                        TransactionAmount:pd[2],
                                        PaymentMode:pd[15],
                                        TxGateway:"",
                                        IssuerCode:"",
                                        TxnDateTime: pd[10],
                                        IsCOD:"false",
                                        /*Details of un-registered buyer*/
                                        FirstName:pd[6],
                                        LastName:pd[7],
                                        Email:pd[12],
                                        AddressStreet1: "",
                                        AddressStreet2: "",
                                        AddressFirstName:"",
                                        AddressZip:pd[9],
                                        AddressState: "",
                                        AddressCountry: "",
                                        MobileNo: "",
                                        /*user account if buyer registered*/
                                        //Buyer *Cust //`json:"buyer"`
                                        //Cart Cart `json:"cart"`
                                        Cstid: cs.Id,
                                        Crtid: crtd,
                                        SecSign: string(secSign),
                                        Notified: true,
                                }
                        }
                        if byr.Firstname != "Guest" {
                                if pd[1] != "" && pd[1] == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleCitrusPost2")
                                        } else {
                                                k,err = addOrder(c, rdr)
                                                hndl(err, "handleCitrusPost3")
                                        }
                                }
                        } else { // checkout as guest
                                if pd[1] != "" && pd[1] == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleCitrusPost4")
                                        } else {
                                                k,err = addOrder(c, rdr)
                                                hndl(err, "handleCitrusPost5")
                                        }
                                }
                        }
                } else {
                        handleAccountError(c, w, r, 4, []byte("HMAC Error"))
                }
                //_,err = addOrder(c, rdr)
                //hndl(err, "handleAccountOrderP3")
                //handleAccountOrders(w, r)
                return
        }

        func handleAccountOrderCnfP(w http.ResponseWriter, r *http.Request ) { // test CITRUS return 
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                crtd := session.Values["crtd"].(int64)
                hndl(err, "handleAccountOrderCnfP0")
                lggd := session.Values["lggd"].(string)
                cs, err := getCstmr(c, lggd)
                hndl(err, "handleAccountOrderCnfP1")
                //crt := getCart(c, crtd, &cs)
                //hndl(err, "handleAccountOrderP1")
                // get posted form data
                err = r.ParseForm()
                hndl(err, "handleAccountOrderP2")
                //validation
                v := r.Form
                pd := make([]string,0)
                pd = append(pd, v.Get("TxId")) // 0
                pd = append(pd, v.Get("TxStatus")) // 1
                for _,p := range pd {
                        if p == "" {
                                handleAccountError(c, w, r, 4, []byte("Missing required data"))
                                return
                        }
                }
                pdi, err := strconv.ParseInt(pd[0], 10, 64)
                hndl(err, "handleAccountOrderCnfP2")
                var rdr Order
                flag := true
                if flag {
                        byr,err := getCstmr(c,lggd)
                        hndl(err, "handleAccountOrderCnfP3")
                        x,k,err := getOrder(c, pdi)
                        hndl(err, "handleAccountOrderCnfP4")
                        if k != nil {
                                x.Id = pdi  //1,000,000-9,999,999
                                x.TxId = pd[0]
                                x.TxStatus = "SUCCESS"
                                x.Cstid = cs.Id
                                x.Crtid = crtd
                                x.Notified = true
                        } else {
                                rdr = Order{
                                        Id:pdi,  //1,000,000-9,999,999
                                        TxId:pd[0],
                                        TxStatus:pd[1],
                                        Cstid: cs.Id,
                                        Crtid: crtd,
                                        Notified: true,
                                }
                        }
                        if byr.Firstname != "Guest" {
                                if pd[1] != "" && pd[1] == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleAccountOrderCnfP5")
                                                //success! archive purchased cart and empty user cart
                                                _, err := svPrchsdCart(c, cs.Email, crtd, pdi)
                                                hndl(err, "handleAccountOrderCnfP6")
                                                mptyCart(c, crtd)
                                        } else {
                                                k,err = addOrder(c, rdr)
                                                hndl(err, "handleAccountOrderCnfP7")
                                        }
                                }
                        } else { // checkout as guest
                                if pd[1] != "" && pd[1] == "SUCCESS" {
                                        if k != nil {
                                                x.Notified = true
                                                x,k,err = pdtOrder(c, x)
                                                hndl(err, "handleAccountOrderCnfP8")
                                                _, err := svPrchsdCart(c, cs.Email, crtd, pdi)
                                                hndl(err, "handleAccountOrderCnfP9")
                                                mptyCart(c, crtd)
                                        } else {
                                                k,err = addOrder(c, rdr)
                                                hndl(err, "handleAccountOrderCnfP10")
                                        }
                                }
                        }
                } else {
                        handleAccountError(c, w, r, 4, []byte("HMAC Error"))
                }
                //_,err = addOrder(c, rdr)
                //hndl(err, "handleAccountOrderP3")
                //handleAccountOrders(w, r)
                return
        }

        func handleAccountError(c appengine.Context, w http.ResponseWriter, r *http.Request, errtyp int, errmsg []byte) {
                var s0 []Good
                var c0 []Category
                c0, err := getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                cs, err := getCstmr(c, "Guest")
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                session, err := ckstore.Get(r, "account_path")
                handle(err)
                switch errtyp {
                case 1,2:
                        session.Values["lggd"] = "Guest"
                        session.Values["crtd"] = int64(0)
                        session.Save(r, w)
                        data := Render4{string(errmsg), cs, s0, c0, getCart(c,0,&cs)}
                        err = acc_errs[errtyp].ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                case 3:
                        cs,err = getCstmr(c, session.Values["lggd"].(string))
                        hndl(err, "handleAccountError0")
                        data := Render4{string(errmsg), cs, s0, c0, getCart(c, session.Values["crtd"].(int64), &cs)}
                        err = acc_errs[errtyp].ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                default:
                        session.Values["lggd"] = "Guest"
                        session.Values["crtd"] = int64(0)
                        session.Save(r, w)
                        data := Render4{string(errmsg), cs, s0, c0, getCart(c, 0, &cs)}
                        err = acc_errs[errtyp].ExecuteTemplate(w, "base", data)
                        handle(err)
                        return
                }
        }

	func logout(w http.ResponseWriter, r *http.Request) {
		session, err := ckstore.Get(r, "account_path")
		hndl(err, "logout0" )
		gdsession, err := ckstore.Get(r, "goods_path")
		hndl(err, "logout1")
		
		session.Values["lggd"] = "Guest"
		session.Values["crtd"] = int64(0)
		session.Options = &sessions.Options{
			Path:"/account/",
		}
		gdsession.Values["lggd"] = "Guest"
		gdsession.Values["crtd"] = int64(0)
		gdsession.Options = &sessions.Options{
			Path:"/goods/",
		}
		err = sessions.Save(r, w)
		hndl(err, "logout2")
		return
	}

// handlers for products display
        func handleCatalogwTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template , rndr interface{}) {
                c := appengine.NewContext(r)
                session, err := ckstore.Get(r, "goods_path")
                if err != nil {
                        http.Error(w, err.Error(), 500)
                        return
                }
                var s0 []Good
                var c0 []Category
                c0, err = getCategories(c)
                if err != nil {
                        srvErr(c, w, err)
                        return
                }
                ctg,_ := mux.Vars(r)["category"]
                sctg,_ := mux.Vars(r)["subcategory"]
                if session.Values["lggd"] != nil { // returning user
                        cs, err := getCstmr(c, session.Values["lggd"].(string))
                        crtd := session.Values["crtd"].(int64)
                        if err != nil {
                                srvErr(c, w, err)
                                return
                        }
                        switch rndr := rndr.(type) {
                        case []string:
                                s0,err = getSctgryGds(c, rndr)
                                data := Render5{ctg,sctg,cs, getCart(c, crtd, &cs), s0,c0,rndr}
                                err = tmpl.ExecuteTemplate(w,"base", data)
                                handle(err)
                                return
                        case Good:
                                data := Render6{ctg,sctg,cs,getCart(c, crtd, &cs), rndr,c0,rndr.Deets}
                                err = tmpl.ExecuteTemplate(w,"base", data)
                                handle(err)
                                return
                        case []Good:
                                data := Render5{ctg,sctg,cs,getCart(c, crtd, &cs), rndr,c0,nil}
                                err = tmpl.ExecuteTemplate(w,"base", data)
                                handle(err)
                                return
                        default:
                                data := Render4{"",cs,s0,c0,getCart(c,crtd,&cs)}
                                err = tmpl.ExecuteTemplate(w,"base", data)
                                handle(err)
                                return
                        }
                }
                ckchck(w, r) //if session is nil check whether cookies enabled
        }

        func handleCategory(w http.ResponseWriter, r *http.Request) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                ctg,_ := vars["category"]
                //handle(err)
                sctgs, err := getSubcats(c,ctg)
                handle(err)
                handleCatalogwTemplate(w, r, tmpl_cat_cat, sctgs)
        }

        func handleSubcategory(w http.ResponseWriter, r *http.Request) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                ctg,_ := vars["category"]
                sctg,_ := vars["subcategory"]
                //handle(err)
                gds, err := getGds(c,ctg,sctg)
                handle(err)
                handleCatalogwTemplate(w, r, tmpl_cat_sub_cat, gds)
        }

        func handleGoodDeets(w http.ResponseWriter, r *http.Request) {
                c := appengine.NewContext(r)
                vars := mux.Vars(r)
                gd,err := strconv.ParseInt(vars["id"],10,64)
                handle(err)
                gdd, err := getGood(c,gd)
                handle(err)
                handleCatalogwTemplate(w, r, tmpl_cat_gds_dts, gdd)
        }

        func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc{
                return func(w http.ResponseWriter, r *http.Request) { 
                        m := validPath.FindStringSubmatch(r.URL.Path)
                        if m == nil {
                                http.NotFound(w,r)
                                return
                        }
                        fn(w, r) //, m[2])
                }
        }

// init
        func init() {
                r := mux.NewRouter()
                s := r.PathPrefix("/admin").Subrouter()
                //s.HandleFunc("/", makeHandler(handleGoods)).Methods("GET")
                //s.HandleFunc("/{id}/", makeHandler(handleGood)).Methods("GET")
                //s.HandleFunc("/{id}/details", makeHandler(handleGoodDetails)).Methods("GET")
                s.HandleFunc("/search", makeHandler(handleGoodsSearch)).Methods("GET")
                s.HandleFunc("/goods", makeHandler(handleGoodsList)).Methods("GET")
                s.HandleFunc("/goods/json", makeHandler(handleGoodsListJSON)).Methods("GET")
                s.HandleFunc("/goods/entry", makeHandler(handleGoodEntry)).Methods("GET")
                s.HandleFunc("/goods/upload", makeHandler(handleGoodsEntry)).Methods("GET")
                s.HandleFunc("/goods/upload", makeHandler(handleGoodsCreate)).Methods("POST")
                s.HandleFunc("/accounts", makeHandler(handleAccountList)).Methods("GET")
                s.HandleFunc("/goods/entry/new", makeHandler(handleGoodCreate)).Methods("GET")
                s.HandleFunc("/orders", makeHandler(handleOrdersList)).Methods("GET")
                //s.HandleFunc("/entry/new", makeHandler(handleGoodCreateJS)).Methods("POST")
                //s.HandleFunc("/entry/new1", makeHandler(handleGoodCreateJS1)).Methods("POST")
                s.HandleFunc("/goods/edit/{id}", makeHandler(handleGoodEdit)).Methods("GET")
                //s.HandleFunc("/edit", makeHandler(handleGoodEditPostJS)).Methods("POST")
                s1 := r.PathPrefix("/account").Subrouter()
                s1.HandleFunc("/", makeHandler(handleAccount)).Methods("GET")
                s1.HandleFunc("/", makeHandler(handleAccountEdit)).Methods("POST")
                s1.HandleFunc("/cks", makeHandler(handleCookieCheck)).Methods("GET")
                s1.HandleFunc("/register", makeHandler(handleAccountRegister)).Methods("GET")
                s1.HandleFunc("/login", makeHandler(handleAccountLogin)).Methods("GET")
                s1.HandleFunc("/login", makeHandler(handleAccountLoginP)).Methods("POST")
                s1.HandleFunc("/logout", makeHandler(handleAccountLogout)).Methods("GET")
                s1.HandleFunc("/create", makeHandler(handleAccountCreateP)).Methods("POST")
                s1.HandleFunc("/forgot", makeHandler(handleAccountForgot)).Methods("GET")
                s1.HandleFunc("/forgot", makeHandler(handleAccountForgotP)).Methods("POST")
                s1.HandleFunc("/edit", makeHandler(handleAccountEdit)).Methods("GET")
                s1.HandleFunc("/edit", makeHandler(handleAccountEditP)).Methods("POST")
                s1.HandleFunc("/orders", makeHandler(handleAccountOrders)).Methods("GET")
                s1.HandleFunc("/orders/place", makeHandler(handleAccountOrderP)).Methods("POST")
                s1.HandleFunc("/orders/cnf", makeHandler(handleAccountOrderCnfP)).Methods("POST")
                s1.HandleFunc("/neworder", makeHandler(handleAccountNewOrderP)).Methods("POST")
                s1.HandleFunc("/returns", makeHandler(handleAccountReturns)).Methods("GET")
                s1.HandleFunc("/cart", makeHandler(handleAccountCart)).Methods("GET")
                s1.HandleFunc("/cart/add/{id}", makeHandler(handleAccountCartAdd)).Methods("GET")
                s1.HandleFunc("/cart/rm/{id}", makeHandler(handleAccountCartRm)).Methods("GET")
                s1.HandleFunc("/cart/pdt", makeHandler(handleAccountCartPdt)).Methods("POST")
                s1.HandleFunc("/cart/chckt", makeHandler(handleAccountChckt)).Methods("GET")
                s2 := r.PathPrefix("/goods").Subrouter()
                s2.HandleFunc("/search", makeHandler(handleGoodsSearch)).Methods("GET")
                s2.HandleFunc("/{category}", makeHandler(handleCategory)).Methods("GET")
                s2.HandleFunc("/{category}/{subcategory}", makeHandler(handleSubcategory)).Methods("GET")
                s2.HandleFunc("/{category}/{subcategory}/{id}", makeHandler(handleGoodDeets)).Methods("GET")
                s3 := r.PathPrefix("/promo").Subrouter()
                s3.HandleFunc("/", makeHandler(handleMainPage))
                s3.HandleFunc("/reg", makeHandler(handleRegPage))
                s3.HandleFunc("/how", makeHandler(handleHowPage))
                s3.HandleFunc("/flock", makeHandler(handleFlockPage))
                s3.HandleFunc("/confirm/", makeHandler(handleConfirmPage))
                http.Handle("/", r)
        }

        func loadInitial() {
                raw,err := ioutil.ReadFile("categories")
                handle(err)
                var categories []Category
                err = json.Unmarshal(raw, &categories)
                handle(err)
        }

// datastore access functions
        func nxtOrderId(c appengine.Context) int64 {
                ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                gi := new(Counter)
                if err := datastore.Get(c, ik, gi); err != nil {
                        gi.TtlGd = 1000
                        gi.NxtGd = 1001
                        gi.TtlCst = 10000
                        gi.NxtCst = 10001
                        gi.TtlRdr = 1000000
                        gi.NxtRdr = 1000001
                        if _, err1 := datastore.Put(c, ik, gi); err1 != nil {
                                //http.Error(w, err.Error(), 500)
                                hndl(err1, "nxtOrderId0")
                                return 1000001
                        }
                        return gi.NxtRdr
                }
                if len(gi.HlsRdr) == 0 {
                        tmp := gi.NxtRdr
                        gi.NxtRdr++
                        _,err := datastore.Put(c, ik, gi)
                        hndl(err, "nxtOrderId1")
                        return tmp
                } else {
                        tmp := gi.HlsRdr[0] // re-use oldest incomplete order id
                        gi.HlsRdr = gi.HlsRdr[1:]
                        _,err := datastore.Put(c, ik, gi)
                        hndl(err, "nxtOrderId2")
                        return tmp
                }
        }

        func getOrder(c appengine.Context, id int64) (Order, *datastore.Key, error) {
                var rdrs []Order
                var rdr Order
                q0 := datastore.NewQuery("order").Filter("Id = ", id)
                k, err := q0.GetAll(c, &rdrs)
                hndl(err, "getOrder0")
                if k != nil {
                        return rdrs[0], k[0], nil
                }
                return rdr,nil,err
        }

        func getOrders(c appengine.Context, cd int64) ([]Order, []*datastore.Key, error) {
                var rdrs []Order
                //var rdr Order
                q0 := datastore.NewQuery("order").Filter("Cstid = ", cd)
                k, err := q0.GetAll(c, &rdrs)
                hndl(err, "getOrders")
                if k != nil {
                        return rdrs, k, nil
                }
                return rdrs,nil,err
        }

        func getCmpltOrders(c appengine.Context, cstml string) ([]Order, error) {
                cs,err := getCstmre(c, cstml)
                hndl(err, "getCmpltOrders0")
                q := datastore.NewQuery("order").Filter("TxStatus = ", "SUCCESS").Filter("Cstid = ", cs.Id)
                var crt []Order
                k,err := q.GetAll(c, &crt)
                if k != nil {
                        return crt, nil
                }
                return crt, err
        }
        func addOrder(c appengine.Context, rdr Order) (*datastore.Key, error) { // doesn't provide error about existing order - just returns. Use pdtOrder to update orders
                var rdrs []Order
                q0 := datastore.NewQuery("order").Filter("Id = ", rdr.Id)
                ks, err := q0.GetAll(c, &rdrs)
                hndl(err, "addOrder0")
                if ks != nil { // order id already exists in datastore - just return it. Don't update Order in datastore
                        return ks[0], err
                }
                k, err := datastore.Put(c, datastore.NewKey(c, "order", "", rdr.Id, nil), &rdr) //assumes order id is next available
                hndl(err, "addOrder1")
                if k != nil {
                        ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                        gi := new(Counter)
                        err = datastore.Get(c, ik, gi)
                        hndl(err, "addOrder2")
                        //gi.TtlRdr += 1
                        gi.NxtRdr += 1
                        _, err := datastore.Put(c, ik, gi)
                        hndl(err, "addOrder3")
                        return k, err
                }
                return nil, err
        }

        func pdtOrder(c appengine.Context, rdr Order) (Order, *datastore.Key, error) {
                _, k, err := getOrder(c, rdr.Id)
                hndl(err, "pdtOrder0")
                if k != nil {
                        k, err = datastore.Put(c, k, &rdr)
                        hndl(err, "pdtOrder1")
                        if k != nil {
                                return rdr, k, nil
                        }
                        return rdr, k, err
                }
                return rdr, nil, err
        }

        func rmOrder(c appengine.Context, rdrd int64) bool {
                var rdrs []Order
                q0 := datastore.NewQuery("order").Filter("Id = ", rdrd)
                k, err := q0.GetAll(c, &rdrs)
                hndl(err, "rmOrder0")
                if k != nil {
                        err = datastore.Delete(c, k[0])
                        hndl(err, "rmOrder1")
                        ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                        gi := new(Counter)
                        err = datastore.Get(c, ik, gi)
                        hndl(err, "rmOrder2")
                        gi.HlsRdr = append(gi.HlsRdr, rdrd)
                        return true
                }
                return false
        }

        func cnfOrder(c appengine.Context, rdr Order) (Order, *datastore.Key, error) {
                x, k, err := pdtOrder(c, rdr)
                hndl(err, "cnfOrder0")
                if k != nil {
                        // order's been confirmed so increase total orders counter
                        ik := datastore.NewKey(c, "counter", "thekey", 0, nil)
                        gi := new(Counter)
                        err = datastore.Get(c, ik, gi)
                        hndl(err, "cnfOrder1")
                        gi.TtlRdr += 1
                        _, err := datastore.Put(c, ik, gi)
                        hndl(err, "cnfOrder2")
                        return x, k, nil
                }
                return x, k, err
        }

        func getCategories(c appengine.Context)  ([]Category, error){
                var cs []Category
                q0 := datastore.NewQuery("category")
                key, err := q0.GetAll(c, &cs)
                handle(err)
                if key != nil {
                        return cs, nil
                }else {
                        raw,err := ioutil.ReadFile("categories")
                        handle(err)
                        var categories []Category
                        err = json.Unmarshal(raw, &categories)
                        handle(err)
                        for _,cat := range categories {
                          if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"category", nil), &cat); err != nil {
                                handle(err)
                          }
                        }
                        _, err = q0.GetAll(c, &cs)
                        handle(err)
                        return cs, nil
                        //return nil, errors.New("Something's wrong")
                }
        }

        func getSubcats(c appengine.Context, ctgry string) ([]string, error) {
                q := datastore.NewQuery("category").Filter("Name = ", ctgry)
                var ctgrs []Category
                var sctgrs []string
                key, err := q.GetAll(c, &ctgrs)
                if key != nil {
                        return ctgrs[0].Subcategories, nil
                }
                if err != nil {
                        return sctgrs, err
                }
                return sctgrs, nil
        }

        func getGds(c appengine.Context, ctgry string, sctgry string) ([]Good, error) {
                q := datastore.NewQuery("product").Filter("Category = ", ctgry).Filter("Subcategory = ", sctgry)
                var gds []Good
                key, err := q.GetAll(c, &gds)
                if key != nil {
                        return gds, nil
                }
                if err != nil {
                        return gds, err
                }
                return gds, nil
        }

        func getSctgryGds(c appengine.Context, sctgrys []string) ([]Good, error) {
                q := datastore.NewQuery("product")
                var gds []Good
                gds1 := make([]Good, len(sctgrys))
                for i,sc := range sctgrys {
                        q.Filter("Subcategory = ", sc)
                        key, err := q.GetAll(c, &gds)
                        if key != nil {
                                gds1[i] = gds[0]
                        }
                        handle(err)
                }
                return gds1, nil
        }

        func getGood(c appengine.Context, yd int64) (Good, error) {
                q := datastore.NewQuery("product").Filter("Id = ", yd)
                var gds []Good
                key, err := q.GetAll(c, &gds)
                if key != nil {
                        return gds[0], nil
                }
                if err != nil {
                        return Good{}, err
                }
                return Good{}, nil
        }

        func getCstmr(c appengine.Context, firstname string)  (Cust, error){
                var cs []Cust
                q0 := datastore.NewQuery("cstmr").Filter("Firstname =", firstname)
                key, err := q0.GetAll(c, &cs)
                handle(err)
                if key != nil {
                        return cs[0], nil
                }else { //return guest
                        q0 := datastore.NewQuery("cstmr").Filter("Firstname =", "Guest")
                        key, err := q0.GetAll(c, &cs)
                        handle(err)
                        if key != nil {
                                return cs[0], nil
                        } else {
                                raw,err := ioutil.ReadFile("cstmrs")
                                handle(err)
                                var cstmrs []Cust
                                err = json.Unmarshal(raw, &cstmrs)
                                handle(err)
                                for _,cstmr := range cstmrs {
                                        if _, err := datastore.Put(c, datastore.NewKey(c,"cstmr","",0,nil), &cstmr); err != nil {
                                                handle(err)
                                        }
                                }
                                _, err = q0.GetAll(c, &cs)
                                handle(err)
                                return cs[0], nil
                        }
                        //return nil, errors.New("Something's wrong")
                }
        }

        func getCstmre(c appengine.Context, email string)  (Cust, error){
                var cs []Cust
                q0 := datastore.NewQuery("cstmr").Filter("Email =", email)
                key, err := q0.GetAll(c, &cs)
                handle(err)
                if key != nil {
                        return cs[0], nil
                }else { //return
                        q0 := datastore.NewQuery("cstmr").Filter("Firstname =", "Guest")
                        key, err := q0.GetAll(c, &cs)
                        handle(err)
                        if key != nil {
                                return cs[0], nil
                        } else {
                                raw,err := ioutil.ReadFile("cstmrs")
                                handle(err)
                                var cstmrs []Cust
                                err = json.Unmarshal(raw, &cstmrs)
                                handle(err)
                                for _,cstmr := range cstmrs {
                                        if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"cstmr", nil), &cstmr); err != nil {
                                                handle(err)
                                        }
                                }
                                _, err = q0.GetAll(c, &cs)
                                handle(err)
                                return cs[0], nil
                        }
                        //return nil, errors.New("Something's wrong")
                }
        }

        func getCstmrk(c appengine.Context, email string)  (Cust, *datastore.Key, error){
                var cs []Cust
                q0 := datastore.NewQuery("cstmr").Filter("Email =", email)
                key, err := q0.GetAll(c, &cs)
                handle(err)
                if key != nil {
                        return cs[0], key[0], nil
                }else {
                        q0 := datastore.NewQuery("cstmr").Filter("Firstname =", "Guest")
                        key, err := q0.GetAll(c, &cs)
                        handle(err)
                        if key != nil {
                                return cs[0], key[0], nil
                        } else {
                                raw,err := ioutil.ReadFile("cstmrs")
                                handle(err)
                                var cstmrs []Cust
                                err = json.Unmarshal(raw, &cstmrs)
                                handle(err)
                                for _,cstmr := range cstmrs {
                                        if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"cstmr", nil), &cstmr); err != nil {
                                                handle(err)
                                        }
                                }
                                key, err = q0.GetAll(c, &cs)
                                handle(err)
                                return cs[0],key[0],nil
                        }
                        //return nil, errors.New("Something's wrong")
                }
        }

        func getCstmr1(c appengine.Context, w http.ResponseWriter, r *http.Request)  (Cust, error){
                var dmmy Cust
                session, err := ckstore.Get(r, "account_path")
                if err != nil {
                        return dmmy, err //http.Error(w, err.Error(), 500)
                }
                var firstname string
                if session.Values["lggd"] != nil {
                        firstname = session.Values["lggd"].(string)
                } else {
                        firstname = "Guest"
                }
                var cs []Cust
                q0 := datastore.NewQuery("cstmr").Filter("Firstname =", firstname)
                key, err := q0.GetAll(c, &cs)
                handle(err)
                if key != nil {
                        return cs[0], nil
                }else { //return guest
                        q0 := datastore.NewQuery("cstmr").Filter("Firstname =", "Guest")
                        key, err := q0.GetAll(c, &cs)
                        handle(err)
                        if key != nil {
                                return cs[0], nil
                        } else {
                                raw,err := ioutil.ReadFile("cstmrs")
                                handle(err)
                                var cstmrs []Cust
                                err = json.Unmarshal(raw, &cstmrs)
                                handle(err)
                                for _,cstmr := range cstmrs {
                                        if _, err := datastore.Put(c, datastore.NewIncompleteKey(c,"cstmr", nil), &cstmr); err != nil {
                                                handle(err)
                                        }
                                }
                                _, err = q0.GetAll(c, &cs)
                                handle(err)
                                return cs[0], nil
                        }
                        //return nil, errors.New("Something's wrong")
                }
        }

        func mkCart(c appengine.Context, cstml string) int64 {
                cc,cck,err := getCstmrk(c, cstml)
                hndl(err,"mkCart1")
                if cck != nil {
                        if cc.Crtid > 0  && getCart(c, cc.Crtid, &cc).Id > 0 {
                                return cc.Crtid
                        }
                        crt := &Cart{cc.Crtid, 0,make([]int64,0), make([]int,0), false, make([]float64,0), 0, 0}
                        crtk := datastore.NewKey(c, "cart","", cc.Crtid, cck)
                        crtk, err := datastore.Put(c, crtk, crt) //re-use crtk
                        hndl(err,"mkCart1")
                        if crtk != nil {
                                cc.Crtid = (*crt).Id
                                crtk, err = datastore.Put(c, cck, cc)
                                hndl(err,"mkCart1")
                                if crtk != nil {
                                        return cc.Id
                                }
                        }
                }
                return 0
        }

        func svPrchsdCart(c appengine.Context, cstml string, crtd int64, rdrd int64) (bool, error) { // duplicates cart with id w 9 digits
                cs, err := getCstmre(c, cstml)
                hndl(err, "svPrchsdCart0")
                crt, k := getCartk(c, crtd, &cs)
                crt.Id = rdrd * 100
                crt.Cstid = cs.Id
                crt.Prchsd = true
                if k != nil {
                        log.Println("Saving purchased cart: ", crt.Id)
                        nk := datastore.NewKey(c, "cart", "", rdrd*100, nil)
                        k, err = datastore.Put(c, nk, &crt)
                        hndl(err, "svPrchsdCart1")
                        return true, nil
                }
                return false, err
        }

        func getCart(c appengine.Context, crtd int64, cst *Cust) Cart {
                log.Println("Getting: ", crtd, cst.Firstname)
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crt []Cart
                var crtr Cart
                key, err := q.GetAll(c, &crt )
                //log.Println("Cart id: ", crt[0])
                if key != nil { //found something
                        if cst.Crtid == crt[0].Id { //check if cart belongs to session user
                                return crt[0]
                        } else {
                                return crt[0]
                        }
                }
                hndl(err,"getCart0")
                return crtr
        }

        func getCartk(c appengine.Context, crtd int64, cst *Cust) (Cart, *datastore.Key) {
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crt []Cart
                var crtr Cart
                key, err := q.GetAll(c, &crt )
                if key != nil { //found something
                        if cst.Crtid == crt[0].Id { //check if cart belongs to session user
                                return crt[0], key[0]
                        } else {
                                return crt[0], key[0]
                        }
                }
                hndl(err,"getCart0")
                return crtr, nil
        }

        func addCart(c appengine.Context, crtd int64, gd int64) Cart {
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crts []Cart
                key, err := q.GetAll(c, &crts)
                if key != nil { //found something
                        crts[0].Ids = append(crts[0].Ids, gd)
                        crts[0].Qntts = append(crts[0].Qntts, 1)
                        gdd,err := getGood(c,gd)
                        handle(err)
                        crts[0].Ttls = append(crts[0].Ttls, gdd.Price * 1)
                        crts[0].Ttl += gdd.Price * 1
                        _,err = datastore.Put(c,key[0],&crts[0])
                }
                handle(err)
                return crts[0]
        }

        func pdtCart(c appengine.Context, crtd int64, qts []int) Cart {
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crts []Cart
                key, err := q.GetAll(c, &crts)
                if key != nil { //found something
                        if len(qts)==len(crts[0].Ids) {
                                crts[0].Ttl = 0
                                for i,gd := range crts[0].Ids { //gds {
                                        //crts[0].Ids[i] = gds[i]
                                        crts[0].Qntts[i] = qts[i]
                                        gdd,err := getGood(c,gd)
                                        hndl(err,"pdtCart0")
                                        crts[0].Ttls[i] = gdd.Price * float64(qts[i])
                                        crts[0].Ttl += gdd.Price * float64(qts[i])
                                }
                                _,err = datastore.Put(c,key[0],&crts[0])
                        }
                }
                hndl(err,"pdtCart1")
                return crts[0]
        }

        func rmCart(c appengine.Context, crtd int64, gd int64) Cart {
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crts []Cart
                key, err := q.GetAll(c, &crts)
                if key != nil { //found something
                        for i,g := range crts[0].Ids {
                                if int64(g) == gd {
                                        crts[0].Ids = append(crts[0].Ids[:i], crts[0].Ids[i+1:]...)
                                        crts[0].Qntts = append(crts[0].Qntts[:i], crts[0].Qntts[i+1:]...)
                                        crts[0].Ttl -= crts[0].Ttls[i]
                                        crts[0].Ttls = append(crts[0].Ttls[:i], crts[0].Ttls[i+1:]...)
                                        _,err = datastore.Put(c,key[0],&crts[0])
                                }
                        }
                }
                hndl(err,"rmCart1")
                return crts[0]
        }

        func mptyCart(c appengine.Context, crtd int64) Cart {
                q := datastore.NewQuery("cart").Filter("Id = ", crtd)
                var crts []Cart
                key, err := q.GetAll(c, &crts)
                if key != nil { //found something
                        crt := &Cart{crts[0].Id, 0,make([]int64,0), make([]int,0), false, make([]float64,0), 0, 0}
                        _,err = datastore.Put(c, key[0], crt)
                }
                hndl(err,"rmCart1")
                return crts[0]
        }

        func getHstCrts(c appengine.Context, cstml string) ([]Cart, error) { // historical carts
                cs,err := getCstmre(c, cstml)
                hndl(err, "getHstr0")
                q := datastore.NewQuery("cart").Filter("Id > ", 100000000).Filter("Cstid = ", cs.Id)
                var crt []Cart
                k,err := q.GetAll(c, &crt)
                if k != nil {
                        return crt, nil
                }
                return crt, err
        }

        func getGoods(c appengine.Context, ids []int64) []Good {
                var crt, crt1 []Good
                for _, id := range ids {
                        q := datastore.NewQuery("product").Filter("Id = ", id)
                        key, err := q.GetAll(c, &crt )
                        if key != nil { //found something
                                crt1 = append(crt1, crt[0])
                        }
                        hndl(err,"getGoods0")
                }
                return crt1
        }

// utils
        func ckchck(w http.ResponseWriter, r *http.Request) {
                ck := new(http.Cookie)
                ck.Name = "COOKIES"
                ck.Value = "YES"
                ck.Path = "/account/cks"
                ck.Expires = time.Now().Add(time.Minute)
                http.SetCookie(w,ck)
                http.Redirect(w, r, "/account/cks", 307)
                return
        }

        func checkMac(msg []byte, msgMAC []byte, key []byte) bool {
                mac := hmac.New(sha1.New,key)
                mac.Write(msg)
                xpctdMAC := mac.Sum(nil)
                return hmac.Equal(msgMAC, xpctdMAC)
        }

        func hMac(msg []byte, key []byte) []byte {
                mac := hmac.New(sha1.New, key)
                mac.Write(msg)
                src := mac.Sum(nil)
                dst := make([]byte, hex.EncodedLen(len(src)))
                hex.Encode(dst, src)
                return dst
        }

// error handlers
        func handle(err error) {
                if err != nil {
                        log.Println("Error: %v", err)
                }
        }

        func hndl(err error, fnctn string) {
                if err != nil {
                        log.Println("Error in %v: %v", fnctn, err)
                }
        }

        func srvErr(ctx appengine.Context, w http.ResponseWriter, err error) {
                w.WriteHeader(http.StatusInternalServerError)
                w.Header().Set("Content-Type", "text/plain")
                io.WriteString(w, "Internal Server Error")
                log.Fatalf("Server error: %v", err)
                //http.Error(w, err.Error(), http.StatusInternalServerError)
        }
