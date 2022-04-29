package portscan

type Port struct {
	Port  int
	State State
}
type Result []Port

// GetPorts returns the ports from  r with state s.
func (r *Result) GetPortsInt(s State) []int {

	v := make([]int, 0)

	for i := range *r {
		if (*r)[i].State == s {
			v = append(v, (*r)[i].Port)
		}
	}

	return v
}

func (r *Result) GetPorts(s State) []Port {

	v := make([]Port, 0)

	for i := range *r {
		if (*r)[i].State == s {
			v = append(v, (*r)[i])
		}
	}

	return v
}

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
