package device

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	serial "github.com/wwhai/goserial"
)

// 传输形式：
// `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
// const rawtcp string = "rawtcp"
// const rawudp string = "rawudp"
// const rs485rawserial string = "rs485rawserial"
// const rs485rawtcp string = "rs485rawtcp"

type _CommonConfig struct {
	Frequency   int    `json:"frequency" validate:"required"`
	AutoRequest bool   `json:"autoRequest" validate:"required"`
	Transport   string `json:"transport" validate:"required"`
	WaitTime    int    `json:"waitTime" validate:"required"`
	Timeout     int    `json:"timeout" validate:"required"`
}
type _UartConfig struct {
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}
type _ProtocolArg struct {
	In  string `json:"in" validate:"required"` // 十六进制字符串
	Out string `json:"out"`                    // 十六进制字符串
}
type _Protocol struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	RW          int    `json:"rw" validate:"required"`         // 1:RO 2:WO 3:RW
	BufferSize  int    `json:"bufferSize" validate:"required"` // 缓冲区大小
	Timeout     int    `json:"timeout" validate:"required"`    // 指令的等待时间, 在 Timeout 范围读 BufferSize 个字节, 否则就直接失败
	//---------------------------------------------------------------------
	// 下面都是校验算法相关配置:
	// -- 例如对[Byte1,Byte2,Byte3,Byte4,Byte5,Byte6,Byte7]用XOR算法比对
	//    从第一个开始，第五个结束[Byte1,Byte2,Byte3,Byte4,Byte5], 比对值位置在第六个[Byte6]
	// 伪代码：XOR(Byte[ChecksumBegin:ChecksumEnd]) == Byte[ChecksumValuePos]
	//---------------------------------------------------------------------
	Checksum         string // 校验算法，目前暂时支持: CRC16, XOR
	ChecksumValuePos string // 校验值比对位
	ChecksumBegin    uint   // 校验算法起始位置
	ChecksumEnd      uint   // 校验算法结束位置
	AutoRequest      bool   // 是否开启轮询
	AutoRequestGap   uint   // 轮询间隔
	//---------------------------------------------------------------------
	ProtocolArg _ProtocolArg `json:"protocol" validate:"required"` // 参数
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   _UartConfig          `json:"uartConfig" validate:"required"`
	DeviceConfig map[string]_Protocol `json:"deviceConfig" validate:"required"`
}
type CustomProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	serialPort serial.Port // 现阶段暂时支持串口
	// tcpConn    *net.TCPConn // rawtcp 以后支持
	// udpConn    *net.UDPConn // rawudp 以后支持
	mainConfig _CustomProtocolConfig
	locker     sync.Locker
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.locker = &sync.Mutex{}
	mdev.mainConfig = _CustomProtocolConfig{
		CommonConfig: _CommonConfig{},
		UartConfig:   _UartConfig{},
		DeviceConfig: map[string]_Protocol{},
	}
	mdev.status = typex.DEV_DOWN
	mdev.errorCount = 0
	return mdev

}

//  初始化
func (mdev *CustomProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !contains([]string{"N", "E", "O"}, mdev.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	if !contains([]string{`rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`},
		mdev.mainConfig.CommonConfig.Transport) {
		return errors.New("parity value only one of 'rawtcp','rawudp','rs485rawserial','rs485rawserial'")
	}
	// parse hex format
	for _, v := range mdev.mainConfig.DeviceConfig {
		if _, err := hex.DecodeString(v.ProtocolArg.In); err != nil {
			errMsg := fmt.Sprintf("hex.DecodeString(ProtocolArg.In) failed:%s", v.ProtocolArg.In)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}
		if v.ProtocolArg.Out != "" {
			if _, err := hex.DecodeString(v.ProtocolArg.Out); err != nil {
				errMsg := fmt.Sprintf("hex.DecodeString(ProtocolArg.Out) failed:%s", v.ProtocolArg.Out)
				glogger.GLogger.Error(errMsg)
				return fmt.Errorf(errMsg)
			}
		}

	}
	return nil
}

// 启动
func (mdev *CustomProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if mdev.mainConfig.CommonConfig.Transport == "rs485rawserial" {
		config := serial.Config{
			Address:  mdev.mainConfig.UartConfig.Uart,
			BaudRate: mdev.mainConfig.UartConfig.BaudRate,
			DataBits: mdev.mainConfig.UartConfig.DataBits,
			Parity:   mdev.mainConfig.UartConfig.Parity,
			StopBits: mdev.mainConfig.UartConfig.StopBits,
			Timeout:  time.Duration(mdev.mainConfig.CommonConfig.Timeout) * time.Second,
		}
		serialPort, err := serial.Open(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		for _, pp := range mdev.mainConfig.DeviceConfig {
			if pp.AutoRequest {
				npp := pp
				go func(ctx context.Context, npp _Protocol) {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							{
							}
						}
						//
						hexs, err0 := hex.DecodeString(npp.ProtocolArg.In)
						if err0 != nil {
							glogger.GLogger.Error(err0)
							mdev.errorCount++
							continue
						}
						// 将指令写进去
						mdev.locker.Lock()
						if core.GlobalConfig.AppDebugMode {
							log.Println("[AppDebugMode] Write data:", hexs)
						}
						if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
							glogger.GLogger.Error(err1)
							mdev.errorCount++
							continue
						}
						result := [1024]byte{} // 全局buf, 默认是1kb, 应该能覆盖绝大多数报文了
						// ctx, cancel := context.WithTimeout(typex.GCTX, time.Duration(npp.Timeout)*time.Microsecond)
						for {
							select {
							case <-ctx.Done():
								{
									glogger.GLogger.Error("read timeout")
									goto L1
								}
							default:
								// 协议的超时时间, 比如60ms
								time.Sleep(time.Duration(mdev.mainConfig.CommonConfig.WaitTime) * time.Microsecond)
								pos := 0
								for i := 0; i < npp.BufferSize; i++ {
									n, err2 := mdev.serialPort.Read(result[pos : pos+1])
									if err2 != nil {
										glogger.GLogger.Error(n, err2)
										mdev.errorCount++
										pos = i // 当某一次读取失败了会尝试重新读取, 尝试次数为 errorCount上限
										i = pos
										continue
									}
									pos++
								}
								goto L0
							}
						L0:
							if core.GlobalConfig.AppDebugMode {
								log.Println("[AppDebugMode] Read data:", result[:npp.BufferSize])
							}
							if npp.Checksum == "CRC16" {
								// 检查字节
								// check-crc()
								glogger.GLogger.Debug("启用了CRC16校验法, 但是暂时没有实现，这里默认校验完成")
							}
							if npp.Checksum == "XOR" {
								// 检查字节
								// check-xor()
								glogger.GLogger.Debug("启用了XOR校验法, 但是暂时没有实现，这里默认校验完成")
							}
							// 返回给lua参数是十六进制大写字符串
							mdev.RuleEngine.WorkDevice(mdev.Details(),
								hex.EncodeToString(result[:npp.BufferSize]))
							goto L1
						}
					L1:
						//-----------
						mdev.locker.Unlock()
						time.Sleep(time.Duration(npp.AutoRequestGap) * time.Microsecond)
					}
				}(mdev.Ctx, npp)
			}
		}
		mdev.serialPort = serialPort
		mdev.status = typex.DEV_UP
		return nil
	}

	return fmt.Errorf("unsupported transport:%s", mdev.mainConfig.CommonConfig.Transport)
}

// 从设备里面读数据出来
func (mdev *CustomProtocolDevice) OnRead(cmd int, data []byte) (int, error) {
	pp, ok := mdev.mainConfig.DeviceConfig[fmt.Sprintf("%d", cmd)]
	if ok {
		hexs, err0 := hex.DecodeString(pp.ProtocolArg.In)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			mdev.errorCount++
			return 0, err0
		}
		mdev.locker.Lock()
		// Send
		if core.GlobalConfig.AppDebugMode {
			log.Println("[AppDebugMode] Write data:", hexs)
		}
		if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
			glogger.GLogger.Error(err1)
			mdev.errorCount++
			return 0, err1
		}
		result := [50]byte{}
		ctx, _ := context.WithTimeout(typex.GCTX, time.Duration(pp.Timeout)*time.Microsecond)
		for {
			select {
			case <-ctx.Done():
				{
					return 0, errors.New("read timeout")
				}
			default:
				pos := 0
				for i := 0; i < pp.BufferSize; i++ {
					n, err2 := mdev.serialPort.Read(result[pos : pos+1])
					if err2 != nil {
						glogger.GLogger.Error(n, err2)
						mdev.errorCount++
						pos = i
						i = pos
						continue
					}
					pos++
				}
				goto RETURN
			}
		}
	RETURN:
		mdev.locker.Unlock()
		if core.GlobalConfig.AppDebugMode {
			log.Println("[AppDebugMode] Read data:", result[:pp.BufferSize])
		}
		// 返回结果
		copy(data, result[:pp.BufferSize])
		return pp.BufferSize, nil
	}
	return 0, errors.New("unknown read command")

}

// 把数据写入设备
// 根据第二个参数来找配置进去的自定义协议, 必须进来一个可识别的指令
// 其中cmd常为0,为无意义参数
func (mdev *CustomProtocolDevice) OnWrite(_ int, data []byte) (int, error) {
	pp, ok := mdev.mainConfig.DeviceConfig[string(data)]
	if ok {
		hexs, err0 := hex.DecodeString(pp.ProtocolArg.In)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			mdev.errorCount++
			return 0, err0
		}
		mdev.locker.Lock()
		// Send
		if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
			glogger.GLogger.Error(err1)
			mdev.errorCount++
			return 0, err1
		}
		mdev.locker.Unlock()
		return 0, nil
	}
	return 0, errors.New("unknown write command")
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if mdev.errorCount >= 5 {
		mdev.status = typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *CustomProtocolDevice) Stop() {
	mdev.status = typex.DEV_STOP
	mdev.CancelCTX()

}

// 设备属性，是一系列属性描述
func (mdev *CustomProtocolDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *CustomProtocolDevice) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *CustomProtocolDevice) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *CustomProtocolDevice) Driver() typex.XExternalDriver {
	return nil
}
func (mdev *CustomProtocolDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}