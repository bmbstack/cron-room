package models

type BmbUser struct {
	ID            int64   `json:"-" gorm:"column:id; type:int(11); primary_key; auto_increment"`
	UserID        string  `json:"userID" gorm:"column:user_id; type:varchar(50); not null; unique_index:uix_bmb_user_user_id; index:idx_bmb_user_user_id"`
	Username      string  `json:"username" gorm:"column:username; type:varchar(20); not null"`
	HeadPhoto     string  `json:"headPhoto" gorm:"column:head_photo; type:varchar(200); not null"`
	Sex           int64   `json:"sex" gorm:"column:sex; type:tinyint(1); not null; default:1"`
	Pid           int64   `json:"pid" gorm:"column:pid; type:tinyint(4); not null; default:1"`
	Cid           int64   `json:"cid" gorm:"column:cid; type:tinyint(4); not null; default:1"`
	BaseModel
}

func init() {
	RegisterModels(&BmbUser{})
}

func (BmbUser) TableName() string {
	return "bmb_user"
}

func (this *BmbUser) FindUserByUserID(userID string) (one BmbUser) {
	DB.Where("user_id=? AND is_deleted=0", userID).First(&one)
	return one
}
