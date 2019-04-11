package proxy

/*
Пакет proxy выполняет запросы,
и работает с историей запросов.

Сами операции выполняются в воркерах. В ворекеры
операции передаются через каналы в виде реализации
интерфейса job. По сути решение - вариация паттерна Action
*/
import (
	"axxon_proxy/task"
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	workerCount    = 10
	jobChannelSize = workerCount
)

var client Client

func init() {
	client = &http.Client{
		Timeout: time.Second * 10,
	}
}

type Proxy interface {
	FetchTask(
		method string,
		path string,
		headers *map[string][]string,
		body *[]byte,
	) (task.Task, error)
	GetList(offset *int, count *int) []task.Task
	DeleteTask(ID string) error
	StopWorkers()
}

type job interface {
	exec(tm task.TaskManager)
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type proxy struct {
	jobChan chan job
	wg      sync.WaitGroup
}

func NewProxy(tm task.TaskManager) Proxy {

	jobChannel := make(chan job, jobChannelSize)

	var p = &proxy{
		jobChan: jobChannel,
		wg:      sync.WaitGroup{},
	}

	for i := 0; i < workerCount; i++ {
		go worker(tm, &p.wg, jobChannel)
		p.wg.Add(1)
	}

	return p
}

type fetchTaskJob struct {
	method  string
	path    string
	headers *map[string][]string
	body    *[]byte
	result  chan task.Task
	errChan chan error
}

type getTasksJob struct {
	offset *int
	count  *int
	result chan []task.Task
}

type deleteTaskJob struct {
	ID  string
	err chan error
}

func (p *proxy) FetchTask(
	method string,
	path string,
	headers *map[string][]string,
	body *[]byte,
) (task.Task, error) {
	job := &fetchTaskJob{
		method:  method,
		path:    path,
		headers: headers,
		body:    body,
		result:  make(chan task.Task),
		errChan: make(chan error),
	}

	p.jobChan <- job

	select {
	case tsk := <-job.result:
		return tsk, nil
	case err := <-job.errChan:
		return task.Task{}, err
	}
}

func (p *proxy) GetList(offset *int, count *int) []task.Task {
	job := &getTasksJob{
		offset: offset,
		count:  count,
		result: make(chan []task.Task),
	}

	p.jobChan <- job
	return <-job.result
}

func (p *proxy) DeleteTask(ID string) error {
	job := &deleteTaskJob{
		ID:  ID,
		err: make(chan error),
	}

	p.jobChan <- job
	return <-job.err
}

func (p *proxy) StopWorkers() {
	close(p.jobChan)
	p.wg.Wait()
}

func (j *fetchTaskJob) exec(tm task.TaskManager) {

	var body *bytes.Reader
	if j.body != nil {
		body = bytes.NewReader(*j.body)
	} else {
		body = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(j.method, j.path, body)
	if err != nil {
		j.errChan <- err
		return
	}

	if j.headers != nil {
		req.Header = *j.headers
	}

	resp, err := client.Do(req)
	if err != nil {
		j.errChan <- err
		return
	}
	resp.Body.Close()

	t, err := task.ParseTask(resp)
	if err != nil {
		j.errChan <- err
		return
	}

	tm.AddTask(t)
	j.result <- t
}

func (g *getTasksJob) exec(tm task.TaskManager) {
	g.result <- tm.GetList(g.offset, g.count)
}

func (d *deleteTaskJob) exec(tm task.TaskManager) {
	ok := tm.DeleteTask(d.ID)
	if !ok {
		err := fmt.Errorf("The task with such ID is missing. ID: %v", d.ID)
		d.err <- err
	}

	d.err <- nil
}

func worker(tm task.TaskManager, wg *sync.WaitGroup, jobChan <-chan job) {
	defer wg.Done()
	for job := range jobChan {
		job.exec(tm)
	}
}
