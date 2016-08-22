package wattwerks

import (
        _"appengine/datastore"

)

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
        Related []int64 `json:"related" cap:"6"`
        Prices []float64 `json:"prices" cap:"6"`
        Volumes []int `json:"volumes" cap:"6"`
        //PriceVolume map[int]float64 `json:"priceVolume"`
        ParameterNames []string `json:"parameterNames" cap:"12"`
        ParameterValues []string `json:"parameterValues" cap:"12"`
        //Parameters map[string]string `json:"parameters"`
        Features []string `json:"features" cap:"12"`
        Items []string `json:"items" cap:"6"`
        UrlImgs1 string `json:"urlImgs1"`
        UrlImgs2 string `json:"urlImgs2"`
        UrlImgs3 string `json:"urlImgs3"`
        UrlFile string `json:"urlFile"`
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
