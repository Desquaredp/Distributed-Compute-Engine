package client_resourceman

import (
	"errors"
	messages "src/messages/client_resourceman"
)

type ResponseInterface interface {
	GetType() string
}

type JobResponse struct {
	responseType string
	Err          error
}

func (jr *JobResponse) GetType() string {
	return jr.responseType
}

// HandleMRJobRequest sends a job request to the ResourceManager.
func (p *ProtoHandler) HandleMRJobRequest(inputFile, outputFile []string, reducerCount []int32, jobBinaryPath [][]byte) {

	p.logger.Info("Sending job request to the ResourceManager.")
	msg := &messages.ClientJobRequest{
		ClientMessage: &messages.ClientJobRequest_Jobs_{
			Jobs: &messages.ClientJobRequest_Jobs{

				Jobs: make([]*messages.ClientJobRequest_Job, len(inputFile)),
			},
		},
	}

	for i := range inputFile {
		msg.ClientMessage.(*messages.ClientJobRequest_Jobs_).Jobs.Jobs[i] = &messages.ClientJobRequest_Job{
			InputFileName:  inputFile[i],
			OutputFileName: outputFile[i],
			ReducerCount:   reducerCount[i],
			JobSo:          jobBinaryPath[i],
		}
	}

	p.msgHandler.ClientRequestSend(msg)

}

func (p *ProtoHandler) fetchJobResponse(msg *messages.ResourceManagerMessage_JobResponse_) (res ResponseInterface) {

	res = &JobResponse{
		responseType: "job",
		Err:          errors.New(msg.JobResponse.ErrorCode.String()),
	}

	return
}

type ProgressUpdate struct {
	ResType            string
	ProgressUpdateType string
	TaskProgress       *TaskProgress
}
type TaskProgress struct {
	TaskType   string
	TaskNumber int32
	TotalTasks int32
}

func (pu *ProgressUpdate) GetType() string {
	return pu.ResType
}

func (p *ProtoHandler) fetchProgressResponse(msg *messages.ResourceManagerMessage_ProgressUpdate_) (Res ResponseInterface) {

	if msg.ProgressUpdate == nil {
		return nil
	}

	Res = &ProgressUpdate{
		ResType:            "progress",
		ProgressUpdateType: msg.ProgressUpdate.UpdateType.String(),
		TaskProgress:       &TaskProgress{},
	}

	if msg.ProgressUpdate.TaskProgress != nil {

		Res.(*ProgressUpdate).TaskProgress.TaskType = msg.ProgressUpdate.TaskProgress.TaskType.String()
		Res.(*ProgressUpdate).TaskProgress.TaskNumber = msg.ProgressUpdate.TaskProgress.TaskNumber
		Res.(*ProgressUpdate).TaskProgress.TotalTasks = msg.ProgressUpdate.TaskProgress.TotalTasks

	} else {
		Res.(*ProgressUpdate).TaskProgress = nil
	}

	return

}
