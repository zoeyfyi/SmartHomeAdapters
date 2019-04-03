package main

import (
	"fmt"
	"log"
)

var thermostatParameters = map[string]Parameter{
	"LowAngle": IntParameter{
		ID:          "LowAngle",
		Name:        "Low Angle",
		Description: "This is the smallest angle your thermostat can be at",
		Min:         0,
		Max:         180,
		Default:     80,
	},
	"HighAngle": IntParameter{
		ID:          "HighAngle",
		Name:        "High Angle",
		Description: "This is the greatest angle your thermostat can be at",
		Min:         0,
		Max:         180,
		Default:     100,
	},
	"LowTemperature": IntParameter{
		ID:          "LowTemperature",
		Name:        "Low Temperature",
		Description: "This is the lowest temperature marked on your thermostat",
		Min:         273,
		Default:     293,
		Max:         373,
	},
	"HighTemperature": IntParameter{
		ID:          "HighTemperature",
		Name:        "High Temperature",
		Description: "This is the highest temperature marked on your thermostat",
		Min:         273,
		Default:     323,
		Max:         373,
	},
}

type Thermostat struct{}

func (s *Thermostat) Name() string {
	return "thermostat"
}

func (s *Thermostat) Description() string {
	return "Thermostat robot for controlling your temperature."
}
func (s *Thermostat) Type() UsecaseType {
	return RangeUsecaseType
}

func (s *Thermostat) DefaultParameters() []Parameter {
	parameters := make([]Parameter, 0, len(thermostatParameters))
	for _, p := range thermostatParameters {
		parameters = append(parameters, p)
	}

	return parameters
}

func (s *Thermostat) GetParameter(id string) *Parameter {
	param, ok := thermostatParameters[id]
	if !ok {
		return nil
	}
	return &param
}

func (s *Thermostat) DefaultToggleStatus() bool {
	log.Fatalf("DefaultToggleStatus called on thermostat")
	return false
}

func (s *Thermostat) DefaultRangeStatus() (int, int, int) {
	return 0, 0, 0
}

func (s *Thermostat) Toggle(value bool, parameters []Parameter, controller RobotController) error {

	log.Fatalf("Toggle called on thermostat")
	return nil

}

func (s *Thermostat) Range(value int64, parameters []Parameter, controller RobotController) error {
	log.Printf("value = %d", value)

	// need to convert value to an angle
	var angle float64
	var lowAngle int64
	var lowTemp int64
	var highAngle int64
	var highTemp int64
	// TODO - this better
	for _, p := range parameters {
		if p, ok := p.(IntParameter); ok {
			if p.ID == "LowAngle" {
				lowAngle = p.Current
			}
			if p.ID == "LowTemperature" {
				lowTemp = p.Current
			}
			if p.ID == "HighAngle" {
				highAngle = p.Current
			}
			if p.ID == "HighTemperature" {
				highTemp = p.Current
			}
		}
	}

	log.Printf(
		"angle = %f, lowAngle = %d, lowTemp = %d, highAngle = %d, highTemp = %d",
		angle, lowAngle, lowTemp, highAngle, highTemp,
	)

	// calculate the angle required to give the desired temperature
	temperatureRatio := float64(value-lowTemp) / float64(highTemp-lowTemp)
	angle = float64(lowAngle) + float64(highAngle-lowAngle)*temperatureRatio

	log.Printf("angle = %f, temperature ratio = %f", angle, temperatureRatio)

	if angle < 0 || angle > 180 {
		log.Printf("angle invalid!")
		return fmt.Errorf("problem with parameters")
	}

	return controller.SetServo(int64(angle))

}
