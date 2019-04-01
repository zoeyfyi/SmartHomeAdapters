package main

import (
	"log"
)

var boltlockParameters = map[string]Parameter{
	"LockAngle": IntParameter{
		ID:      "LockAngle",
		Name:    "Lock Angle",
		Min:     0,
		Max:     180,
		Default: 80,
	},
	"UnlockAngle": IntParameter{
		ID:      "UnlockAngle",
		Name:    "Unlock Angle",
		Min:     0,
		Max:     180,
		Default: 100,
	},
}

type Boltlock struct{}

func (s *Boltlock) Name() string {
	return "Boltlock"
}

func (s *Boltlock) Type() UsecaseType {
	return ToggleUsecaseType
}

func (s *Boltlock) DefaultParameters() []Parameter {
	parameters := make([]Parameter, 0, len(boltlockParameters))
	for _, p := range boltlockParameters {
		parameters = append(parameters, p)
	}

	return parameters
}

func (s *Boltlock) GetParameter(id string) *Parameter {
	param, ok := boltlockParameters[id]
	if !ok {
		return nil
	}
	return &param
}

func (s *Boltlock) DefaultToggleStatus() bool {
	return false
}

func (s *Boltlock) DefaultRangeStatus() (int, int, int) {
	log.Fatalf("DefaultRangeStatus called on boltlock")
	return 0, 0, 0
}

func (s *Boltlock) Toggle(value bool, parameters []Parameter, controller RobotController) error {
	log.Printf("toggling boltlock to: %t", value)

	angleID := "LockAngle"
	if !value {
		angleID = "UnlockAngle"
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

func (s *Boltlock) Range(value int64, parameters []Parameter, controller RobotController) error {
	log.Fatalf("Range called on boltlock")
	return nil
}
