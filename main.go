package main

import (
  "github.com/godbus/dbus/v5"
  "github.com/shibumi/iwd"
  "fmt"
  "encoding/json"
  "errors"
)

type Ri3_block struct {
  Icon string `json:"icon"`
  State string `json:"state"`
  Text string `json:"text"`
}

func main() {
  defer somethingBadHappened()
  c,err := dbus.ConnectSystemBus()
  if err != nil {
    outputError(err.Error())
    return
  }
  i := IwdStatus{i: iwd.New(c)}
  s, err := i.getConnectedStation()
  if err != nil {
    outputError(err.Error())
    return
  }
  n := i.getNetworkName(s.ConnectedNetwork)
  result,_ := json.Marshal(Ri3_block{
    Icon: "net_wireless",
    State: "Good",
    Text: n,
  })
  fmt.Printf("%v\n", string(result))
}

func outputError(message string) {
  result,_ := json.Marshal(Ri3_block{
    Icon: "net_down",
    State: "Warning",
    Text: message,
  })
  fmt.Printf("%v\n", string(result))
}

func (s *IwdStatus) getConnectedStation() (iwd.Station, error) {
  for _,station := range s.i.Stations {
    if station.State == "connected" {
      return station, nil
    }
  }
  return iwd.Station{}, errors.New("No connection")
}

func (s *IwdStatus) getNetworkName(path dbus.ObjectPath) string {
  for _,network := range s.i.Networks {
    if network.Path == path {
      return network.Name
    }
  }
  return ""
}

func somethingBadHappened() {
  if r := recover(); r != nil {
    outputError("error")
  }
}

type IwdStatus struct {
  i iwd.Iwd
}

