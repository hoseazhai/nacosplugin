package nacosplugin

import (
	"fmt"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/openlog"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net/url"
)

const (
	nacosRegistry = "nacos"
)

type Registrator struct {
	Name           string
	registryClient naming_client.INamingClient
	opts           *registry.Options
}

func genSid(appID string, serviceName string, version string, env string) string {
	//sid := strings.Join([]string{appID, serviceName, version, env}, "-")
	sid := serviceName
	return sid
}

// RegisterService register service
func (r *Registrator) RegisterService(ms *registry.MicroService) (string, error) {
	sid := genSid(ms.AppID, ms.ServiceName, ms.Version, ms.Environment)
	return sid, nil
}

// RegisterServiceInstance register service instance
func (r *Registrator) RegisterServiceInstance(sid string, instance *registry.MicroServiceInstance) (string, error) {
	//api, _ := r.registryClient.Api()
	vo := msiToInstVo(sid, instance)

	//instanceId, err := api.RegisterInstanceWithVo(vo)
	//if err != nil {
	//	openlog.GetLogger().Error("register instance failed.")
	//	return "", err
	//}
	//naming_client.INamingClient.
	success, err := r.registryClient.RegisterInstance(vo)
	if err != nil || success {
		openlog.GetLogger().Error("register instance failed.")
		return "", err
	}
	openlog.Debug(fmt.Sprintf("register instance sucess : serviceName:%s  ip:%s.",
		vo.ServiceName, vo.Ip))
	return vo.ServiceName, nil
}

func msiToInstVo(sid string, instance *registry.MicroServiceInstance) vo.RegisterInstanceParam {
	eps := registry.GetProtocolList(instance.EndpointsMap)
	openlog.GetLogger().Debug(fmt.Sprintf("epsmap = %v, eps=%v", instance.EndpointsMap, eps))
	u, _ := url.Parse(eps[0])
	return vo.RegisterInstanceParam{
		Ip:          u.Hostname(),
		ServiceName: sid,
		GroupName:   "jecloud",
		ClusterName: "default",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	}
}

// RegisterServiceAndInstance register service and instance
func (r *Registrator) RegisterServiceAndInstance(ms *registry.MicroService, instance *registry.MicroServiceInstance) (string, string, error) {
	sid := genSid(ms.AppID, ms.ServiceName, ms.Version, ms.Environment)
	vo := msiToInstVo(sid, instance)
	success, err := r.registryClient.RegisterInstance(vo)
	if err != nil || success {
		openlog.GetLogger().Error("register Service And Instance failed.")
		return "", "", err
	}
	openlog.Debug(fmt.Sprintf("register Service And Instance sucess : serviceName:%s  ip:%s sid :",
		vo.ServiceName, vo.Ip, sid))

	return sid, vo.ServiceName, nil
}

// UnRegisterMicroServiceInstance unregister micro-service instances
func (r *Registrator) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	//sid := genSid(ms.AppID, ms.ServiceName, ms.Version, ms.Environment)
	//vo := msiToInstVo(sid, instance)
	param := vo.DeregisterInstanceParam{
		GroupName:   "jecloud",
		Ip:          microServiceInstanceID,
		ServiceName: microServiceID,
		Ephemeral:   true,
	}
	success, err := r.registryClient.DeregisterInstance(param)
	if err != nil || !success {
		openlog.GetLogger().Error("unregister  Instance failed.")
		return err
	}
	openlog.Debug(fmt.Sprintf("unregister Instance sucess : :microServiceID%s  microServiceInstanceID:%s ",
		microServiceID, microServiceInstanceID))
	if err != nil {
		openlog.GetLogger().Error(fmt.Sprintf("unregister instance failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID))
		return err
	}
	return nil
}

// WSHeartbeat : Keep instance heartbeats.
func (r *Registrator) WSHeartbeat(microServiceID string, microServiceInstanceID string, callback func()) (bool, error) {
	param := vo.SelectOneHealthInstanceParam{
		ServiceName: microServiceID,
		GroupName:   "jecloud",
	}
	_, err := r.registryClient.SelectOneHealthyInstance(param)
	if err != nil {
		openlog.Error(fmt.Sprintf("Heartbeat failed, microServiceID/instanceID: %s/%s. %s",
			microServiceID, microServiceInstanceID, err))
		return false, err
	}
	openlog.Debug(fmt.Sprintf("heartbeat success, microServiceID/instanceID: %s/%s.",
		microServiceID, microServiceInstanceID))
	return true, nil
}

// Heartbeat check heartbeat of micro-service instance
func (r *Registrator) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	//r.registryClient.ServerHealthy()
	//api, _ := r.registryClient.Api()
	//err := api.SendHeartbeat(microServiceID, microServiceInstanceID)
	//if err != nil {
	//	openlog.GetLogger().Error(fmt.Sprintf("Heartbeat failed, microServiceID/instanceID: %s/%s. %s", microServiceID, microServiceInstanceID, err))
	//	return false, err
	//}
	return true, nil
}

// AddDependencies add dependencies
func (r *Registrator) AddDependencies(request *registry.MicroServiceDependency) error {
	return nil
}

// AddSchemas add schema
func (r *Registrator) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return nil
}

// UpdateMicroServiceInstanceStatus update micro-service instance status
func (r *Registrator) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	param := vo.UpdateInstanceParam{
		Ip:          microServiceInstanceID,
		ServiceName: microServiceID,
	}
	success, err := r.registryClient.UpdateInstance(param)

	if err != nil || !success {
		openlog.GetLogger().Error(fmt.Sprintf("UpdateMicroServiceInstanceStatus failed, microServiceID/instanceID = %s/%s, status=%s.", microServiceID, microServiceInstanceID, status))
		return err
	}
	openlog.GetLogger().Debug(fmt.Sprintf("UpdateMicroServiceInstanceStatus success, microServiceID/instanceID = %s/%s, status=%s.", microServiceID, microServiceInstanceID, status))
	return nil
}

// UpdateMicroServiceProperties update micro-service properities
func (r *Registrator) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return nil
}

// UpdateMicroServiceInstanceProperties update micro-service instance properities
func (r *Registrator) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	//api, _ := r.registryClient.Api()
	//err := api.UpdateMeta(microServiceID, microServiceInstanceID, properties)
	//if err != nil {
	//	openlog.GetLogger().Error(fmt.Sprintf("UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID))
	//	return err
	//}
	openlog.GetLogger().Debug(fmt.Sprintf("UpdateMicroServiceInstanceProperties success, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID))
	return nil
}

// Close close the file
func (r *Registrator) Close() error {
	return nil
}

func NewNacosRegistry(opts registry.Options) registry.Registrator {

	//serviceUrls := make([]string, 0)
	serviceConfigs := []constant.ServerConfig{}
	for _, addr := range opts.Addrs {
		//url := fmt.Sprintf("http://%s/eureka", addr)
		nacosConfig := *constant.NewServerConfig(addr, 8848, constant.WithContextPath("/nacos"))
		serviceConfigs = append(serviceConfigs, nacosConfig)
	}
	//sc := []constant.ServerConfig{
	//	*constant.NewServerConfig("192.168.20.26", 8848, constant.WithContextPath("/nacos")),
	//}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: serviceConfigs,
		},
	)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return &Registrator{
		Name:           "",
		registryClient: client,
	}
}

func Init() {
	registry.InstallRegistrator(nacosRegistry, NewNacosRegistry)
	//registry.InstallServiceDiscovery(nacosRegistry, NewEurekaDiscovery)
}
