package dao

type ListRequest struct {
	Page     int
	PageSize int
}

type ProductListRequest struct {
	ListRequest
}
