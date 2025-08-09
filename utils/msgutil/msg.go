package msgutil

type Data map[string]interface{}
type Message struct {
	Data Data
}

func NewMessage() *Message {
	return &Message{
		Data: make(Data),
	}
}
func (m Message) set(key string, value interface{}) Message {
	m.Data[key] = value
	return m
}
func (m Message) Done() Data {
	return m.Data
}
func UserAlreadyExists() Data {
	return NewMessage().set("message", "User already exists").Done()
}
func UserCreatedSuccessfully() Data {
	return NewMessage().set("message", "User created successfully").Done()
}
func SomethingWentWrongMsg() Data {
	return NewMessage().set("message", "Something went wrong").Done()
}
func UserLoggedInSuccessfully() Data {
	return NewMessage().set("message", "Login Successful").Done()
}
func InvalidRequestMsg() Data {
	return NewMessage().set("message", "Invalid request").Done()
}
func LogoutSuccessfully() Data {
	return NewMessage().set("message", "Logout Successfully").Done()
}
func UserUnauthorized() Data { return NewMessage().set("message", "User unauthorized").Done() }
func AccessForbiddenMsg() Data {
	return NewMessage().set("message", "Access forbidden").Done()
}
func CampaignCreatedSuccessfully() Data {
	return NewMessage().set("message", "Campaign created successfully").Done()
}
func CampaignUpdatedSuccessfully() Data {
	return NewMessage().set("message", "Campaign updated successfully").Done()
}
func CampaignNotFound() Data { return NewMessage().set("message", "Campaign not found").Done() }
func CampaignDeletedSuccessfully() Data {
	return NewMessage().set("message", "Campaign deleted successfully").Done()
}
