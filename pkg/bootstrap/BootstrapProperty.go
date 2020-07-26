package bootstrap

var (
	//HTTPConfig 配置变量
	HTTPConfig         HTTPConf
	//DiscoverConfig 配置变量
	DiscoverConfig     DiscoverConf
	//ConfigServerConfig 配置变量
	ConfigServerConfig ConfigServerConf
	//RPCConfig 配置变量
	RPCConfig          RPCConf
)

//HTTPConf 配置
type HTTPConf struct {
	Host string
	Port string
}

//RPCConf 配置
type RPCConf struct {
	Port string
}

//DiscoverConf 服务发现与注册配置
type DiscoverConf struct {
	Host        string
	Port        string
	ServiceName string
	Weight      int
	InstanceID  string
}

//ConfigServerConf 配置中心
type ConfigServerConf struct {
	ID        string   // 应用名
	Profile   string   // 环境名
	Label     string   // 分支
}
