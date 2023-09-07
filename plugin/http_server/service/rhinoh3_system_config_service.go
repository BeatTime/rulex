// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"
)

//--------------------------------------------------------------------------------------
// 注意: 这些设置主要是针对RhinoH3 Ubuntu16.04 的，有可能在不同的发行版有不同的指令，不一定通用
// ！！！！ Warning: MUST RUN WITH SUDO or ROOT USER  ！！！！
//--------------------------------------------------------------------------------------
/*
*
* 专门针对H3的一些系统指令封装
*
 */
func GetWlanList() ([]string, error) {
	// 执行 nmcli 命令来获取WIFI列表
	cmd := exec.Command("nmcli", "--fields", "SSID,MODE,FREQ,SIGNAL,BARS,SECURITY", "device", "wifi", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error executing nmcli: %v", err)
	}
	lines := strings.Split(string(output), "\n")
	var wifiList []string
	wifiList = append(wifiList, lines...)
	return wifiList, nil
}

/*
*
* 关闭WIFI开关
*
 */
func DisableWifi() error {
	cmd := exec.Command("nmcli", "radio", "wifi", "off")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error disabling Wi-Fi: %v", err)
	}
	return nil
}

/*
*
* 打开WIFI开关
*
 */
func EnableWifi() error {
	cmd := exec.Command("nmcli", "radio", "wifi", "on")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error enabling Wi-Fi: %v", err)
	}
	return nil
}

/*
*
* 验证时间格式 YYYY-MM-DD HH:MM:SS
*
 */
func isValidTimeFormat(input string) bool {
	expectedFormat := "2006-01-02 15:04:05"
	_, err := time.Parse(expectedFormat, input)
	return err == nil
}

/*
*
* 获取当前系统时间
*
 */
func GetSystemTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*
*
*
设置时间，格式为 "YYYY-MM-DD HH:MM:SS"
*
*/
func SetSystemTime(newTime string) error {
	if !isValidTimeFormat(newTime) {
		return fmt.Errorf("Invalid time format:%s, must be 'YYYY-MM-DD HH:MM:SS'", newTime)
	}
	// newTime := "2023-08-10 15:30:00"
	cmd := exec.Command("date", "-s", newTime)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

/*
* amixer 设置音量, 输入参数是个数值, 每次增加或者减少1%
*        amixer set 'Line Out' 1 | grep 'Front Left:' | awk -F '[][]' '{print $2}'
*
 */
func SetVolume(v int) (string, error) {
	shellCmd := "amixer set 'Line Out' %s | grep 'Front Left:' | awk -F '[][]' '{print $2}'"
	if v > -100 && v < 100 {
		var cmd *exec.Cmd
		if v < 0 {
			cmd = exec.Command("sh", "-c", fmt.Sprintf(shellCmd, fmt.Sprintf("%v%%-", math.Abs(float64(v)))))
		}
		if v > 0 {
			cmd = exec.Command("sh", "-c", fmt.Sprintf(shellCmd, fmt.Sprintf("%v%%+", math.Abs(float64(v)))))
		}
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		volume := strings.TrimSpace(string(output))
		return volume, nil
	}
	return "", fmt.Errorf("Invalid volume:%v, must be in range [0,100]", v)

}

/*
*
  - 获取音量百分比 20%
    amixer get Master | grep 'Front Left:' | awk -F '[][]' '{print $2}'

*
*/
func GetVolume() (string, error) {
	// 创建一个 Command 对象，执行多个命令通过管道连接
	cmd := exec.Command("sh", "-c",
		"amixer get 'Line Out' | grep 'Front Left:' | awk -F '[][]' '{print $2}'")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	volume := strings.TrimSpace(string(output))
	return volume, nil
}

/*
*
* 时区
*
 */
func GetTimeZone() string {
	z, _ := time.Now().Zone()
	return z

}

// SetTimeZone 设置系统时区
// timezone := "Asia/Shanghai"
func SetTimeZone(timezone string) error {
	cmd := exec.Command("timedatectl", "set-timezone", timezone)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error: %v\nOutput: %s", err, string(output))
	}
	return nil
}

/*
*
* Ubuntu16.04 刷新DNS，
*
 */
func ReloadDNS16() error {
	// 执行 /etc/init.d/networking 命令来重新加载DNS缓存
	cmd := exec.Command("/etc/init.d/networking", "restart")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error reloading DNS: %v", err)
	}

	return nil
}

/*
*
* Ubuntu18+ 刷新DNS，
*
 */
func ReloadDNS18xx() error {
	// 执行 systemd-resolved 命令来重新加载DNS缓存
	cmd := exec.Command("systemctl", "reload", "systemd-resolved.service")

	// 执行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error reloading DNS: %v", err)
	}

	return nil
}

/*
*
* 设置默认网卡
*
 */
func SetDefaultRoute(newGatewayIP string) error {
	cmd := exec.Command("ip", "route", "add", "default", "via", newGatewayIP)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

/*
*
rer@revb-h3:~$ nmcli device status

	DEVICE           TYPE      STATE         CONNECTION
	usb0             ethernet  connected     Wired connection 1
	wlx0cc6551c5026  wifi      connected     AABBCC
	eth1             ethernet  connected     eth1
	eth0             ethernet  disconnected  --
	lo               loopback  unmanaged     --

*
*/
type DeviceStatus struct {
	DEVICE     string `json:"device"`
	TYPE       string `json:"type"`
	STATE      string `json:"state"`
	CONNECTION string `json:"connection"`
}

func GetCurrentNetConnection() ([]DeviceStatus, error) {
	nmcliCmd := exec.Command("nmcli", "device", "status")
	nmcliOutput, err := nmcliCmd.Output()
	if err != nil {
		return nil, err
	}
	nmcliOutputStr := string(nmcliOutput)
	deviceStatuses := parseNmcliOutput(nmcliOutputStr)
	return deviceStatuses, nil
}
func parseNmcliOutput(output string) []DeviceStatus {
	var deviceStatuses []DeviceStatus

	// 按行分割输出
	lines := strings.Split(output, "\n")

	// 如果没有输出行，返回空切片
	if len(lines) == 0 {
		return deviceStatuses
	}

	// 获取列名
	headers := strings.Fields(lines[0])

	// 遍历剩余的行，每行是一个设备状态
	for _, line := range lines[1:] {
		fields := strings.Fields(line)

		// 如果列数不匹配列名数，跳过该行
		if len(fields) != len(headers) {
			continue
		}

		// 创建一个新的设备状态结构体，并填充数据
		var status DeviceStatus
		for i, header := range headers {
			switch header {
			case "DEVICE":
				status.DEVICE = fields[i]
			case "TYPE":
				status.TYPE = fields[i]
			case "STATE":
				status.STATE = fields[i]
			case "CONNECTION":
				status.CONNECTION = fields[i]
			}
		}

		// 将设备状态添加到切片
		deviceStatuses = append(deviceStatuses, status)
	}

	return deviceStatuses
}
