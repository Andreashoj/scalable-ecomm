package domain

type ProductSearch struct {
	Order Order
	Sort  Sort
}

type Order string
type Sort string

const (
	OrderAscending  Order = "ascending"
	OrderDescending Order = "descending"
)

const (
	SortDate Sort = "date"
)

func (s *Sort) IsValid() bool {
	return true
}

func (s *Sort) ToSQL() string {
	return "created_at"
}

func (s *Order) IsValid() bool {
	return true
}

func (s *Order) ToSQL() string {
	if *s == OrderAscending {
		return "ASC"
	}

	if *s == OrderDescending {
		return "DESC"
	}

	return "ASC"
}

func NewProductSearch(order, sort string) *ProductSearch {
	return &ProductSearch{
		Order: Order(order),
		Sort:  Sort(sort),
	}
}
