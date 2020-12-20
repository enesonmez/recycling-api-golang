package info

func (info *Info) InfoConstructer(isSuccess bool, status int, message string) {
	info.IsSuccess = isSuccess
	info.Status = status
	info.Message = message
}
