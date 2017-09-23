package xeserver

import (
	"net"
	"bufio"
	"strings"
)

const (
	userERR        = "userERR"
	userOK         = "userOK"
	userGetAlluser = "userGetAlluser"
)

type Client struct {
	ID     string
	Conn   net.Conn
	Server *server
	SendOK bool
}

type server struct {
	ID        string
	IP        string
	PORT      string
	ConnnSitu bool
	Users     []string
	Usersconn map[string]*Client
}

var newuser func(c *Client)

func NewUser(callback func(c *Client)) {
	newuser = callback
}

var userclosed func(c *Client)

func UserClosed(callback func(c *Client)) {
	userclosed = callback
}

var autopassmess func(sender *Client, reciever *Client, message string)

func AutoPassMess(callback func(sender *Client, reciever *Client, message string)) {
	autopassmess = callback
}

func (s server) Close() {
	for _, item := range s.Usersconn {
		item.Conn.Close()
	}
	for _, item := range s.Users {
		delete(s.Usersconn, item)
	}
	s.Close()
	s.ConnnSitu = false
}

func (c *Client) close() {
	delete(c.Server.Usersconn, c.ID)
	c.Server.Users = remove(c.Server.Users, c.ID)
	c.Server.ConnnSitu = false
	c.Conn.Close()
	userclosed(c)
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.Conn)
	for {
		temp, err := reader.ReadString('\n')
		if (err != nil) {
			break;
		}
		if (temp != "") {
			temp = strings.Replace(temp, "\n", "", -1)
			temp = strings.Replace(temp, "\r", "", -1)
			if (strings.ContainsAny(temp, "@")) {
				go redirectMess(c, c.ID, temp)
			} else if (temp == userGetAlluser) {
				go sendAlluser(c)
			}
		}
	}
	c.close()
}

func sendAlluser(c *Client) {
	var temp string = userGetAlluser
	for _, item := range c.Server.Users {
		temp += item + ","
	}
	go sendDirectMess(c.Conn, temp)
}

func redirectMess(c *Client, senderID string, s string) {
	recid := strings.Split(s, "@")[0]
	mess := strings.Replace(s, recid+"@", "", -1)
	total := c.Server.ID + "@" + senderID + "@" + mess
	cli, okk := c.Server.Usersconn[recid]
	if okk {
		autopassmess(c, cli, mess)
		if c.SendOK == true {
			go sendDirectMess(cli.Conn, total)
		} else {
			c.SendOK = true
		}
	}
}

func StartXE(serverID string, ip string, port string) *server {
	s, err := net.Listen("tcp", ip+":"+port)
	if (err != nil) {
		server := new(server)
		server.Usersconn = make(map[string]*Client)
		server.ID = serverID
		server.IP = ip
		server.PORT = port
		server.ConnnSitu = false
		return server
	} else {
		server := new(server)
		server.Usersconn = make(map[string]*Client)
		server.ID = serverID
		server.IP = ip
		server.PORT = port
		server.ConnnSitu = true
		go serverStarting(s, server)
		return server
	}
}

func serverStarting(server net.Listener, serv *server) {
	defer server.Close()
	for {
		conn, _ := server.Accept()
		go newClient(serv, conn)
	}
}

func newClient(serv *server, conn net.Conn) {
	specialID := getSpecialID(conn)
	_, okk := serv.Usersconn[specialID]
	if okk {
		sendDirectMess(conn, userERR)
		conn.Close()
	} else {
		client := &Client{
			ID:     specialID,
			Conn:   conn,
			Server: serv,
			SendOK: true,
		}
		serv.Usersconn[specialID] = client
		serv.Users = append(serv.Users, specialID)

		go client.listen()
		go newuser(client)
		go sendDirectMess(conn, userOK)
	}
}

func sendDirectMess(conn net.Conn, s string) {
	writer := bufio.NewWriter(conn)
	writer.WriteString(s + "\n")
	go writer.Flush()
}

func getSpecialID(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	for {
		temp, _ := reader.ReadString('\n')
		if (temp != "") {
			temp = strings.Replace(temp, "\n", "", -1)
			if strings.ContainsAny(temp, "@") {
				temp = strings.Replace(temp, "@", "", -1)
				temp = strings.Replace(temp, "\r", "", -1)
				return temp
			}
		}
	}
}
