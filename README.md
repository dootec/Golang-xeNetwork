# xeNetwork
xeNetwork is a simple, handy and private network that can access any your multiple clients that via internet with TCP protocol. All you need to do set special identity of your client and send message. Don't worry, Background will be configured by xeNetwork.

The project can be used communication local or remote networks. If you want to access the remote server, modem settings of the server should be set to accept connections.


# Installation
```
go get -u github.com/dootec/xeNetwork
```


# Client Example

```
package main

import (
	"github.com/dootec/xeNetwork/xeclient"
	"fmt"
	"strconv"
)

func main() {

	user1 := xeclient.StartXE("Ege", "127.0.0.1", "6677")
	user2 := xeclient.StartXE("Golang", "127.0.0.1", "6677")
	userX := xeclient.StartXE("Foo", "127.0.0.1", "8001")
	userY := xeclient.StartXE("Admin", "127.0.0.1", "8001")

	go xeclient.BlinkMess(func(serverID string, receiver *xeclient.Client, senderID string, message string) {
		go fmt.Println("[NewM] Heyy '" + receiver.ID + "' has your new message from '" + senderID + "' and its message : '" + message + "'\t[It was automatically sent from the '" + serverID + "' server.]")
	})

	go xeclient.BlinkConn(func(c *xeclient.Client, b bool) {
		go fmt.Println("[ConC] " + c.ID + "'s connection was " + strconv.FormatBool(b))
	})

	user1.SendMessage(user2.ID, "Hi Golang")
	user2.SendMessage(user1.ID, "Hi Ege")
	userX.SendMessage(userY.ID, "We ¦ Open Source")
	userY.SendMessage(userX.ID, "Yes, me too")
	fmt.Println("Messages have sended!")
	fmt.Scanln()

	userX.UpdateUsers()
	fmt.Println(userX.Users)
	fmt.Println("Users id of userX has Updated")
	fmt.Scanln()

	user1.CloseXE()
	user2.CloseXE()
	userX.CloseXE()
	userY.CloseXE()

	fmt.Println("Connections have closed!")
	fmt.Scanln()
}
```


# Server Example

```
package main

import (
	"fmt"
	"github.com/dootec/xeNetwork/xeserver"
)

func main() {

	main := xeserver.StartXE("Main", "127.0.0.1", "6677")
	if main.ConnnSitu == true {
		fmt.Println(main.ID + " Server is working")
	}
	other := xeserver.StartXE("other", "127.0.0.1", "8001")
	if other.ConnnSitu == true {
		fmt.Println(other.ID + " Server is working")
	}
	defer main.Close()
	defer other.Close()

	go xeserver.NewUser(func(c *xeserver.Client) {
		go fmt.Println("[User] A new user named '" + c.ID + "' was joined the '" + c.Server.ID + "' server.")
	})

	go xeserver.UserClosed(func(c *xeserver.Client) {
		go fmt.Println("[Exit] A new user named '" + c.ID + "' was exited the '" + c.Server.ID + "' server.")
	})

	go xeserver.AutoPassMess(func(sender *xeserver.Client, reciever *xeserver.Client, message string) {
		go fmt.Println("[Auto] '" + sender.ID + "' is being sending a new message to '" + reciever.ID + "' and its message is : '" + message + "'")
		if message == "We ¦ Open Source" {
			sender.SendOK = false
		}
	})

	fmt.Scanln()
}
```