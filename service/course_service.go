package service

import (
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseService struct {
	courseRepository repository.CourseRepository
}

func NewCourseService(cr repository.CourseRepository) CourseService {
	return CourseService{
		courseRepository: cr,
	}
}

func (cs *CourseService) GetAll() ([]model.Course, error) {
	return cs.courseRepository.GetCoursesByBsonDocument(bson.D{})
}

func (cs *CourseService) GetById(id primitive.ObjectID) (model.Course, error) {
	return cs.courseRepository.GetCourseByBsonDocument(bson.D{{"_id", id}})
}

func (cs *CourseService) GetByName(name string) (model.Course, error) {
	return cs.courseRepository.GetCourseByBsonDocument(bson.D{{"name", name}})
}

func (cs *CourseService) AddCourse(course model.Course) (model.Course, error) {
	return cs.courseRepository.SaveCourse(course)
}

func (cs *CourseService) UpdateCourse(course model.Course) (model.Course, error) {
	return cs.courseRepository.UpdateCourse(course)
}

func (cs *CourseService) DeleteCourse(id primitive.ObjectID) error {
	return cs.courseRepository.DeleteCourse(id)
}
