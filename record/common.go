package record

import (
	"gorm.io/gorm"
)

//func (r *Record) GetID() uuid.UUID {
//	//
//	if r.cachedIDSet.IfFalseSetTrue() {
//		return r.cachedID
//	}
//	r.cachedID = _struct.GetFieldValue(r.NonPtrObj, "ID").(uuid.UUID)
//	return r.cachedID
//}

func (r *Record) GetUserID() interface{} {
	return r.userID
}

func (r *Record) GetDeviceID() interface{} {
	return r.deviceID
}

/*func (r *Record) FieldExists() bool {
	// most of the time the ID will be of type uuid.UUID
	// but there may be also other types...
	if r.GetID() == uuid.Nil {
		return false
	}
	return true
}*/

func (r *Record) getDB() *gorm.DB {
	/*
		We should check if there is a context here...
		if yes then we can get the DB client with this context!
		if there is no context, then we just get the default client
	*/
	// TODO: we can handle the recover here?!

	if r.dbSet.IfFalseSetTrue() {
		return r.DB
	}

	// if r.DB == nil {
	// 	r.DB = db.DB()
	// }
	if r.Ctx != nil {
		// Get DB with context
		r.DB = r.DB.WithContext(r.Ctx)
	}

	return r.DB
}
