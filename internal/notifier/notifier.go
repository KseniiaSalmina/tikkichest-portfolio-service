package notifier

import (
	"fmt"
)

type Sender interface {
	SendNotification(id int, notification string)
}

type Notifier struct {
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

func NewNotifier(sender Sender) *Notifier {
	return &Notifier{sender: sender}
}

func (n *Notifier) Notify(userID int, obj Object, objID int, change Change) {
	notification := fmt.Sprintf("Your %s №%d has been %s", string(obj), objID, string(change))
	n.sender.SendNotification(userID, notification)
}
