package infoserver

import fmt "fmt"

type RobotNotFoundError struct {
	ID string
}

func (e *RobotNotFoundError) Error() string {
	return fmt.Sprintf("No robot with ID \"%s\"", e.ID)
}

type FailedRetreiveStatusError struct{}

func (e *FailedRetreiveStatusError) Error() string {
	return "Failed to retrive status of robot"
}

type InvalidRobotTypeError struct {
	Type string
}

func (e *InvalidRobotTypeError) Error() string {
	return fmt.Sprintf("Invalid robot type \"%s\"", e.Type)
}
