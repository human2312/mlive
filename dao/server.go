package dao

var (
	Db *Dao
)

// 初始化
func NewDao()  {

	Db = New()
}