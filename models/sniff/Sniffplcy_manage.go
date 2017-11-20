// Sniffplcy_manage.go
package sniff

import (
	"apt-web-server_v2/models/modelsPublic"
	"fmt"
)

func Manage() bool {
	rst := Sniff_xml_creat("./xml/sniff_plcy.xml")
	if rst == false {
		modelsPublic.Update_mysql(`update sniff_plcy_issued set status=0,err="生成策略文件失败"`)
		fmt.Println("creat sniff plcy err!!!!")
		return rst
	}
	rst = Tar_xml("./xml/PM_Upgrade.tar.gz")
	if rst == false {
		modelsPublic.Update_mysql(`update sniff_plcy_issued set status=0,err="策略文件打包失败"`)
		fmt.Println("tar err")
		return rst
	}
	rst = Ftpput_xml("./xml/PM_Upgrade.tar.gz")
	if rst == false {
		modelsPublic.Update_mysql(`update sniff_plcy_issued set status=0,err="下发策略失败"`)
		fmt.Println("ftp put err")
		return rst
	}
	return true
}
