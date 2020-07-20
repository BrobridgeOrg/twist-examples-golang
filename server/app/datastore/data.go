package datastore

func CreateAccount(name string, initBalance int) {
	DataBalance[name] = initBalance
	DataReserve[name] = 0
	DataUser[name] = true
}

var DataBalance = make(map[string]int)
var DataReserve = make(map[string]int)
var DataUser = make(map[string]bool)
