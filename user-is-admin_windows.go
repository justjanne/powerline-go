// +build windows

package main

import (
	"golang.org/x/sys/windows"
)

func userIsAdmin() bool {
	var sid *windows.SID

	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	t := windows.Token(0)

	member, err := t.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}
