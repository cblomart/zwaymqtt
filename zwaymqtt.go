package main

import (
  "os"
  "log"
  "flag"
  "fmt"
  "time"
  "strings"
  "regexp"
  "errors"
  "strconv"
  "encoding/json"
  "net/http"
  "io/ioutil"

  MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
  pfl "github.com/davecheney/profile"
)

type MqttUpdate struct {
  Topic string
  Value string
}

type Gateway struct {
  Key string
  Topic string
  Value string
  Write bool
  Type string
  Args []string
}

//command line variable
var zway_server string
var zway_username string
var zway_password string
var zway_home string
var zway_refresh int
var mqtt_server string
var mqtt_username string
var mqtt_password string
var mqtt_protocol string
var debug bool
var profile string

//used variables
var zway_timestamp int 
var zway_dataapi = "/ZWaveAPI/Data/"
var zway_zautoapi = "/ZAutomation/api/v1/"
var zway_runapi = "/ZWaveAPI/Run/"
var zway_cookiename = "ZWAYSession"
var http_client = new(http.Client)
var zway_cookie = new(http.Cookie)
var gateways []Gateway
var zway_retries int


//ZWay enumerations
const (
  BASIC_TYPE_CONTROLER = 1
  BASIC_TYPE_STATIC_CONTROLER = 2
  BASIC_TYPE_SLAVE = 3
  BASIC_TYPE_ROUTING_SLAVE = 4
  GENERIC_TYPE_THERMOSTAT = 8
  GENERIC_TYPE_BINARY_SWITCH = 16
  GENERIC_TYPE_MULTILEVEL_SWITCH = 17
  GENERIC_TYPE_SWITCH_REMOTE = 18
  GENERIC_TYPE_SWITCH_TOGGLE = 19
  GENERIC_TYPE_SECURITY_PANEL = 23
  GENERIC_TYPE_BINARY_SENSOR = 32
  GENERIC_TYPE_MULTILEVEL_SENSOR = 33
  GENERIC_TYPE_METER = 49
  GENERIC_TYPE_ENTRY_CONTROL = 64
  GENERIC_TYPE_ALARM_SENSOR = 161
  COMMAND_CLASS_NO_OPERATION = 0
  COMMAND_CLASS_BASIC = 32
  COMMAND_CLASS_CONTROLLER_REPLICATION = 33
  COMMAND_CLASS_APPLICATION_STATUS = 34
  COMMAND_CLASS_ZIP_SERVICES = 35
  COMMAND_CLASS_ZIP_SERVER = 36
  COMMAND_CLASS_SWITCH_BINARY = 37
  COMMAND_CLASS_SWITCH_MULTILEVEL = 38
  COMMAND_CLASS_SWITCH_ALL = 39
  COMMAND_CLASS_SWITCH_TOGGLE_BINARY = 40
  COMMAND_CLASS_SWITCH_TOGGLE_MULTILEVEL = 41
  COMMAND_CLASS_CHIMNEY_FAN = 42
  COMMAND_CLASS_SCENE_ACTIVATION = 43
  COMMAND_CLASS_SCENE_ACTUATOR_CONF = 44
  COMMAND_CLASS_SCENE_CONTROLLER_CONF = 45
  COMMAND_CLASS_ZIP_CLIENT = 46
  COMMAND_CLASS_ZIP_ADV_SERVICES = 47
  COMMAND_CLASS_SENSOR_BINARY = 48
  COMMAND_CLASS_SENSOR_MULTILEVEL = 49
  COMMAND_CLASS_METER = 50
  COMMAND_CLASS_ZIP_ADV_SERVER = 51
  COMMAND_CLASS_ZIP_ADV_CLIENT = 52
  COMMAND_CLASS_METER_PULSE = 53
  COMMAND_CLASS_THERMOSTAT_HEATING = 56
  COMMAND_CLASS_METER_TABLE_CONFIG = 60
  COMMAND_CLASS_METER_TABLE_MONITOR = 61
  COMMAND_CLASS_METER_TABLE_PUSH = 62
  COMMAND_CLASS_THERMOSTAT_MODE = 64
  COMMAND_CLASS_THERMOSTAT_OPERATING_STATE = 66
  COMMAND_CLASS_THERMOSTAT_SET_POINT = 67
  COMMAND_CLASS_THERMOSTAT_FAN_MODE = 68
  COMMAND_CLASS_THERMOSTAT_FAN_STATE = 69
  COMMAND_CLASS_CLIMATE_CONTROL_SCHEDULE = 70
  COMMAND_CLASS_THERMOSTAT_SETBACK = 71
  COMMAND_CLASS_DOOR_LOCK_LOGGING = 76
  COMMAND_CLASS_SCHEDULE_ENTRY_LOCK = 78
  COMMAND_CLASS_BASIC_WINDOW_COVERING = 80
  COMMAND_CLASS_MTP_WINDOW_COVERING = 81
  COMMAND_CLASS_SCHEDULE = 83
  COMMAND_CLASS_CRC_16_ENCAP = 86
  COMMAND_CLASS_ASSOCIATION_GROUP_INFO = 89
  COMMAND_CLASS_DEVICE_RESET_LOCALLY = 90
  COMMAND_CLASS_CENTRAL_SCENE = 91
  COMMAND_CLASS_IP_ASSOCIATION = 92
  COMMAND_CLASS_ANTITHEFT = 93
  COMMAND_CLASS_ZWAVEPLUS_INFO = 94
  COMMAND_CLASS_MULTI_INSTANCE = 96
  COMMAND_CLASS_DOOR_LOCK = 98
  COMMAND_CLASS_USER_CODE = 99
  COMMAND_CLASS_BARRIER_OPERATOR = 102
  COMMAND_CLASS_CONFIGURATION = 112
  COMMAND_CLASS_ALARM = 113
  COMMAND_CLASS_MANUFACTURER_SPECIFIC = 114
  COMMAND_CLASS_POWER_LEVEL = 115
  COMMAND_CLASS_PROTECTION = 117
  COMMAND_CLASS_LOCK = 118
  COMMAND_CLASS_NODE_NAMING = 119
  COMMAND_CLASS_FIRMWARE_UPDATE = 122
  COMMAND_CLASS_GROUPING_NAME = 123
  COMMAND_CLASS_REMOTE_ASSOCIATION_ACTIVATE = 124
  COMMAND_CLASS_REMOTE_ASSOCIATION = 125
  COMMAND_CLASS_BATTERY = 128
  COMMAND_CLASS_CLOCK = 129
  COMMAND_CLASS_HAIL = 130
  COMMAND_CLASS_WAKEUP = 132
  COMMAND_CLASS_ASSOCIATION = 133
  COMMAND_CLASS_VERSION = 134
  COMMAND_CLASS_INDICATOR = 135
  COMMAND_CLASS_PROPRIETRAY = 136
  COMMAND_CLASS_LANGUAGE = 137
  COMMAND_CLASS_TIME = 138
  COMMAND_CLASS_TIME_PARAMETERS = 139
  COMMAND_CLASS_GEOGRAPHIC_LOCATION = 140
  COMMAND_CLASS_COMPOSITE = 141
  COMMAND_CLASS_MULTICHANNEL_ASSOCIATION = 142
  COMMAND_CLASS_MULTI_CMD = 143
  COMMAND_CLASS_ENERGY_PRODUCTION = 144
  COMMAND_CLASS_MANUFACTURER_PROPRIETRATY = 145
  COMMAND_CLASS_SCREEN_MD = 146
  COMMAND_CLASS_SCREEN_ATTRIBUTES = 147
  COMMAND_CLASS_SIMPLE_AV_CONTROL = 148
  COMMAND_CLASS_AV_CONTENT_DIRECTORY_MD = 149
  COMMAND_CLASS_RENDERER_STATUS = 150
  COMMAND_CLASS_AV_CONTENT_SEARCH_MD = 151
  COMMAND_CLASS_SECURITY = 152
  COMMAND_CLASS_AV_TAGGING_MD = 153
  COMMAND_CLASS_IP_CONFIGURATION = 154
  COMMAND_CLASS_ASSOCIATION_COMMAND_CONFIGURATION = 155
  COMMAND_CLASS_ALARM_SENSOR = 156
  COMMAND_CLASS_SILENCE_ALARM = 157
  COMMAND_CLASS_SENSOR_CONFIGURATION = 158
  COMMAND_CLASS_MARK = 239
  COMMAND_CLASS_NON_INEROPERABLE = 240
)

var ZWaveClassNames = [...]string{
COMMAND_CLASS_NO_OPERATION: "command no operation",
COMMAND_CLASS_BASIC: "command basic",
COMMAND_CLASS_CONTROLLER_REPLICATION: "command controler replication",
COMMAND_CLASS_APPLICATION_STATUS: "command application status",
COMMAND_CLASS_ZIP_SERVICES: "command zip services",
COMMAND_CLASS_ZIP_SERVER: "command zip server",
COMMAND_CLASS_SWITCH_BINARY: "command switch binary",
COMMAND_CLASS_SWITCH_MULTILEVEL: "command switch multilevel",
COMMAND_CLASS_SWITCH_ALL: "commad switch all",
COMMAND_CLASS_SWITCH_TOGGLE_BINARY: "command switch toggle binary",
COMMAND_CLASS_SWITCH_TOGGLE_MULTILEVEL: "command switch toggle multilevel",
COMMAND_CLASS_CHIMNEY_FAN: "command chimney fan",
COMMAND_CLASS_SCENE_ACTIVATION: "command scene activation",
COMMAND_CLASS_SCENE_ACTUATOR_CONF: "command scene actuator configuration",
COMMAND_CLASS_SCENE_CONTROLLER_CONF: "command scene controler configuration",
COMMAND_CLASS_ZIP_CLIENT: "command zip client",
COMMAND_CLASS_ZIP_ADV_SERVICES: "command zip adv services",
COMMAND_CLASS_SENSOR_BINARY: "command sensor binary",
COMMAND_CLASS_SENSOR_MULTILEVEL: "command sensor multilevel",
COMMAND_CLASS_METER: "command meter",
COMMAND_CLASS_ZIP_ADV_SERVER: "command zip adv server",
COMMAND_CLASS_ZIP_ADV_CLIENT: "command zip adv client",
COMMAND_CLASS_METER_PULSE: "command meter pulse",
COMMAND_CLASS_THERMOSTAT_HEATING: "command thermostat heating",
COMMAND_CLASS_METER_TABLE_CONFIG: "command meter table config",
COMMAND_CLASS_METER_TABLE_MONITOR: "command meter table monitor",
COMMAND_CLASS_METER_TABLE_PUSH: "command meter table push",
COMMAND_CLASS_THERMOSTAT_MODE: "command thermostat mode",
COMMAND_CLASS_THERMOSTAT_OPERATING_STATE: "command thermostat operationg state",
COMMAND_CLASS_THERMOSTAT_SET_POINT: "command thermostat set point",
COMMAND_CLASS_THERMOSTAT_FAN_MODE: "command thermostat fan mode",
COMMAND_CLASS_THERMOSTAT_FAN_STATE: "command thermostat fan state",
COMMAND_CLASS_CLIMATE_CONTROL_SCHEDULE: "command climate control schedule",
COMMAND_CLASS_THERMOSTAT_SETBACK: "command thermostat setback",
COMMAND_CLASS_DOOR_LOCK_LOGGING: "command door lock logging",
COMMAND_CLASS_SCHEDULE_ENTRY_LOCK: "command schedule entry lock",
COMMAND_CLASS_BASIC_WINDOW_COVERING: "command basic window covering",
COMMAND_CLASS_MTP_WINDOW_COVERING: "command mtp window covering",
COMMAND_CLASS_SCHEDULE: "command shedule",
COMMAND_CLASS_CRC_16_ENCAP: "command crc 16 encap",
COMMAND_CLASS_ASSOCIATION_GROUP_INFO: "command association group info",
COMMAND_CLASS_DEVICE_RESET_LOCALLY: "command device reset locally",
COMMAND_CLASS_CENTRAL_SCENE: "command central scene",
COMMAND_CLASS_IP_ASSOCIATION: "command ip association",
COMMAND_CLASS_ANTITHEFT: "command antitheft",
COMMAND_CLASS_ZWAVEPLUS_INFO: "command zwaveplus info",
COMMAND_CLASS_MULTI_INSTANCE: "command multi instance",
COMMAND_CLASS_DOOR_LOCK: "command door lock",
COMMAND_CLASS_USER_CODE: "command user code",
COMMAND_CLASS_BARRIER_OPERATOR: "command barrier operator",
COMMAND_CLASS_CONFIGURATION: "command configuration",
COMMAND_CLASS_ALARM: "command alarm",
COMMAND_CLASS_MANUFACTURER_SPECIFIC: "commad manufacturer specific",
COMMAND_CLASS_POWER_LEVEL: "command power level",
COMMAND_CLASS_PROTECTION: "command protection",
COMMAND_CLASS_LOCK: "command lock",
COMMAND_CLASS_NODE_NAMING: "command node naming",
COMMAND_CLASS_FIRMWARE_UPDATE: "command firmware update",
COMMAND_CLASS_GROUPING_NAME: "command grouping name",
COMMAND_CLASS_REMOTE_ASSOCIATION_ACTIVATE: "command remote association activte",
COMMAND_CLASS_REMOTE_ASSOCIATION: "command remote association",
COMMAND_CLASS_BATTERY: "command battery",
COMMAND_CLASS_CLOCK: "command clock",
COMMAND_CLASS_HAIL: "command hail",
COMMAND_CLASS_WAKEUP: "command wakeup",
COMMAND_CLASS_ASSOCIATION: "command association",
COMMAND_CLASS_VERSION: "command version",
COMMAND_CLASS_INDICATOR: "command indicator",
COMMAND_CLASS_PROPRIETRAY: "command proprietary",
COMMAND_CLASS_LANGUAGE: "command language",
COMMAND_CLASS_TIME: "command time",
COMMAND_CLASS_TIME_PARAMETERS: "command time parameters",
COMMAND_CLASS_GEOGRAPHIC_LOCATION: "command geographic location",
COMMAND_CLASS_COMPOSITE: "command position",
COMMAND_CLASS_MULTICHANNEL_ASSOCIATION: "command multichannel association",
COMMAND_CLASS_MULTI_CMD: "command multi cmd",
COMMAND_CLASS_ENERGY_PRODUCTION: "command energy production",
COMMAND_CLASS_MANUFACTURER_PROPRIETRATY: "command manufacturer proprietary",
COMMAND_CLASS_SCREEN_MD: "command screen md",
COMMAND_CLASS_SCREEN_ATTRIBUTES: "command screen attributes",
COMMAND_CLASS_SIMPLE_AV_CONTROL: "command simple av control",
COMMAND_CLASS_AV_CONTENT_DIRECTORY_MD: "command av content directory",
COMMAND_CLASS_RENDERER_STATUS: "command renderer status",
COMMAND_CLASS_AV_CONTENT_SEARCH_MD: "command av content search md",
COMMAND_CLASS_SECURITY: "command security",
COMMAND_CLASS_AV_TAGGING_MD: "command av tagging md",
COMMAND_CLASS_IP_CONFIGURATION: "command ip configuration",
COMMAND_CLASS_ASSOCIATION_COMMAND_CONFIGURATION:
  "command association command configuration",
COMMAND_CLASS_ALARM_SENSOR: "command alarm sensor",
COMMAND_CLASS_SILENCE_ALARM: "command silence alarm",
COMMAND_CLASS_SENSOR_CONFIGURATION: "command sensor configuration",
COMMAND_CLASS_MARK: "command mark",
COMMAND_CLASS_NON_INEROPERABLE: "command non interoperable",
}

var ZWaveTypeNames = [...]string{
  BASIC_TYPE_CONTROLER: "basic controler",
  BASIC_TYPE_STATIC_CONTROLER: "basic static controler",
  BASIC_TYPE_SLAVE: "basic slave",
  BASIC_TYPE_ROUTING_SLAVE: "basic routing slave",
  GENERIC_TYPE_THERMOSTAT: "generic thermostat",
  GENERIC_TYPE_BINARY_SWITCH: "generic binary switch",
  GENERIC_TYPE_MULTILEVEL_SWITCH: "generic multilevel switch",
  GENERIC_TYPE_SWITCH_REMOTE: "generic switch remote",
  GENERIC_TYPE_SWITCH_TOGGLE: "generic switch toggle",
  GENERIC_TYPE_SECURITY_PANEL: "generic security panel",
  GENERIC_TYPE_BINARY_SENSOR: "generic binary sensor",
  GENERIC_TYPE_MULTILEVEL_SENSOR: "generic multilevel sensor",
  GENERIC_TYPE_METER: "generic meter",
  GENERIC_TYPE_ENTRY_CONTROL: "generic entry control",
  GENERIC_TYPE_ALARM_SENSOR: "generic alarm sensor",
}

func (g *Gateway) ToString() string {
  w := "->"
  if g.Write { w = "<>" }
  return fmt.Sprintf("%s %s %s (%s)", g.Key, w, g.Topic, g.Type)
}

func (g *Gateway) GetValue(update map[string]interface{}) string {
  switch g.Type {
  case "string":
    value, err := jsonStringValue(g.Key + "." + g.Value,update)
    if err == nil {
      return value
    }
  case "int":
    value, err := jsonFloatValue(g.Key + "." + g.Value,update)
    if err == nil {
      return fmt.Sprintf("%d", int(value))
    }
  case "float":
    value, err := jsonFloatValue(g.Key + "." + g.Value,update)
    if err == nil {
      v := fmt.Sprintf("%.3f", value)
      if strings.Contains(v,".") {
        v = strings.TrimRight(v,"0.")
      }
      return v
    }
  case "bool":
    value, err := jsonBoolValue(g.Key + "." + g.Value,update)
    if err == nil {
      return fmt.Sprintf("%t", value)
    }
  }
  return ""
}

func init() {
  //initialize command line parameters
  flag.StringVar(&zway_server,"s","localhost:8083","Z-Way server name or ZWAY_SERVER environment variable")
  flag.StringVar(&zway_username,"u","admin","Z-Way username or ZWAY_USERNAME environment variable")
  flag.StringVar(&zway_password,"p","","Z-Way passsword or ZWAY_PASSWORD environment variable")
  flag.StringVar(&zway_home,"h","razberry","mqtt topic root or ZWAY_HOME environment variable")
  flag.StringVar(&mqtt_server,"m","localhost:1883","MQTT server or MQTT_SERVER environment variable")
  flag.StringVar(&mqtt_username,"mu","","MQTT username or MQTT_USERNAME environment variable")
  flag.StringVar(&mqtt_password,"mp","","MQTT password or MQTT_PASSWORD environment variable")
  flag.StringVar(&mqtt_protocol,"proto","tcp","MQTT protocol tcp/ws/tls or MQTT_PROTOCOL environment variable")
  flag.IntVar(&zway_refresh,"r",30,"Z-Way refresh rate in seconds or ZWAY_REFRESH environment variable")
  flag.BoolVar(&debug,"v",false,"Show debug messages")
  flag.StringVar(&profile,"profile","","Profile execution (cpu/mem/all)")
  flag.Parse()
  
  // check defaults against environment variables
  if zway_server == "localhost:8083" && len(os.Getenv("ZWAY_SERVER")) > 0 {
      zway_server = os.Getenv("ZWAY_SERVER")
  }
  
  if zway_username == "admin" && len(os.Getenv("ZWAY_USERNAME")) > 0 {
      zway_username = os.Getenv("ZWAY_USERNAME")
  }
  
  if len(zway_password) == 0 && len(os.Getenv("ZWAY_PASSWORD")) > 0 {
      zway_password = os.Getenv("ZWAY_PASSWORD")
  }
  
  if zway_home == "razberry" && len(os.Getenv("ZWAY_HOME")) > 0 {
      zway_home = os.Getenv("ZWAY_HOME")
  }
  
  if zway_refresh == 30 && len(os.Getenv("ZWAY_REFRESH")) > 0 {
      zway_refresh, _ = strconv.Atoi(os.Getenv("ZWAY_REFRESH"))
  }
  
  if mqtt_server == "localhost:1883" && len(os.Getenv("MQTT_SERVER")) > 0 {
      mqtt_server = os.Getenv("MQTT_SERVER")
  }
  
  if len(mqtt_username) == 0 && len(os.Getenv("MQTT_USERNAME")) > 0 {
      mqtt_username = os.Getenv("MQTT_USERNAME")
  }
  
  if len(mqtt_password) == 0 && len(os.Getenv("MQTT_PASSWORD")) > 0 {
      mqtt_password = os.Getenv("MQTT_PASSWORD")
  }  
  
  if mqtt_protocol == "tcp" && len(os.Getenv("MQTT_PROTOCOL")) > 0 {
    mqtt_protocol = os.Getenv("MQTT_PROTOCOL")
  }
  
  if !debug && len(os.Getenv("ZWAYMQTT_DEBUG")) > 0 {
    if os.Getenv("ZWAYMQTT_DEBUG") == "true" {
      debug = true
    }
  }
  
  //standardise hostname values to <host>:<port>
  zway_match, err := regexp.MatchString(":[0-9]+$",zway_server)
  if err != nil {
    log.Fatal(fmt.Sprintf("Could not use regexp: %s", err))
  }
  if zway_match == false {
    log.Print("Setting port 8083 on given Z-Way server")
    zway_server = zway_server + ":8083"
  }
  mqtt_match, err := regexp.MatchString(":[0-9]+$",mqtt_server)
  if err != nil {
    log.Fatal(fmt.Sprintf("Could not use regexp: %s", err))
  }
  if mqtt_match == false {
    log.Print("Setting port 1883 on given MQTT server")
    mqtt_server = mqtt_server + ":1883"
  }
}

func getzway() string {
  if (debug) { log.Print("Getting Z-Way update.") }
  url := fmt.Sprintf("http://%s%s%d", zway_server, zway_dataapi, zway_timestamp)
  req, err := http.NewRequest("GET",url,nil)
  if err != nil {
    log.Printf("Error initializing request: %s", err)
  }
  if zway_cookie != nil {
    req.AddCookie(zway_cookie)
  }
  rsp, err := http_client.Do(req)
  if err != nil {
    log.Printf("Could not make zway update: %s", err)
    return ""
  }
  defer rsp.Body.Close()
  bdy, err := ioutil.ReadAll(rsp.Body)
  if err != nil {
    log.Printf("could not read body: %s", err)
  }
  return string(bdy)
}

func authzway() {
  //getting Zway authentication cookie
  url := fmt.Sprintf("http://%s%slogin", zway_server, zway_zautoapi)
  login := fmt.Sprintf("{\"login\": \"%s\", \"password\": \"%s\"}",
    zway_username, zway_password)
  req, err := http.NewRequest("POST",url,strings.NewReader(login))
  if err != nil {
    log.Printf("Error initializing request: %s", err)
  }
  req.Header.Set("Content-Type", "application/json")
  rsp, err := http_client.Do(req)
  if err != nil {
    log.Fatalf("Could not login to Z-Way: %s", err)
  }
  cookies := rsp.Cookies()
  for i := range cookies {
    if cookies[i].Name == zway_cookiename && cookies[i].Path == "/" {
      zway_cookie = cookies[i]
      break
    }
  }
  if zway_cookie == nil {
    log.Fatal("Z-Way cookie not found.")
  }
}

func jsonValue(key string, target map[string]interface{}) (interface{}, error) {
  //if the value is directly found... return it
  if target[key] != nil {
    return target[key], nil
  }
  current := target
  keys := strings.Split(key,".")
  for i := range keys[:len(keys)-1] {
    value := current[keys[i]]
    if value == nil {
      return nil, fmt.Errorf("Json Key not existent (%s)", keys[i])
    }
    current = value.(map[string]interface{})
  }
  key = keys[len(keys)-1]
  value := current[key]
  if value != nil {
    return value, nil
  }
  return nil, errors.New("Json Value non existent.")
}

func jsonStringValue(key string, target map[string]interface{}) (string, error) {
  iface, err := jsonValue(key,target)
  if err != nil {
    return "", err
  }
  return iface.(string), nil
}

func jsonIntValue(key string, target map[string]interface{}) (int, error) {
  iface, err := jsonValue(key,target)
  if err != nil {
    return 0, err
  }
  return iface.(int), nil
}

func jsonFloatValue(key string, target map[string]interface{}) (float64, error) {
  iface, err := jsonValue(key,target)
  if err != nil {
    return 0.0, err
  }
  return iface.(float64), nil
}

func jsonMapValue(key string, target map[string]interface{}) (map[string]interface{}, error) {
  iface, err := jsonValue(key,target)
  if err != nil {
    return nil, err
  }
  return iface.(map[string]interface{}), nil
}

func jsonBoolValue(key string, target map[string]interface{}) (bool, error) {
  iface, err := jsonValue(key,target)
  if err != nil {
    return false, err
  }
  return iface.(bool), nil
}

func zwaygetcmdclassdata(cmdClasses map[string]interface{}, cmdClass int) (map[string]interface{}, error) {
  iface := cmdClasses[strconv.Itoa(cmdClass)]
  if iface == nil {
    return nil, errors.New("Command class not implemented by instance")
  }
  class := iface.(map[string]interface{})
  data, err := jsonMapValue("data",class)
  if err != nil {
    return nil, err
  }
  return data, nil
}

func normName(name string) string {
  //trim
  res := strings.Trim(name," /")
  //lower
  res = strings.ToLower(res)
  //spaces
  res = strings.Replace(res," ","_",-1)
  //percents
  res = strings.Replace(res,"%","pc",-1)
  //deg
  res = strings.Replace(res,"°","",-1)
  return res
}

func zwayparsedevices(update map[string]interface{}) {
  log.Print("Parse Z-Way devices")
  for node, info := range update {
    m := info.(map[string]interface{})
    basicType, err := jsonFloatValue("data.basicType.value",m)
    if err != nil {
      log.Printf("basic type not found: %s", err)
      continue
    }
    givenName, err := jsonStringValue("data.givenName.value",m)
    if err != nil {
      log.Printf("given name not found: %s", err)
      continue
    }
    // specificType := int(jsonFloatValue("data.specificType.value",m))
    isControler := false
    switch int(basicType) {
    case BASIC_TYPE_CONTROLER:
      isControler = true
    case BASIC_TYPE_STATIC_CONTROLER:
      isControler = true
    }
    // skip if controller
    if isControler {
      log.Printf("Skipping node %s: %s", node, ZWaveTypeNames[int(basicType)])
      continue
    }
    // skip if no name
    if len(givenName) == 0 {
      log.Printf("given name empty")
      continue
    }
    // parsing instances
    instances, err := jsonMapValue("instances",m)
    if err != nil {
      continue
    }
    for i := range instances {
      // get instance
      instance := instances[i].(map[string]interface{})
      // get command classes from instance
      commandClasses, err := jsonMapValue("commandClasses",instance)
      if err != nil {
        log.Printf("command classes not found: %s", err)
        continue
      }
      // check if instance has battery informations
      data, err := zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_BATTERY)
      if err == nil {
        nkey := fmt.Sprintf("devices.%s.instances.%s.commandClasses.%d.data",
          node, i, COMMAND_CLASS_BATTERY)
        topic := fmt.Sprintf("%s/sensors/analogic/%s/%s/battery",
          zway_home, normName(givenName),i)
        _, err = jsonFloatValue("last.value",data)
        if err == nil  {
          gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
            Value: "last.value", Write:false, Type: "float"})
        }
      }
      // check if instance has switch binary
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_SWITCH_BINARY)
      if err == nil {
        nkey := fmt.Sprintf("devices.%s.instances.%s.commandClasses.%d.data",
           node, i, COMMAND_CLASS_SWITCH_BINARY)
        topic := fmt.Sprintf("%s/actuators/binary/%s/%s/switch",
          zway_home, normName(givenName), i)
        _, err = jsonBoolValue("level.value",data)
        if err == nil {
          gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
            Value: "level.value", Write:true, Type: "bool"})          
        }
      }
      // check if instance has switch multilevel
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_SWITCH_MULTILEVEL)
      if err == nil {
        nkey := fmt.Sprintf("devices.%s.instances.%s.commandClasses.%d.data",
          node, i, COMMAND_CLASS_SWITCH_MULTILEVEL)
        topic := fmt.Sprintf("%s/actuators/analogic/%s/%s/switch",
          zway_home, normName(givenName),i)
        _, err = jsonFloatValue("level.value", data)
        if err == nil {
          gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
            Value: "level.value", Write:true, Type: "float"})        
        }
      }
      // check if instance has sensor binary
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_SENSOR_BINARY)
      if err == nil {
        sensorType := "generic"
        nkey := fmt.Sprintf("devices.%s.instances.%s.commandClasses.%d.data",
          node, i, COMMAND_CLASS_SENSOR_BINARY)
        topic := fmt.Sprintf("%s/sensors/binary/%s/%s/%s",
           zway_home, normName(givenName), i, sensorType)
        _, err = jsonBoolValue("level.value",data)
        if err == nil {
            gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
               Value: "level.value", Write:false, Type: "bool"})
        }
        for k, v := range data {
          if _, err := strconv.Atoi(k); err == nil {
            sensor := v.(map[string]interface{})
            sensorType, err := jsonStringValue("sensorTypeString.value",sensor)
            if err != nil {
              log.Printf("Could not get sensor type: %s", err)
              continue
            }
            nnkey := fmt.Sprintf("%s.%s",nkey,k)
            topic := fmt.Sprintf("%s/sensors/binary/%s/%s/%s",
              zway_home,normName(givenName), i, normName(sensorType))
            _, err = jsonBoolValue("level.value",sensor)
            if err == nil {
              gateways = append(gateways, Gateway{Key: nnkey, Topic: topic,
                Value: "level.value", Write:false, Type: "bool"})
            }
          }
        }
      }
      // check if instance has sensor multilevel
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_SENSOR_MULTILEVEL)
      if err == nil {
        for k, v := range data {
          if _, err := strconv.Atoi(k); err == nil {
            sensor := v.(map[string]interface{})
            sensorType, err := jsonStringValue("sensorTypeString.value",
              sensor)
            if err != nil {
              log.Printf("Could not get sensor type: %s", err)
              continue
            }
            sensorScale, err := jsonStringValue("scaleString.value",
              sensor)
            if err != nil {
              log.Printf("Could not get sensor scale: %s", err)
              continue
            }
            nkey := fmt.Sprintf(
              "devices.%s.instances.%s.commandClasses.%d.data.%s",
              node, i, COMMAND_CLASS_SENSOR_MULTILEVEL,k)
            topic := fmt.Sprintf("%s/sensors/analogic/%s/%s/%s/%s",
              zway_home, normName(givenName), i, normName(sensorType),
              normName(sensorScale))
            _, err = jsonFloatValue("val.value",sensor)
            if err == nil {  
              gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
                Value: "val.value", Write:false, Type: "float"})
            }
          }
        }
      }
      // check if instance has meter
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_METER)
      if err == nil {
        for k, v := range data {
          if _, err := strconv.Atoi(k); err == nil {
            sensor := v.(map[string]interface{})
            sensorType, err := jsonStringValue("sensorTypeString.value",
              sensor)
            if err != nil {
              log.Printf("Could not get sensor type: %s", err)
              continue
            }
            sensorScale, err := jsonStringValue("scaleString.value",
              sensor)
            if err != nil {
            log.Printf("Could not get sensor scale: %s", err)
              continue
            }
            nkey := fmt.Sprintf(
              "devices.%s.instances.%s.commandClasses.%d.data.%s",
              node, i, COMMAND_CLASS_METER,k)
            topic := fmt.Sprintf("%s/sensors/analogic/%s/%s/%s/%s",
              zway_home, normName(givenName), i, normName(sensorType),
              normName(sensorScale))
            _, err = jsonFloatValue("val.value",sensor)
            if err == nil {  
              gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
                Value: "val.value", Write:false, Type: "float"})
            }
          }
        }
      }
      // check if instance has thermostat set point
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_THERMOSTAT_SET_POINT)
      if err == nil {
        for k, v := range data {
          if _, err := strconv.Atoi(k); err == nil {
            setpoint :=  v.(map[string]interface{})
            setpointType, err := jsonStringValue("modeName.value",
              setpoint)
            if err != nil {
              log.Printf("Could not get set point mode: %s", err)
              continue
            }
            setpointScale, err := jsonStringValue("scaleString.value",
              setpoint)
            if err != nil {
              log.Printf("Could not get setpoint scale: %s", err)
              continue
            }
            nkey := fmt.Sprintf(
              "devices.%s.instances.%s.commandClasses.%d.data.%s",
              node, i, COMMAND_CLASS_THERMOSTAT_SET_POINT,k)
            topic := fmt.Sprintf("%s/actuators/analogic/%s/%s/%s/%s",
              zway_home, normName(givenName), i, normName(setpointType),
              normName(setpointScale))
            _, err = jsonIntValue("val.value",setpoint)
            if err == nil {
              gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
                Value: "val.value", Write:true, Type: "int", Args: []string{ setpointType, } })
            }
          }
        }
      }
      // check if instance has alarm sensor
      data, err = zwaygetcmdclassdata(commandClasses,
        COMMAND_CLASS_ALARM_SENSOR)
      if err == nil {
        for k, v := range data {
          if _, err := strconv.Atoi(k); err == nil {
            alarm :=  v.(map[string]interface{})
            alarmType, err := jsonStringValue("typeString.value",
              alarm)
            if err != nil {
              log.Printf("Could not get alarm type: %s", err)
              continue
            }
            nkey := fmt.Sprintf(
              "devices.%s.instances.%s.commandClasses.%d.data.%s",
              node, i, COMMAND_CLASS_ALARM_SENSOR,k)
            topic := fmt.Sprintf("%s/sensors/analogic/%s/%s/%s",
              zway_home, normName(givenName), i, normName(alarmType))
            _, err = jsonIntValue("sensorState.value",alarm)
            if err == nil {
              gateways = append(gateways, Gateway{Key: nkey, Topic: topic,
                Value: "sensorState.value", Write:false, Type: "int"})
            }
          }
        }
      }
    }
  }
}

func zwayupdategateways(update map[string]interface{}, mqtt_updates chan<- MqttUpdate) {
  if (debug) { log.Print("Update Z-Way devices") }
  for _, g := range gateways {
    //Z-Way is always true
    value := g.GetValue(update)
    if len(value) > 0 {
      if (debug) { log.Printf("ZWAY: %s / Value: %s", g.ToString(), value ) }
      mqtt_updates <- MqttUpdate{Topic: g.Topic, Value: value}
    }
  }
}

func normalizeJson(json map[string]interface{}) map[string]interface{} {
  for k, v := range json {
    if strings.IndexRune(k,'.') > -1 {
      keys := strings.Split(k,".")
      nkey := keys[0]
      rest := strings.Join(keys[1:len(keys)],".")
      tmp := make(map[string]interface{})
      tmp[rest] = v.(map[string]interface{})
      if json[nkey] != nil {
        for k2, v2 := range json[nkey].(map[string]interface{}) {
          tmp[k2] = v2
        }
      }
      json[nkey] = normalizeJson(tmp)
      delete(json, k)
    }
  }
  return json
}

func checkzwayupdate(update string,mqtt_updates chan<- MqttUpdate) {
  var f interface{}
  err := json.Unmarshal([]byte(update), &f)
  if err != nil {
    log.Printf("Error decoding json: %s", err)
  }
  m := f.(map[string]interface{})
  m = normalizeJson(m)
  if zway_timestamp == 0 {
    devices, err := jsonMapValue("devices",m)
    if err != nil {
      log.Printf("devices not found: %s", err)
      return
    }
    zwayparsedevices(devices)
  }
  zwayupdategateways(m,mqtt_updates)
  zway_timestampf, err := jsonFloatValue("updateTime",m)
  if err != nil {
    log.Printf("timestamp not found: %s", err)
    return
  }
  zway_timestamp = int(zway_timestampf)
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
  topic := msg.Topic()
  value := string(msg.Payload())
  for _, g := range gateways {
    if g.Topic == topic { g.Set(value) }
  }
}

func (g *Gateway) Set(value string) {
  if !g.Write {
    if (debug) { log.Printf("MQTT: %s / Readonly", g.ToString()) }
    return
  }
  if g.Get() == value {
    if (debug) { log.Printf("MQTT: %s / Value not changed", g.ToString()) }
    return
  }
  //check value
  switch g.Type {
    case "int":
      if strings.Contains(value,".") {
        value = strings.TrimRight(value,"0.")
      }
      i, err := strconv.Atoi(value)
      if err != nil {
        log.Printf("MQTT: %s / value not int: %s", g.ToString(), value)
        return
      }
      value = fmt.Sprintf("%d",i)
    case "float":
      if strings.Contains(value,".") {
        value = strings.TrimRight(value,"0.")
      }
      f, err := strconv.ParseFloat(value,64)
      if err != nil {
        log.Printf("MQTT: %s / value not float: %s", g.ToString(), value)
        return
      }
      value = fmt.Sprintf("%.3f", f)
  }
  log.Printf("MQTT: %s / Value: %s ", g.ToString(), value)
  key := g.Key
  r := regexp.MustCompile("\\.([0-9]+)(\\.|$)")
  key = r.ReplaceAllString(key, "[$1].")
  r = regexp.MustCompile("\\.data$")
  key = r.ReplaceAllString(key,"")
  args := ""
  if g.Args != nil {
      for _, v := range g.Args {
          args += fmt.Sprintf("'%s',", v)
      }
  } 
  result, _ := zwayget(zway_runapi,fmt.Sprintf("%s.Set(%s%s)", key, args, value))
  if result != "null" {
    log.Printf("Error updating value: %s", result)
  }
}

func (g *Gateway) Get() string {
  if (debug) { log.Print("Setting Z-Way value.") }
  key := g.Key
  r := regexp.MustCompile("\\.([0-9]+)\\.")
  key = r.ReplaceAllString(key, "[$1].")
  result, _ := zwayget(zway_runapi, fmt.Sprintf("%s.%s", key, g.Value))
  return result
}

func zwayget(api string, path string) (string, error) {
  url := fmt.Sprintf("http://%s%s%s", zway_server, api, path)
  if (debug) { log.Printf("Http Get on Z-Way: %s", url) }
  req, err := http.NewRequest("GET",url,nil)
  if err != nil {
    return "", err
  }
  if zway_cookie != nil {
    req.AddCookie(zway_cookie)
  }
  rsp, err := http_client.Do(req)
  if err != nil {
    return "", err
  }
  defer rsp.Body.Close()
  bdy, err := ioutil.ReadAll(rsp.Body)
  if err != nil {
    return "", err
  }
  result := string(bdy)
  return result, nil
}


func main() {
  //start profiling
  if len(profile) > 0 {
    log.Print("Profiling enabled")
    cfg := pfl.Config{}
    if profile=="mem" || profile=="all" {
      cfg.MemProfile = true
    }
    if profile=="cpu" || profile=="all" {
      cfg.CPUProfile = true
    }
    defer pfl.Start(&cfg).Stop()
  }
  //print informations given
  log.Print("Starting Z-Way to mqtt gateway...")
  log.Printf("Z-Way server: %s", zway_server)
  if len(zway_password) > 0 {
    log.Printf("Z-Way user: %s", zway_username)
  } else {
    log.Print("Not using authentication as no password given.")
  }
  log.Printf("Z-Way refresh rate: %d", zway_refresh)
  log.Printf("MQTT server: %s", mqtt_server)

  //authtenticate to zway
  if len(zway_password) > 0 {
    authzway()
  }

  //connect and subscribe to mqtt
  //prepare
  opts := MQTT.NewClientOptions()
  opts.AddBroker(mqtt_protocol+"://"+mqtt_server)
  opts.SetClientID("ZWayMQTT")
  opts.SetDefaultPublishHandler(f)
  opts.SetAutoReconnect(true)
  if len(mqtt_username) > 0 && len(mqtt_password) > 0 {
    opts.SetUsername(mqtt_username)
    opts.SetPassword(mqtt_password)
  }

  //Connect
  mqtt := MQTT.NewClient(opts)
  if token := mqtt.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  //create the control channel
  quit := make(chan struct{})
  defer close(quit)

  //create zway update channel
  zway_updates := make(chan string,3)
  defer close(zway_updates)

  //create mqtt update channel
  mqtt_updates := make(chan MqttUpdate,20)
  defer close(mqtt_updates)

  //create the zway refresh timer
  refreshes := time.NewTicker(time.Second * time.Duration(zway_refresh)).C

  //make initial refreshe
  zway_updates <- getzway()

  //subscribe only when zway started
  subject := zway_home + "/actuators/#"
  if token := mqtt.Subscribe(subject, 1, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }

  //start refreshes
  go func() {
    for _ = range refreshes {
      update := getzway()
      if len(update) > 0 {
        zway_updates <- getzway()
      } else {
        log.Print("Got empty zwave response...")
        if zway_retries < 3 {
          log.Printf("Reinitializing Z-Way for the %d time.", zway_retries)
          authzway()
          zway_retries++
        } else {
          log.Print("Already tested 3 times: stop")
          <-quit
          return
        }
      }
    }
  }()

  //start update parsing
  go func() {
    for zway_update := range zway_updates {
      checkzwayupdate(zway_update,mqtt_updates)
    }
  }()

  //star mqtt updating
  go func() {
    for mqtt_update := range mqtt_updates {
      token := mqtt.Publish(mqtt_update.Topic, 1, true, mqtt_update.Value)
      token.Wait()
    }
  }()

  //start the main loop
  for {
    select {
    case <- quit:
      return
    }
  }
}
