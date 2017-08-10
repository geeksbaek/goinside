package goinside

func getTestGuestSession() (ms *GuestSession, err error) {
	return Guest("ㅇㅇ", "123")
}
