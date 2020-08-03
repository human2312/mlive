package shop

var (
	shop = Shop{}
)

func List() ([]Shop) {
	shop.UserId = 1
	return shop.List()
}

func Add()  {
	shop.UserId = 1
	shop.Add()
}

func Edit()  {
	
}

func Del()  {
	
}
