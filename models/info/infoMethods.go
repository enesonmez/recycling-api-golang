package info

func (info *Info) InfoConstructer(isSuccess bool, message string) {
	info.IsSuccess = isSuccess
	info.Message = message
}
