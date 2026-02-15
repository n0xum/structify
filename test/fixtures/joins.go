package fixtures

type JoinUser struct {
	ID       int64  `db:"pk"`
	Username string `db:"unique"`
	Email    string
	Active   bool   `db:"default:true"`
}

type JoinProduct struct {
	ID          int64   `db:"pk"`
	Name        string
	Price       float64 `db:"check:price >= 0"`
	Description string
}

type JoinOrder struct {
	ID        int64  `db:"pk"`
	UserID    int64  `db:"fk:join_user,id,on_delete:CASCADE"`
	Status    string `db:"enum:pending,processing,shipped,delivered,cancelled"`
	CreatedAt int64  `db:"default:extract(epoch from now())"`
}

type JoinOrderItem struct {
	ID        int64   `db:"pk"`
	OrderID   int64   `db:"fk:join_order,id,on_delete:CASCADE"`
	ProductID int64   `db:"fk:join_product,id,on_delete:SET_NULL"`
	Quantity  int     `db:"check:quantity > 0,default:1"`
	Price     float64
}
