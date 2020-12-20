package page

import (
	. "../info"
	. "../user"
)

type InfoPage struct {
	Information Info
}

type UserPage struct {
	Information Info
	User        User
}
