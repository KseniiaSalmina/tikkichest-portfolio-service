package sender

type Sender interface {
	Send(id int, event Event)
}

type Manager struct {
	sender Sender
}

// Object can be Portfolio, Craft or Content
type Object string

const (
	Portfolio Object = "portfolio"
	Craft     Object = "craft"
	Content   Object = "content"
)

// Change can be CreateObj, UpdateObj or DeleteObj
type Change string

const (
	CreateObj Change = "created"
	UpdateObj Change = "changed"
	DeleteObj Change = "deleted"
)

type Event struct {
	Object   Object `json:"object"`
	ObjectID int    `json:"object_id"`
	Change   Change `json:"change"`
}

func NewManager(sender Sender) *Manager {
	return &Manager{sender: sender}
}

func (n *Manager) SendEvent(userID int, obj Object, objID int, change Change) {
	event := Event{
		Object:   obj,
		ObjectID: objID,
		Change:   change,
	}
	n.sender.Send(userID, event)
}
