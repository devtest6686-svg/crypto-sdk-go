package consts

type CallbackEventType string

// 事件类型常量 TODO:
const (
	EventTypeTransactionCreated CallbackEventType = "transaction_created"
	EventTypeTransactionUpdated CallbackEventType = "transaction_updated"
	EventTypeWalletCreated      CallbackEventType = "wallet_created"
	EventTypeWalletUpdated      CallbackEventType = "wallet_updated"
	EventTypeAddressGenerated   CallbackEventType = "address_generated"
	EventTypeTest               CallbackEventType = "test"
)
