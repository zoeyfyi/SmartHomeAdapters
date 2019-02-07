package infoserver

import fmt "fmt"

type RobotNotFoundError struct {
	ID string
}

func (e *RobotNotFoundError) Error() string {
	return fmt.Sprintf("No robot with ID \"%s\"", e.ID)
}

type StatusRequestFailed struct {
	Message string
}

func (e *StatusRequestFailed) Error() string {
	if e.Message == "" {
		return "Failed to retrive status of robot"
	} else {
		return fmt.Sprintf("Failed to retrive status of robot: %s", e.Message)
	}
}

type InvalidRobotTypeError struct {
	Type string
}

func (e *InvalidRobotTypeError) Error() string {
	return fmt.Sprintf("Invalid robot type \"%s\"", e.Type)
}

type RobotNotTogglableError struct {
	ID   string
	Type string
}

func (e *RobotNotTogglableError) Error() string {
	return fmt.Sprintf("Robot \"%s\" of type \"%s\" cannot be toggled", e.ID, e.Type)
}

type ToggleRequestFailed struct {
	Message string
}

func (e *ToggleRequestFailed) Error() string {
	return fmt.Sprintf("Toggle request failed: %s", e.Message)
}
