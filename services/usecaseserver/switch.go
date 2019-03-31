package main

import (
	"log"
)

var switchParameters = map[string]Parameter{
	"OnAngle": IntParameter{
		ID:      "OnAngle",
		Name:    "On Angle",
		Min:     0,
		Max:     180,
		Default: 80,
	},
	"OffAngle": IntParameter{
		ID:      "OffAngle",
		Name:    "Off Angle",
		Min:     0,
		Max:     180,
		Default: 100,
	},
}

type Switch struct{}

func (s *Switch) Name() string {
	return "switch"
}

func (s *Switch) Type() UsecaseType {
	return ToggleUsecaseType
}

func (s *Switch) DefaultParameters() []Parameter {
	parameters := make([]Parameter, len(switchParameters))
	for _, p := range switchParameters {
		parameters = append(parameters, p)
	}

	return parameters
}

func (s *Switch) GetParameter(id string) *Parameter {
	param, ok := switchParameters[id]
	if !ok {
		return nil
	}
	return &param
}

func (s *Switch) DefaultToggleStatus() bool {
	return false
}

func (s *Switch) DefaultRangeStatus() (int, int, int) {
	log.Fatalf("DefaultRangeStatus called on switch")
	return 0, 0, 0
}

func (s *Switch) Toggle(value bool, parameters []Parameter, controller RobotController) error {
	angleID := "OnAngle"
	if !value {
		angleID = "OffAngle"
	}

	angle := (*s.GetParameter(angleID)).(IntParameter).Current
	return controller.SetServo(angle)
}

func (s *Switch) Range(value int64, parameters []Parameter, controller RobotController) error {
	log.Fatalf("Range called on switch")
	return nil
}
