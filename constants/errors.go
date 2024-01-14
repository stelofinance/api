package constants

const (
	ErrorG000 = "G000: Invalid JSON Body"
	ErrorG001 = "G001: Invalid URL Parameter"
	ErrorG002 = "G002: Invalid URL Query"

	ErrorS000 = "S000: Unexpected Server Error"

	ErrorU000 = "U000: Username Taken"
	ErrorU001 = "U001: Invalid Username or Password"
	// ErrorU002 = "U002: Wallet Not Found"
	ErrorU003 = "U003: User Not Found"
	ErrorU004 = "U004: Session Not Found"

	ErrorW000 = "W000: Wallet Not Found"
	// ErrorW001 = "W001: Asset Not Found"
	ErrorW002 = "W002: Wallet or Asset Not Found"
	ErrorW003 = "W003: User Already Assigned To Wallet"
	// ErrorW004 = "W004: Transaction(s) Not Found"
	ErrorW005 = "W005: Session Not Found"
	ErrorW006 = "W006: Wallet Not Owned"
	ErrorW007 = "W007: User Not Assigned"
	ErrorW008 = "W008: Can't Delete Primary Wallet"
	// ErrorW009 = "W009: Can't Add Webhook To Primary Wallet"
	ErrorW010 = "W010: Webhook Errored"
	ErrorW011 = "W011: Address Taken"
	ErrorW012 = "W012: Can't Transfer Primary Wallet"

	ErrorA000 = "A000: Session Required"
	ErrorA001 = "A001: Invalid Session"
	ErrorA002 = "A002: Guests Only"
	ErrorA003 = "A003: Invalid Admin Key"

	ErrorI000 = "I000: Asset Not Found"
	ErrorI001 = "I001: Invalid Asset"

	ErrorP000 = "P000: Invalid Auth Form"
	ErrorP001 = "P001: Invalid Channel"

	ErrorH000 = "H000: Can't Create Warehouses"
)
