package wattwerks

import (

)

type Counter struct {
        TtlGd int64 // 1000-9999
        NxtGd int64 // next lowest available ID
        TtlCst int64 // 10,000-99,999
        NxtCst int64
        TtlRdr int64
        NxtRdr int64
        HlsRdr []int64 // incomplete order ids that get returned from PG wo completion/confirmation
}

