package testutil

type FakeMutex struct{}

func (m FakeMutex) Lock()   {}
func (m FakeMutex) Unlock() {}
