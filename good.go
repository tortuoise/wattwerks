package wattwerks

import (
        _"appengine/datastore"

)

type Category struct {
        Name string `json:"category"`
        Subcategories []string `json:"subcategories"`
}
type Good struct {
        Id int64 `json:"id,string" schema:"-"` //1,000-9,999
        Code string `json:"code" schema:"code"`
        Category string `json:"category" schema:"category"`
        Subcategory string `json:"subCategory" schema:"subcategory"`
        Brand string `json:"brand" schema:"brand"`
        Desc string `json:"desc" schema:"desc"`
        Price float64 `json:"price,string" schema:"price"`
        Url string `json:"url" schema:"url"`
        Urlimg string `json:"urlImg" schema:"urlImg"`
        Featured bool `json:"featured,string" schema:"featured"`
        Hidden bool `json:"hidden,string" schema:"hidden"`
        Deets GoodDeets `json:"goodDeets" schema:"Deets"`
}
type GoodDeets struct {
        //Id int64 `json:"id"`
        DescDetails string `json:"descDetails" schema:"descDetails"`//path to file with details
        Tax float64 `json:"tax,string" schema:"tax"`//percent
        Price float64 `json:"price,string" schema:"price"`
        Stock int `json:"stock,string" schema:"stock"`
        Related []int64 `json:"related" cap:"6" schema:"-"`
        Prices []float64 `json:"prices" cap:"6" schema:"-"`
        Volumes []int `json:"volumes" cap:"6" schema:"-"`
        //PriceVolume map[int]float64 `json:"priceVolume"`
        ParameterNames []string `json:"parameterNames" cap:"12" schema:"-"`
        ParameterValues []string `json:"parameterValues" cap:"12" schema:"-"`
        //Parameters map[string]string `json:"parameters"`
        Features []string `json:"features" cap:"12" schema:"-"`
        Items []string `json:"items" cap:"6" schema:"-"`
        UrlImgs1 string `json:"urlImgs1" schema:"urlImgs"`
        UrlImgs2 string `json:"urlImgs2" schema:"-"`
        UrlImgs3 string `json:"urlImgs3" schema:"-"`
        UrlFile string `json:"urlFile" schema:"urlFile"`
}
type RelatedGoods struct {
        Gd Good `json:"good,string"`
        Rgds []Good `json:"goods,string"`
}


type GoodDatabase interface {

        ListGoods() ([]*Good, error)

        Add(good *Good) (int64, error)

        Get(id int64) (*Good, error)

        GetGoodByCategory(cat Category) ([]*Good, error)

        GetGoodBySubcategory(subcat string) ([]*Good, error)

        Update(good *Good) (int64, error)

        Delete(id int64) error

        Close() error

}
