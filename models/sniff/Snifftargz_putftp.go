// Snifftargz_putftp.go
package sniff

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	//"os"

	ftp4go "github.com/shenshouer/ftp4go"
)

type serverInfo struct {
	Type  string `xml:"type,attr"`
	IP    string `xml:"serverIP"`
	Port  int    `xml:"serverPort"`
	User  string `xml:"serverUser"`
	Pwd   string `xml:"serverPwd"`
	Path  string `xml:"serverPath"`
	Dpath string `xml:"serverDescpath"`
}
type server struct {
	Name string     `xml:"serverName"`
	Info serverInfo `xml:"serverInfo"`
}
type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version int      `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

func put_file(path string, s_ftp server, localfile string, serfile string) bool {
	ftpClient := ftp4go.NewFTP(0) // 1 for debugging
	//connect
	//_, err := ftpClient.Connect("10.66.2.114", 21, "")
	_, err := ftpClient.Connect(s_ftp.Info.IP, s_ftp.Info.Port, "")
	if err != nil {
		fmt.Printf("The connection %s failed \n", s_ftp.Name)
		return false
	}
	defer ftpClient.Quit()
	//_, err = ftpClient.Login("test", "123456", "")
	_, err = ftpClient.Login(s_ftp.Info.User, s_ftp.Info.Pwd, "")
	if err != nil {
		fmt.Printf("The %s user or pwd failed \n", s_ftp.Name)
		return false
	}
	if path == "" {
		//_, err = ftpClient.Cwd("snap")
		_, err = ftpClient.Cwd(s_ftp.Info.Path)
		if err != nil {
			fmt.Printf("The %s path no have \n", s_ftp.Name)
			return false
		}
	} else {
		_, err = ftpClient.Cwd(path)
		if err != nil {
			fmt.Println("The path no have")
			return false
		}
	}

	//err = ftpClient.UploadFile("test.go1", "./test.go1", false, nil)
	err = ftpClient.UploadFile(serfile, localfile, false, nil)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	/*
		cwd, err := ftpClient.Dir()
		if err != nil {
			fmt.Println("The Pwd command failed")
			return false
		}
		fmt.Println("The current folder is", cwd)
	*/
	return true
}
func parser_server(S_cfg Servers, file string) Servers {
	/*
		file, err := os.Open("server.xml") // For read access.
		if err != nil {
			fmt.Printf("error: %v", err)
			return S_cfg
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
	*/
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return S_cfg
	}
	err = xml.Unmarshal(data, &S_cfg)
	if err != nil {
		fmt.Printf("error: %v", err)
		return S_cfg
	}
	return S_cfg
}
func Ftpput_xml(localfile string) bool {
	var dpicfg Servers
	var precfg Servers
	dpicfg = parser_server(dpicfg, "./conf/dpi_server.xml")
	if dpicfg.Version == 0 {
		return false
	}
	//fmt.Println(dpicfg)
	for _, s_ftp := range dpicfg.Svs {
		if s_ftp.Info.Type == "ftp" {
			//fmt.Println(i, s_ftp)
			rst := put_file("", s_ftp, localfile, "PM_Upgrade.tar.gz")
			if rst == false {
				fmt.Printf("upload %s err \n", localfile)
				//continue
				return false
			}
			rst = put_file(s_ftp.Info.Dpath, s_ftp, "./conf/PM_Upgrade_actionlist.xml", "PM_Upgrade_actionlist.xml")
			if rst == false {
				fmt.Printf("upload %s err \n", "PM_Upgrade_actionlist.xml")
				//continue
				return false
			}
			//fmt.Printf("upload %s ok \n", localfile)

		}
	}
	precfg = parser_server(precfg, "./conf/pre_server.xml")
	if precfg.Version == 0 {
		return false
	}
	//fmt.Println(precfg)
	for _, s_ftp := range precfg.Svs {
		if s_ftp.Info.Type == "ftp" {
			//fmt.Println(i, s_ftp)
			rst := put_file("", s_ftp, "./rules/rule", "rule")
			if rst == false {
				fmt.Printf("upload %s err \n", "rule")
				//continue
				return false
			}
			//fmt.Printf("upload %s ok \n", localfile)
		}
	}
	return true
}
