package enums

type GlobalClientScope string

const (
	All GlobalClientScope = `*`

	ClientManagementScope GlobalClientScope = `client_management`
	UserManagementScope   GlobalClientScope = `user_management`
)
