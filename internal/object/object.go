package object

type Obj struct {
	St  *Settings // settings
	Sts *Statuses
	Ck  *CookieInfo
}

func (o *Obj) NewObjects(st *Settings, sts *Statuses, ck *CookieInfo) {
	if st == nil {
		o.St = NewSettings()
	} else {
		o.St = st
	}
	if sts == nil {
		o.Sts = NewStatuses()
	} else {
		o.Sts = sts
	}
	if ck == nil {
		o.Ck = NewCookieInfo()
	} else {
		o.Ck = ck
	}
}
