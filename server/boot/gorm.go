package boot

import (
	"fmt"
	"github.com/flipped-aurora/gva/global"
	"github.com/gookit/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"os"
	"strings"
)

type DatabaseInfo struct {
	Value        string `gorm:"column:Value"`
	VariableName string `gorm:"column:Variable_name"`
}

var Mysql = new(_mysql)

type _mysql struct {
	db  *gorm.DB
	err error

	old       string // 配置文件第一次读取数据库数据
	input     string
	version   string
	character string
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: gorm连接mysql数据库
func (m *_mysql) Initialize() {
	_config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	if global.Config.LogMode {
		_config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		_config.Logger = logger.Default.LogMode(logger.Silent)
	}
	m.db, m.err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       global.Config.Dsn(), // DSN data source name
		DefaultStringSize:         191,                 // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,               // 根据当前 MySQL 版本自动配置
	}), _config)
	if m.err != nil {
		fmt.Println()
		fmt.Println(global.Config.Dsn())
		if err := global.Config.CreateTable(); err != nil {
			color.Debug.Printf("[%s] --> 创建数据库失败!", global.Config.DbType)
		} else {
			m.Initialize()
		}
	} else {
		global.Db = m.db
	}
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 检查数据库是否存在
func (m *_mysql) CheckDatabase() {
	var unknownDatabase = fmt.Sprintf("Unknown database '%v'", global.Config.DbName)
	if m.err != nil {
		if strings.Split(m.err.Error(), ": ")[1] == unknownDatabase {
			color.Debug.Print("\n[Mysql] -->配置文件的数据库名为:")
			color.LightGreen.Printf(" {%v} ", global.Config.DbName)
			color.Debug.Println("不存在!\n")
			color.Debug.Println("您的配置文件所配置的数据库不存在,请选择:")
			color.Debug.Print("0:请自行创建配置文件所配置的数据库名为:")
			color.LightGreen.Printf(" {%v} \n", global.Config.DbName)
			color.Debug.Print("1:尝试使用sql为您创建配置文件所配置的数据库名为:")
			color.LightGreen.Printf(" {%v} \n", global.Config.DbName)
			color.Debug.Println("2:忽略错误! 注意: 如果不修复, 将会退出初始化数据的进程!")
			color.Warn.Println("\n注意!!!!!!!")
			color.Warn.Println("输入1之后,如果配置文件的mysql用户名为root才会有百分百的权限去创建数据库,不是root的话就会跳过创建数据库步骤!\n")
			color.Debug.Println("请输入指令:")
			if n, _ := fmt.Scanln(&m.input); n != 0 {
				if m.input == "1" {
					if global.Config.Username == "root" {
						m.database()
					} else {
						color.Debug.Print("\n很抱歉,您的配置文件的mysql用户名配置不是root,不确定你有无权限创建数据库,为您跳过创建数据库操作,请自行创建配置文件所配置的数据库名为:")
						color.LightGreen.Printf(" {%v} \n", global.Config.DbName)
					}
				} else if m.input == "2" {
					os.Exit(0)
				} else {
					color.Warn.Println("[Mysql] --> 请自行创建数据库!")
					os.Exit(0)
				}
			}
		}
	}
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 检查数据库编码是不是utf8mb4
func (m *_mysql) CheckUtf8mb4() {
	if m.character != "utf8mb4" {
		color.Info.Println("您当前的数据库编码不是utf8mb4,请选择:")
		color.Info.Println("0:请自行修改数据的编码为utf8mb4!")
		color.Info.Println("1:尝试使用sql为您修改编码为utf8mb4!")
		color.Info.Println("2:忽略错误! 注意如果不修复,生成初始数据的时候也许或许可能有几率报错的喔!")
		color.Info.Println("请输入指令:")
		if n, _ := fmt.Scanln(&m.input); n != 0 {
			if m.input == "1" {
				m.utf8mb4()
			} else if m.input == "2" {
				return
			} else {
				color.Warn.Println("[Mysql] --> 请自行修改数据的编码为utf8mb4!")
				os.Exit(0)
			}
		}
	}
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 打印数据库基本信息
func (m *_mysql) Info() {

	// 获取数据库版本
	color.Debug.Println("[Mysql] -->获取数据库版本中.......")
	if err := global.Db.Raw("SELECT VERSION() AS version;").Scan(&m.version).Error; err != nil {
		color.Info.Printf("[Mysql] -->获取数据库版本失败! err: %v", err)
		m.version = "未知版本~~~"
	}

	// 获取数据库编码
	var info DatabaseInfo
	color.Debug.Println("\n[Mysql] -->获取数据库编码中.......")
	if err := global.Db.Raw("show variables like 'character_set_database' ").Scan(&info).Error; err != nil {
		color.Error.Printf("[Mysql] -->获取数据库编码失败! err:%v\n", err)
		m.character = "未知编码~~~"
	}
	color.Debug.Println("\n[Mysql] -->获取数据库编码成功!\n")
	m.character = info.Value

	color.Debug.Println("\n[Mysql] -->获取数据库版本成功!")
	color.Debug.Print("\n您当前的数据库版本: ")
	color.LightGreen.Printf(" {%v} ", m.version)
	color.Debug.Print(", 使用的数据库是: ")
	color.LightGreen.Printf(" {%v} ", global.Config.DbName)
	color.Debug.Print(", 数据库编码是: ")
	color.LightGreen.Printf(" {%v} \n\n", m.character)
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 设置配置文件的数据库编码为utf8mb4
func (m *_mysql) utf8mb4() {
	color.Debug.Print("\n[Mysql] --> 设置数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.Config.DbName)
	color.Debug.Print("数据库的编码为utf8mb4中.......\n")
	if err := global.Db.Debug().Exec("ALTER DATABASE " + global.Config.DbName + " CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci`").Error; err != nil {
		color.Debug.Print("\n[Mysql] --> 设置数据库名为:")
		color.LightGreen.Printf(" {%v} ", global.Config.DbName)
		color.Debug.Print("数据库的编码为utf8mb4失败!请手动修改数据库名为:")
		color.LightGreen.Printf(" {%v} ", global.Config.DbName)
		color.Debug.Println("的编码为utf8mb4\n")
		return
	}
	color.Info.Print("\n[Mysql] --> 设置数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.Config.DbName)
	color.Debug.Print("的编码为utf8mb4成功!\n")
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 创建配置文件的数据库
func (m *_mysql) database() {
	m.old = global.Config.DbName
	global.Config.DbName = "mysql"
	color.Debug.Printf("\n[Mysql] --> 正在连接 mysql 数据库中.......\n")
	m.Initialize()
	if m.err != nil {
		color.Error.Printf("\n[Mysql] --> 链接 mysql 数据库失败!, err: %v\n", m.err)
		color.Error.Printf("[Mysql] --> 请自行创建配置文件所需的数据库!\n")
		os.Exit(0)
	}
	color.Debug.Printf("\n[Mysql] --> 连接 mysql 数据库成功\n")

	global.Config.DbName = m.old
	color.Debug.Print("\n[Mysql] --> 正在为您创建数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.Config.DbName)
	color.Debug.Print("中.......\n")
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", global.Config.DbName)
	if m.err = global.Db.Exec(sql).Error; m.err != nil {
		color.Debug.Print("\n[Mysql] --> 创建数据库名为:")
		color.LightGreen.Printf(" {%v} ", global.Config.DbName)
		color.Debug.Print("失败!请手动修改数据库名为")
		color.LightGreen.Printf(" {%v} \n", global.Config.DbName)
		color.Error.Print("err:")
		color.LightGreen.Printf(" {%v} \n", m.err)
		os.Exit(0)
		return
	}
	color.Debug.Print("\n[Mysql] --> 正在为您创建数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.Config.DbName)
	color.Debug.Print("成功!\n")
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@description: 处理零值
func (m *_mysql) zero() {
	var info DatabaseInfo
	color.Info.Println("\n[Mysql]--> 获取数据库数据中.......")
	if err := global.Db.Raw("show variables like 'sql_mode';").Scan(&info).Error; err != nil {
		color.Error.Printf("\n[Mysql]-->获取数据库数据失败! err:%v\n", err)
	}
	color.Info.Println("\n[Mysql]--> 处理数据库返回数据.......")
	var values = strings.Split(info.Value, ",")
	info.Value = ""
	for i, value := range values {
		if value == "NO_ZERO_IN_DATE" || value == "NO_ZERO_DATE" {
		} else {
			if i == len(values)-1 {
				info.Value += value
			} else {
				info.Value += value + ","
			}
		}
	}
	if err := global.Db.Exec("set global sql_mode='" + info.Value + "';").Error; err != nil {
		color.Error.Printf("\n[Mysql]--> 设置数据库的零值失效失败! err:%v\n", err)
		return
	}
	color.Info.Println("\n[Mysql]--> 设置数据库零值失效成功")
}
