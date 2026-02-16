package fixtures

type FKUser struct {
	ID       int64  `db:"pk"`
	Username string `db:"unique"`
	Email    string
	Active   bool   `db:"default:true"`
}

type Order struct {
	ID        int64     `db:"pk"`
	UserID    int64     `db:"fk:users,id,on_delete:CASCADE"`
	ProductID int64     `db:"fk:products,id,on_delete:RESTRICT"`
	Status    string    `db:"enum:pending,processing,shipped,delivered"`
	CreatedAt int64     `db:"default:extract(epoch from now())"`
}

type OrderItem struct {
	ID        int64  `db:"pk"`
	OrderID   int64  `db:"fk:orders,id,on_delete:CASCADE,on_update:CASCADE"`
	ProductID int64  `db:"fk:products,id,on_delete:SET_NULL"`
	Quantity  int    `db:"check:quantity > 0"`
	Price     float64
}

type Comment struct {
	ID        int64  `db:"pk"`
	PostID    int64  `db:"fk:posts,id"`
	UserID    int64  `db:"fk:users,id,on_delete:NO_ACTION"`
	Content   string `db:"check:length(content) > 0"`
}
