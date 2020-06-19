package main

type NoAvailableServerError string

func (e NoAvailableServerError) Error() string {
	return string(e)
}
