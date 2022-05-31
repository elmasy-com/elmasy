package portscan

type Port struct {
	Port  int
	State State
}
type Result []Port

// GetPortsInt returns []int{ports} from r with state s.
func (r *Result) GetPortsInt(s State) []int {

	v := make([]int, 0)

	for i := range *r {
		if (*r)[i].State == s {
			v = append(v, (*r)[i].Port)
		}
	}

	return v
}

// GetPorts returns []Port from r with state s.
func (r *Result) GetPorts(s State) []Port {

	v := make([]Port, 0)

	for i := range *r {
		if (*r)[i].State == s {
			v = append(v, (*r)[i])
		}
	}

	return v
}

// Len return the number of ports with state s.
func (r *Result) Len(s State) int {

	v := 0

	for i := range *r {
		if (*r)[i].State == s {
			v++
		}
	}

	return v
}

// Add result to port.
// If port is not in the scanned ports, silently ignore it
func (r *Result) addResult(port int, state State) {

	for i := range *r {
		if (*r)[i].Port == port {
			(*r)[i].State = state
		}
	}
}
