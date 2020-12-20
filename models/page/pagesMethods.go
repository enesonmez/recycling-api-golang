package page

import (
	. "../info"
	. "../user"
)

func (infoPage *InfoPage) InfoPageConstructer(info *Info) {
	infoPage.Information = *info
}

func (userPage *UserPage) UserPageConstructer(info *Info, user *User) {
	userPage.Information = *info
	userPage.User = *user
}
