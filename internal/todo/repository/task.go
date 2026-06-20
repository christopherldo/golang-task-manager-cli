package repository

import (
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"chrisldo.com/todo-cli/internal/todo/models"
)

type TaskRepository struct {
	mu               sync.RWMutex
	cachedTasks      []models.Task
	cachedLastTaskId *int
	store            Store
}

func NewTaskRepository(store Store) *TaskRepository {
	return &TaskRepository{
		cachedTasks: []models.Task{},
		store:       store,
	}
}

func (r *TaskRepository) AppendTaskToDatabase(taskToBeAdded models.Task) error {
	r.mu.Lock()

	r.cachedTasks = append(r.cachedTasks, taskToBeAdded)
	tasksToWrite := slices.Clone(r.cachedTasks)
	r.cachedLastTaskId = &taskToBeAdded.ID

	r.mu.Unlock()

	err := r.store.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("appending task to the database: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetOneTaskFromDatabase(taskId int) (models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx := slices.IndexFunc(r.cachedTasks, func(task models.Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		return models.Task{}, fmt.Errorf("not found with this ID: %d", taskId)
	}

	return r.cachedTasks[idx], nil
}

func (r *TaskRepository) GetAllTasksFromDatabase() []models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.cachedTasks
}

func (r *TaskRepository) UpdateTaskOnDatabase(taskToBeUpdated models.Task) error {
	r.mu.Lock()

	idx := slices.IndexFunc(r.cachedTasks, func(task models.Task) bool {
		return task.ID == taskToBeUpdated.ID
	})

	if idx == -1 {
		r.mu.Unlock()
		return fmt.Errorf("not found with this id: %d", taskToBeUpdated.ID)
	}

	r.cachedTasks[idx] = taskToBeUpdated

	tasksToWrite := slices.Clone(r.cachedTasks)

	r.mu.Unlock()

	err := r.store.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("saving task do the database: %w", err)
	}

	return nil
}

func (r *TaskRepository) MarkTaskAsDone(taskId int) error {
	r.mu.Lock()

	idx := slices.IndexFunc(r.cachedTasks, func(task models.Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		r.mu.Unlock()
		return fmt.Errorf("not found with this ID: %d", taskId)
	}

	r.cachedTasks[idx].IsDone = true
	tasksToWrite := slices.Clone(r.cachedTasks)

	r.mu.Unlock()

	err := r.store.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("marking task as done: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetLastTaskId() int {
	if r.cachedLastTaskId != nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		return *r.cachedLastTaskId
	}

	tasks := r.GetAllTasksFromDatabase()

	if len(tasks) == 0 {
		return 0
	}

	lastIndex := len(tasks) - 1

	lastTask := tasks[lastIndex]

	return lastTask.ID
}

func (r *TaskRepository) DeleteTaskFromDatabase(taskId int) error {
	r.mu.Lock()

	idx := slices.IndexFunc(r.cachedTasks, func(task models.Task) bool {
		return task.ID == taskId
	})

	if idx == -1 {
		r.mu.Unlock()
		return fmt.Errorf("not found with this id: %d", taskId)
	}

	r.cachedTasks = slices.Delete(r.cachedTasks, idx, idx+1)

	tasksToWrite := slices.Clone(r.cachedTasks)

	r.mu.Unlock()

	err := r.store.WriteDatabase(tasksToWrite)

	if err != nil {
		return fmt.Errorf("saving tasks to the database: %w", err)
	}

	return nil
}

func (r *TaskRepository) LoadDatabaseToMemory() error {
	bytes, err := r.store.ReadDatabase()

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &r.cachedTasks)

	if err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	lastTaskId := r.GetLastTaskId()
	r.cachedLastTaskId = &lastTaskId

	return nil
}
