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

func (s *Switch) Description() string {
	return "TODO: switch description"
}

func (s *Switch) Type() UsecaseType {
	return ToggleUsecaseType
}

func (s *Switch) DefaultParameters() []Parameter {
	parameters := make([]Parameter, 0, len(switchParameters))
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
	log.Printf("toggling switch to: %t", value)

	angleID := "OnAngle"
	if !value {
		angleID = "OffAngle"
	}

	// TODO provide better way to get parameters
	var angle int64
	for _, p := range parameters {
		if p, ok := p.(IntParameter); ok && p.ID == angleID {
			angle = p.Current
			break
		}
	}

	log.Printf("setting servo to %d", angle)
	return controller.SetServo(angle)
}

func (s *Switch) Range(value int64, parameters []Parameter, controller RobotController) error {
	log.Fatalf("Range called on switch")
	return nil
}
