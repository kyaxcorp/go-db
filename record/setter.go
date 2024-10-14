package record

func (r *Record) SetUserID(userID interface{}) *Record {
	r.userID = userID
	return r
}

func (r *Record) SetDeviceID(deviceID interface{}) *Record {
	r.deviceID = deviceID
	return r
}
