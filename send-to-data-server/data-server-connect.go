// send data to a remote server using sftp
package main

import (
	"fmt"
	"log"
	//"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

type serverParams struct  {
	server string
	login string
	passwd string
}


// copy a localfile to remote sftp server
// TODO: replace localFile by a []string
func Send_data_to_ftp_server(localFile string, dstPath string,  param serverParams) error {
	config := &ssh.ClientConfig{
		User: param.login,
		Auth: []ssh.AuthMethod{
			ssh.Password(param.passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", param.server, config)
	if err != nil {
		return fmt.Errorf("ssh.dial %s", err)
	}
	defer connection.Close()

	sftp, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatalf("sftp.newclient %s",err)
	}
	defer sftp.Close()


	// Open the source file
	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Println("open src ", err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := sftp.Create(dstPath)
	if err != nil {
		fmt.Println("sftp.Create ", err)
	}
	defer dstFile.Close()

	// write to file
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		fmt.Println("transf")
		fmt.Println(err)
	}
	return nil
}


func main() {
	var params = serverParams{"localhost:22", "xxxxxx", "yyyyyy"}
	
	if err := Send_data_to_ftp_server("datafile.txt", "dstdatafile.txt", params); err != nil {
		fmt.Println(err)
	}
}
