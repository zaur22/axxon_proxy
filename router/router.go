package router

/*
	Обычный router
*/
import (
	"axxon_proxy/proxy"
	"axxon_proxy/task"
	"encoding/json"
	"net/http"
)

type Router struct {
	Proxy proxy.Proxy
}

type FetchTaskRequest struct {
	Method  string               `json:"method"`
	Path    string               `json:"path"`
	Headers *map[string][]string `json:"headers"`
	Body    *[]byte              `json:"body"`
}

type GetTasksRequest struct {
	Offset *int `json:"offset"`
	Count  *int `json:"count"`
}

type DeleteTaskRequest struct {
	ID string `json:"ID"`
}

type FetchTaskResponse struct {
	Task task.Task `json:"result"`
}

type GetTasksResponse struct {
	Tasks []task.Task `json:"result"`
}

func (p *Router) FetchTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	ftr := FetchTaskRequest{}
	err := decoder.Decode(&ftr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	task, err := p.Proxy.FetchTask(ftr.Method, ftr.Path, ftr.Headers, ftr.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := json.Marshal(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(res))
}

func (p *Router) GetTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	gt := GetTasksRequest{}
	err := decoder.Decode(&gt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tasks := p.Proxy.GetList(gt.Offset, gt.Count)

	res, err := json.Marshal(&tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(res))
}

func (p *Router) DeleteTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	del := DeleteTaskRequest{}
	err := decoder.Decode(&del)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	taskRes := p.Proxy.DeleteTask(del.ID)
	res, err := json.Marshal(&taskRes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(res))
	return
}
