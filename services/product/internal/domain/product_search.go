package domain

type ProductSearch struct {
	Order  Order
	Sort   Sort
	Filter []string
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

func (s Sort) IsValid() bool {
	switch s {
	case SortDate:
		return true
	default:
		return false
	}
}

func (s Sort) ToSQL() string {
	return "created_at"
}

func (s Order) IsValid() bool {
	switch s {
	case OrderAscending:
		return true
	case OrderDescending:
		return true
	default:
		return false
	}
}

func (s Order) ToSQL() string {
	switch s {
	case OrderAscending:
		return "ASC"
	case OrderDescending:
		return "DESC"
	default:
		return "ASC"
	}
}

func NewProductSearch(order, sort string, filters []string) *ProductSearch {
	return &ProductSearch{
		Order:  Order(order),
		Sort:   Sort(sort),
		Filter: filters,
	}
}
