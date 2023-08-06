package client_resourceman

import messages "src/messages/client_resourceman"

// JobRequest is a struct that contains the job request information.
type JobRequest struct {
	RequestType string
	Jobs        []Job
}

type Job struct {
	InputFile    string
	OutputFile   string
	ReducerCount int32
	JobBinary    []byte
}

func (jr *JobRequest) GetReqType() string {
	return jr.RequestType
}

// RequestInterface is an interface that contains the request information.
type RequestInterface interface {
	GetReqType() string
}

func (p *ProtoHandler) fetchJobRequest(msg *messages.ClientJobRequest_Jobs_) (req RequestInterface) {

	p.logger.Info("Received job request from the client.")
	if msg.Jobs == nil {
		p.logger.Error("No jobs found in the request.")
		return
	}
	req = &JobRequest{
		RequestType: "job",
		Jobs:        make([]Job, len(msg.Jobs.Jobs)),
	}

	for i, job := range msg.Jobs.Jobs {
		req.(*JobRequest).Jobs[i] = Job{
			InputFile:    job.InputFileName,
			OutputFile:   job.OutputFileName,
			ReducerCount: job.ReducerCount,
			JobBinary:    job.JobSo,
		}
	}
	return req
}

func (p *ProtoHandler) HandleJobResponse(err error) {

	if err == nil {
		p.logger.Info("Job completed successfully.")

		res := &messages.ResourceManagerMessage_JobResponse_{
			JobResponse: &messages.ResourceManagerMessage_JobResponse{
				ErrorCode: messages.ResourceManagerMessage_NO_ERROR,
			},
		}

		wrapper := &messages.ResourceManagerMessage{
			Message: res,
		}

		p.MsgHandler().ResourceManagerResponseSend(wrapper)
		return

	}
	switch err.Error() {

	case "FILE_NOT_FOUND":
		res := &messages.ResourceManagerMessage_JobResponse_{
			JobResponse: &messages.ResourceManagerMessage_JobResponse{
				ErrorCode: messages.ResourceManagerMessage_FILE_NOT_FOUND,
			},
		}

		wrapper := &messages.ResourceManagerMessage{
			Message: res,
		}

		p.MsgHandler().ResourceManagerResponseSend(wrapper)

	case "SERVER_ERROR":

	case "JOB_INVALID":

	}

}
