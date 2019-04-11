package task

/*
	Пакет для хранения тасков и потокобезопасной работы с ними
*/
import (
	"log"
	"net/http"
	"sync"

	"github.com/gofrs/uuid"
)

func NewTaskManager() TaskManager {
	return &story{
		mx:    sync.RWMutex{},
		tasks: make(map[string]Task),
	}
}

type TaskManager interface {
	AddTask(t Task)
	GetList(startFrom *int, count *int) []Task
	DeleteTask(ID string) bool
}

type Task struct {
	ID            string      `json:"id"`
	HTTPStatus    string      `json:"http-status"`
	Header        http.Header `json:"headers"`
	ContentLength int64       `json:"content-length"`
}

type story struct {
	mx    sync.RWMutex
	tasks map[string]Task
}

func (s *story) AddTask(t Task) {
	s.mx.Lock()
	s.tasks[t.ID] = t
	s.mx.Unlock()
}

func (s *story) GetList(offset *int, count *int) []Task {
	var tasks = []Task{}
	if offset == nil {
		offset = new(int)
		*offset = 0
	}
	if count == nil {
		count = new(int)
		*count = len(s.tasks)
	}

	s.mx.RLock()
	for _, val := range s.tasks {
		if *offset > 0 {
			*offset--
		} else {
			tasks = append(tasks, val)
			*count--
		}
	}
	s.mx.RUnlock()
	return tasks
}

func (s *story) DeleteTask(ID string) bool {

	s.mx.RLock()
	_, ok := s.tasks[ID]
	s.mx.RUnlock()
	if !ok {
		return false
	}

	s.mx.Lock()
	delete(s.tasks, ID)
	s.mx.Unlock()

	return true
}

func ParseTask(resp *http.Response) (Task, error) {
	task := Task{}
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
		return task, err
	}

	task.ID = u.String()
	task.HTTPStatus = resp.Status
	task.Header = resp.Header
	task.ContentLength = resp.ContentLength

	return task, nil
}
