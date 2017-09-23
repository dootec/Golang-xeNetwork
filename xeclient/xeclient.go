package xeclient

import (
	"net"
	"bufio"
	"strings"
)

var blinkmess func(serverID string, receiver *Client, senderID string, mess string)

func BlinkMess(callback func(serverID string, receiver *Client, senderID string, mess string)) {
	blinkmess = callback
}

var blinkconn func(c *Client, b bool)

func BlinkConn(callback func(c *Client, b bool)) {
	blinkconn = callback
}

const (
	userERR        = "userERR"
	userOK         = "userOK"
	userGetAlluser = "userGetAlluser"
)

type Client struct {
	ID          string
	IP          string
	Port        string
	Conn        net.Conn
	ConnSitu    bool
	Users       []string
	Temp        string
	TempChanged bool
}

func (c *Client) CloseXE() {
	c.ConnSitu = false
	c.Users = nil
	c.Temp = ""
	c.TempChanged = false
	c.Conn.Close()
	blinkconn(c, false)
}

func (c Client) SendMessage(toID string, mess string) {
	/*
	#Life Cycle of SendMessage protocol

	Clients send message with Protocol1 to Server
	[Protocol1] = [1] @ [2]
	1: Receiver name
	2: Message

	--------------------------------------------------
	Server must get message with Protocol1 from Client
	--------------------------------------------------

	Server send message with Protocol2 to another Client
	[Protocol2] = [1] @ [2] @ [3]
	1: Server Name
	2: Sender Name
	3: Message

	--------------------------------------------------
	Client must get message with Protocol2 from Server
	--------------------------------------------------
	*/
	writer := bufio.NewWriter(c.Conn)
	writer.WriteString(toID + "@" + mess + "\n")
	go writer.Flush()
}

func StartXE(id string, ip string, port string) *Client {
	tcp, _ := net.ResolveTCPAddr("tcp", ip+":"+port)
	conn, _ := net.DialTCP("tcp", nil, tcp)
	if sendGetDirect(conn, "@"+id) == userOK {
		user := &Client{
			ID:       id,
			IP:       ip,
			Port:     port,
			Conn:     conn,
			ConnSitu: true,
		}
		go user.listen()
		return user
	} else {
		go conn.Close()
		return nil
	}
}

func sendGetDirect(conn net.Conn, s string) string {
	writer := bufio.NewWriter(conn)
	writer.WriteString(s + "\n")
	go writer.Flush()
	reader := bufio.NewReader(conn)
	for {
		temp, _ := reader.ReadString('\n')
		if (temp != "") {
			temp = strings.Replace(temp, "\n", "", -1)
			temp = strings.Replace(temp, "\r", "", -1)
			return temp
		}
	}
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.Conn)
	for {
		temp, err := reader.ReadString('\n')
		if err != nil {
			c.CloseXE()
			break
		}
		if temp != "" {
			temp = strings.Replace(temp, "\n", "", -1)
			temp = strings.Replace(temp, "\r", "", -1)
			go protocol(c, temp)
		}
	}
}

func protocol(receiver *Client, s string) {
	if (strings.Count(s, "@") == 2) {
		temp := strings.Split(s, "@")
		message := strings.Replace(s, temp[0]+"@"+temp[1]+"@", "", -1)
		go blinkmess(temp[0], receiver, temp[1], message)
		/*
		serverID := temp[0]
		senderID := temp[1]
		message := strings.Replace(s, serverID+"@"+senderID+"@", "", -1)
		go blinkmess(serverID, receiver, senderID, message)
		*/
	} else {
		receiver.Temp = s
		receiver.TempChanged = true
	}
}

func (c *Client) UpdateUsers() {
	temp := sendGetwTemp(c, userGetAlluser)
	if (strings.Contains(temp, userGetAlluser)) {
		tempsplit := strings.Split(strings.Replace(temp, userGetAlluser, "", -1), ",")
		c.Users = nil
		for _, item := range tempsplit {
			c.Users = append(c.Users, item)
		}
	}
}

func sendGetwTemp(c *Client, s string) string {
	writer := bufio.NewWriter(c.Conn)
	writer.WriteString(s + "\n")
	go writer.Flush()
	for {
		if c.TempChanged {
			c.TempChanged = false
			return c.Temp
		}
	}
}
