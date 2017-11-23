package offlineAssignment

import (
	"apt-web-server_v2/modules/mlog"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func SSHCmd(sshUser, sshPasswd, sshIP, sshCmd string, sshPort int) (err error) {
	defer func() {
		if info := recover(); info != nil {
			mlog.Debug("ssh error:", sshCmd)
			err = errors.New("ssh error")
		}
	}()
	session, err := SSHConnect(sshUser, sshPasswd, sshIP, sshPort)
	if err != nil {
		mlog.Debug("connect error:", err)
		return err
		panic(err)
	}
	defer session.Close()
	err = session.Run(sshCmd)
	if err != nil {
		mlog.Debug("run sshCmd error:", err)
		return err
		panic(err)
	}
	return nil
}
