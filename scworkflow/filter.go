package scworkflow

var (
	wf_list_to_process []string
)

func InitFilter(wf_to_process []string) {
	wf_list_to_process = wf_to_process
}

func ShouldProcessWorkflowByName(name string) bool {
	if len(wf_list_to_process) > 0 {
		for _, w := range wf_list_to_process {
			if name == w {
				return true
			}
		}
		return false
	}
	return true
}
